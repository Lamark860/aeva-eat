import { ref, computed } from 'vue'
import http from '../api/http'

// Color palette for AuthorTag — 4 buckets, deterministic by user id.
const AUTHOR_COLORS = ['terra', 'ochre', 'moss', 'plum']
export function authorColor(userId) {
  if (userId == null) return AUTHOR_COLORS[0]
  const n = Number(userId)
  if (Number.isNaN(n)) {
    // string fallback — sum char codes
    let s = 0
    for (let i = 0; i < userId.length; i++) s += userId.charCodeAt(i)
    return AUTHOR_COLORS[s % AUTHOR_COLORS.length]
  }
  return AUTHOR_COLORS[Math.abs(n) % AUTHOR_COLORS.length]
}

// Returns Monday 00:00 of the week containing `d` (UTC-agnostic local).
function weekStart(d) {
  const x = new Date(d)
  const day = (x.getDay() + 6) % 7 // mon=0..sun=6
  x.setHours(0, 0, 0, 0)
  x.setDate(x.getDate() - day)
  return x
}

const RU_MONTH_GEN = [
  'янв', 'фев', 'мар', 'апр', 'мая', 'июн',
  'июл', 'авг', 'сен', 'окт', 'ноя', 'дек',
]
const RU_MONTH_FULL = [
  'январь', 'февраль', 'март', 'апрель', 'май', 'июнь',
  'июль', 'август', 'сентябрь', 'октябрь', 'ноябрь', 'декабрь',
]
const RU_WD = ['вс', 'пн', 'вт', 'ср', 'чт', 'пт', 'сб']

// "4–9 мая" or "29 апр – 4 мая"
function formatWeekRange(start) {
  const end = new Date(start)
  end.setDate(end.getDate() + 6)
  const sameMonth = start.getMonth() === end.getMonth()
  if (sameMonth) {
    return `${start.getDate()}–${end.getDate()} ${RU_MONTH_GEN[start.getMonth()]}`
  }
  return `${start.getDate()} ${RU_MONTH_GEN[start.getMonth()]} – ${end.getDate()} ${RU_MONTH_GEN[end.getMonth()]}`
}

// caption: "Probka · сб"
export function formatVisitCaption(name, date) {
  const d = new Date(date)
  return `${name} · ${RU_WD[d.getDay()]}`
}

// Older buckets: collapse multiple weeks into a coarser label after ~3 weeks ago.
function bucketKey(start, now) {
  const ageDays = Math.floor((now - start) / 86400000)
  if (ageDays < 7) return { kind: 'week', key: 'this', start }
  if (ageDays < 28) return { kind: 'week', key: start.toISOString().slice(0, 10), start }
  // monthly buckets
  const m = `${start.getFullYear()}-${String(start.getMonth() + 1).padStart(2, '0')}`
  return { kind: 'month', key: m, year: start.getFullYear(), month: start.getMonth() }
}

function formatBucketLabel(bucket) {
  if (bucket.kind === 'week') return formatWeekRange(bucket.start)
  return RU_MONTH_FULL[bucket.month]
}

