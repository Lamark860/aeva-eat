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
  const loading = ref(false)
  const error = ref(null)

  async function load() {
    loading.value = true
    error.value = null
    try {
      // Pull a generous slice; backend sort=new orders by created_at desc.
      const { data } = await http.get('/places', { params: { sort: 'new', limit: 60 } })
      places.value = data.places || []
    } catch (e) {
      error.value = e
      places.value = []
    } finally {
      loading.value = false
    }
  }

  // Group places into buckets (current week + collapsed older).
  const buckets = computed(() => {
    if (!places.value.length) return { current: null, older: [] }
    const now = new Date()
    const thisWeekStart = weekStart(now)

    const current = { label: formatWeekRange(thisWeekStart), start: thisWeekStart, items: [] }
    const olderMap = new Map()

    for (const p of places.value) {
      const date = new Date(p.created_at || p.updated_at || Date.now())
      const ws = weekStart(date)
      if (ws.getTime() === thisWeekStart.getTime()) {
        current.items.push({ ...p, _date: date })
        continue
      }
      const b = bucketKey(ws, now)
      if (!olderMap.has(b.key)) {
        olderMap.set(b.key, { ...b, label: formatBucketLabel(b), items: [] })
      }
      olderMap.get(b.key).items.push({ ...p, _date: date })
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
    loading,
    error,
    buckets,
    load,
    formatVisitCaption,
  }
}
