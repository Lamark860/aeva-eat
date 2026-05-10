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

  // Q4 — группировка place-артефактов по (place_id, date, attendees-set).
  // Это закрывает «раз в месяц туда же» как отдельные карточки, а не одну
  // плитку с устаревшей датой. Wishlist и notes не группируются, т.к. они
  // разовые события.
  const items = computed(() => {
    const placeMap = new Map(places.value.map(p => [p.id, p]))
    const userMap  = new Map(users.value.map(u => [u.id, u]))
    const out = []

    // Группы review-событий: ключ — place_id|YYYY-MM-DD|sorted-attendees.
    const reviewGroups = new Map()
    for (const ev of events.value) {
      if (ev.kind !== 'review_added' || !ev.place_id) continue
      const place = placeMap.get(ev.place_id)
      if (!place) continue
      const day = (ev.occurred_at || '').slice(0, 10)
      const att = (ev.attendees || []).slice().sort((a, b) => a - b)
      const key = `${ev.place_id}|${day}|${att.join(',')}`
      if (!reviewGroups.has(key)) {
        reviewGroups.set(key, {
          _kind: 'place',
          _date: new Date(ev.occurred_at),
          _attendees: att.map(id => userMap.get(id)).filter(Boolean),
          // id для tilt-хэша: нумеровать, чтобы повторные визиты в одно
          // место в разные дни получали разные наклоны.
          id: `${ev.place_id}-${day}`,
          _placeId: ev.place_id,
          ...place,
        })
      }
    }
    out.push(...reviewGroups.values())

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