export function useFeed() {
  const places = ref([])
  const notes = ref([])
  const wishlist = ref([])
  const events = ref([])
  const users = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function load() {
    loading.value = true
    error.value = null
    try {
      // Q4 — лента event-driven: /api/feed диктует хронологию (review/note),
      // /places и /notes даёт rich-данные для рендера, /wishlist/all — планы,
      // /users — резолв attendees-id'шников в имена/аватары.
      const [eventsRes, placesRes, notesRes, wishRes, usersRes] = await Promise.allSettled([
        http.get('/feed', { params: { limit: 100 } }),
        http.get('/places', { params: { sort: 'new', limit: 60 } }),
        http.get('/notes'),
        http.get('/wishlist/all'),
        http.get('/users'),
      ])
      events.value   = eventsRes.status === 'fulfilled' ? (eventsRes.value.data || []) : []
      places.value   = placesRes.status === 'fulfilled' ? (placesRes.value.data.places || []) : []
      notes.value    = notesRes.status  === 'fulfilled' ? (notesRes.value.data || []) : []
      wishlist.value = wishRes.status   === 'fulfilled' ? (wishRes.value.data || []) : []
      users.value    = usersRes.status  === 'fulfilled' ? (usersRes.value.data || []) : []
    } catch (e) {
      error.value = e
      events.value = []
      places.value = []
      notes.value = []
      wishlist.value = []
      users.value = []
    } finally {
      loading.value = false
    }
  }

  // Q4 — группировка place-артефактов по (place_id, day). Все review-события
  // одного места за один день мерджатся: attendees объединяются (тот же
  // визит, но соавторы могли запостить отдельно), videos накапливаются
  // как список URL (один отзыв с видео не должен расширять соседа без видео).
  // Разные дни → разные карточки — это и закрывает «раз в месяц туда же»
  // как отдельные артефакты в ленте.
  const items = computed(() => {
    const placeMap = new Map(places.value.map(p => [p.id, p]))
    const userMap  = new Map(users.value.map(u => [u.id, u]))
    const out = []

    const reviewGroups = new Map()
    for (const ev of events.value) {
      if (ev.kind !== 'review_added' || !ev.place_id) continue
      const place = placeMap.get(ev.place_id)
      if (!place) continue
      const day = (ev.occurred_at || '').slice(0, 10)
      const key = `${ev.place_id}|${day}`
      let g = reviewGroups.get(key)
      if (!g) {
        g = {
          _kind: 'place',
          _date: new Date(ev.occurred_at),
          _attendeeIds: new Set(),
          _videos: [],
          _foodSum: 0, _serviceSum: 0, _vibeSum: 0, _ratingCount: 0,
          id: `${ev.place_id}-${day}`,
          _placeId: ev.place_id,
          ...place,
          // override: per-event-grouped поля. Видео и рейтинги — из событий,
          // не place-уровня. Иначе все визиты места показывали бы
          // одинаковую общую оценку и любой видео-флаг расширял всех.
          videos: undefined,
          video_url: undefined,
          avg_food: undefined,
          avg_service: undefined,
          avg_vibe: undefined,
          // created_at заменим ниже на дату самой ранней встречи в группе —
          // нужно для caption'а «Place · вс/сб» в ArtifactCard.
          created_at: undefined,
        }
        reviewGroups.set(key, g)
      }
      for (const id of (ev.attendees || [])) g._attendeeIds.add(id)
      if (ev.video_url) g._videos.push(ev.video_url)
      if (ev.food_rating != null) {
        g._foodSum    += ev.food_rating
        g._serviceSum += ev.service_rating ?? 0
        g._vibeSum    += ev.vibe_rating    ?? 0
        g._ratingCount += 1
      }
      // _date: берём самое раннее событие дня — карточка датируется началом визита
      const d = new Date(ev.occurred_at)
      if (d < g._date) g._date = d
    }
    for (const g of reviewGroups.values()) {
      g._attendees = [...g._attendeeIds]
        .sort((a, b) => a - b)
        .map(id => userMap.get(id))
        .filter(Boolean)
      g.videos = g._videos
      // Усредняем оценки внутри группы (на случай нескольких отзывов в один день).
      if (g._ratingCount > 0) {
        g.avg_food    = +(g._foodSum    / g._ratingCount).toFixed(1)
        g.avg_service = +(g._serviceSum / g._ratingCount).toFixed(1)
        g.avg_vibe    = +(g._vibeSum    / g._ratingCount).toFixed(1)
        g.review_count = g._ratingCount
      }
      g.created_at = g._date.toISOString()
      delete g._videos
      delete g._attendeeIds
      delete g._foodSum
      delete g._serviceSum
      delete g._vibeSum
      delete g._ratingCount
      out.push(g)
    }

    for (const n of notes.value) {
      out.push({ _kind: 'note', _date: new Date(n.created_at), id: `note-${n.id}`, _note: n })
    }
    for (const w of wishlist.value) {
      const date = w.is_struck && w.struck_at ? w.struck_at : w.created_at
      out.push({
        _kind: 'wishlist',
        _date: new Date(date),
        id: `wish-${w.user_id}-${w.place.id}`,
        _wish: w,
      })
    }
    out.sort((a, b) => b._date - a._date)
    return out
  })

  const buckets = computed(() => {
    const list = items.value
    if (!list.length) return { current: null, older: [] }
    const now = new Date()
    const thisWeekStart = weekStart(now)

    const current = { label: formatWeekRange(thisWeekStart), start: thisWeekStart, items: [] }
    const olderMap = new Map()

    for (const it of list) {
      const ws = weekStart(it._date)
      if (ws.getTime() === thisWeekStart.getTime()) {
        current.items.push(it)
        continue
      }
      const b = bucketKey(ws, now)
      if (!olderMap.has(b.key)) {
        olderMap.set(b.key, { ...b, label: formatBucketLabel(b), items: [] })
      }
      olderMap.get(b.key).items.push(it)
    }

    const older = Array.from(olderMap.values()).sort((a, b) => {
      const at = (a.start || new Date(a.year, a.month, 1)).getTime()
      const bt = (b.start || new Date(b.year, b.month, 1)).getTime()
      return bt - at
    })

    return { current: current.items.length ? current : { ...current, items: [] }, older }
  })

  return {
    places,
    notes,
    loading,
    error,
    buckets,
    load,
    formatVisitCaption,
  }
}
