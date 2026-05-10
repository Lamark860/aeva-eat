<template>
  <div class="sb-paper sb-grain sb-screen find">
    <header class="find-header">
      <div class="sb-section-head" style="padding: 0">
        <h2>Найти</h2>
        <span class="sub">куда сходить сегодня</span>
      </div>
    </header>

    <!-- Search row: stamp-styled input + filter pin -->
    <div class="find-search">
      <div class="search-stamp">
        <span class="search-glyph">⌕</span>
        <input
          v-model="placesStore.filters.search"
          type="text"
          class="search-input"
          placeholder="место, кухня, город…"
          @input="debouncedFetch"
        />
        <button
          v-if="placesStore.filters.search"
          class="clear-x"
          aria-label="Очистить поиск"
          @click="clearSearch"
        >
×
</button>
      </div>
      <button
        class="filter-pin"
        :class="{ on: activeFilterCount > 0 }"
        type="button"
        data-bs-toggle="offcanvas"
        data-bs-target="#placesFilterDrawer"
        aria-label="Фильтры"
      >
        <span class="head" aria-hidden="true"></span>
        <span class="lbl">фильтры</span>
        <span v-if="activeFilterCount > 0" class="count">{{ activeFilterCount }}</span>
      </button>
    </div>

    <!-- 🎲 random + add -->
    <div class="find-actions">
      <button class="dice" type="button" :disabled="!placesStore.places.length" @click="rollDice">
        <span class="glyph">🎲</span>
        <span>мне&nbsp;повезёт</span>
      </button>
      <router-link v-if="auth.isAuthenticated" to="/places/new" class="add-link">
        + добавить место
      </router-link>
    </div>

    <!-- Active filter chips -->
    <div v-if="activeFilterCount > 0" class="find-chips">
      <button v-if="placesStore.filters.is_gem" class="chip" @click="placesStore.filters.is_gem = false; fetchFiltered()">
        ♦ жемчужины ×
      </button>
      <button v-if="placesStore.filters.city" class="chip" @click="placesStore.filters.city = ''; fetchFiltered()">
        {{ placesStore.filters.city }} ×
      </button>
      <button
        v-for="id in placesStore.filters.cuisine_type_ids || []"
        :key="`c-${id}`"
        class="chip"
        @click="removeCuisine(id)"
      >
{{ cuisineName(id) }} ×
</button>
      <button
        v-for="id in placesStore.filters.category_ids || []"
        :key="`cat-${id}`"
        class="chip"
        @click="removeCategory(id)"
      >
{{ categoryName(id) }} ×
</button>
      <button
        v-for="id in placesStore.filters.attended_by || []"
        :key="`att-${id}`"
        class="chip"
        @click="toggleAttended(id); fetchFiltered()"
      >
{{ friendName(id) }} был · ×
</button>
      <button v-if="dateRangeChip" class="chip" @click="clearDateRange">
        {{ dateRangeChip }} ×
      </button>
      <button v-if="placesStore.filters.sort" class="chip" @click="clearSort">
        {{ sortLabel(placesStore.filters.sort) }} ×
      </button>
      <button class="chip reset" @click="resetFilters">сброс</button>
    </div>

    <!-- Shelves (only when no filters/search active) -->
    <template v-if="showShelves">
      <!-- Жемчужины -->
      <section v-if="topGems.length" class="shelf">
        <div class="shelf-head">
          <h3>Жемчужины</h3>
          <router-link to="/gems" class="shelf-all">все →</router-link>
        </div>
        <div class="shelf-row gems">
          <router-link
            v-for="p in topGems"
            :key="p.id"
            :to="`/places/${p.id}`"
            class="shelf-gem"
          >
            <Polaroid
              :src="p.image_url || ''"
              :width="118"
              :height="118"
              :tilt="tiltFor(p.id)"
              gem
              :placeholder="placeholderFor(p.id)"
            >
              <Tape :variant="tapeVariantFor(p.id)" :style="{ top: '-9px', left: '34px', transform: 'rotate(-10deg)' }" />
              <span class="gem-corner"><GemBadge :size="20" /></span>
            </Polaroid>
            <div class="shelf-cap">{{ p.name }}</div>
          </router-link>
        </div>
      </section>

      <!-- По городам — вертикальный список (DESIGN-DECISIONS F1):
           имя серифой слева, count-штамп справа -->
      <section v-if="cityShelf.length" class="shelf">
        <div class="shelf-head">
          <h3>По городам</h3>
        </div>
        <ul class="city-list">
          <li v-for="c in cityShelf" :key="c.name">
            <router-link :to="`/cities/${encodeURIComponent(c.name)}`" class="city-row">
              <span class="city-name">{{ c.name }}</span>
              <span class="city-count">{{ c.count }}</span>
            </router-link>
          </li>
        </ul>
      </section>

      <!-- По кухням -->
      <section v-if="cuisineShelf.length" class="shelf">
        <div class="shelf-head">
          <h3>По кухням</h3>
        </div>
        <div class="shelf-stamps">
          <Stamp
            v-for="c in cuisineShelf"
            :key="c.id"
            kind="moss"
            class="clickable"
            @click="filterByCuisine(c.id)"
          >
{{ c.name }} · {{ c.count }}
</Stamp>
        </div>
      </section>

      <!-- По друзьям — горизонтальная карусель аватарок 60px (DESIGN-DECISIONS §F1) -->
      <section v-if="friends.length" class="shelf">
        <div class="shelf-head">
          <h3>По друзьям</h3>
        </div>
        <div class="shelf-row friends">
          <router-link
            v-for="u in friends"
            :key="u.id"
            :to="`/people/${u.id}`"
            class="friend-tile"
          >
            <span
              class="r-tag sb-author-tag friend-avatar"
              :class="[authorColor(u.id), { 'has-photo': !!u.avatar_url }]"
              :title="u.username"
            >
              <img v-if="u.avatar_url" :src="u.avatar_url" alt="" class="r-ph" />
              <template v-else>{{ (u.username || '?').slice(0, 1).toUpperCase() }}</template>
            </span>
            <span class="friend-name">{{ u.username }}</span>
            <span class="friend-count">{{ u.place_count }}</span>
          </router-link>
        </div>
      </section>
    </template>

    <!-- Results list -->
    <section v-if="!showShelves || (placesStore.filters.search || activeFilterCount > 0)" class="find-results">
      <div v-if="placesStore.loading" class="sb-empty">…</div>
      <div v-else-if="placesStore.places.length === 0" class="sb-empty">
        <div>ничего не нашли</div>
        <div style="margin-top: 8px">
          <button class="add-link" @click="resetFilters">сбросить фильтры</button>
        </div>
      </div>
      <div v-else class="results-list">
        <ResultCard v-for="place in placesStore.places" :key="place.id" :place="place" />
      </div>

      <nav v-if="totalPages > 1" class="find-pagination">
        <button
          class="page-link"
          :disabled="placesStore.page <= 1"
          @click="goToPage(placesStore.page - 1)"
        >
‹
</button>
        <span class="page-meta">{{ placesStore.page }} / {{ totalPages }}</span>
        <button
          class="page-link"
          :disabled="placesStore.page >= totalPages"
          @click="goToPage(placesStore.page + 1)"
        >
›
</button>
      </nav>
    </section>

    <!-- Filter drawer (Bootstrap offcanvas, scrapbook-styled body) -->
    <div
      class="offcanvas offcanvas-bottom find-drawer"
      tabindex="-1"
      id="placesFilterDrawer"
      aria-labelledby="placesFilterDrawerLabel"
      style="height: auto; max-height: 85vh"
    >
      <div class="offcanvas-header">
        <h5 class="offcanvas-title" id="placesFilterDrawerLabel">Фильтры</h5>
        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Закрыть"></button>
      </div>
      <div class="offcanvas-body">
        <div class="mb-3">
          <label class="drawer-label">Город</label>
          <select v-model="placesStore.filters.city" class="form-select">
            <option value="">все</option>
            <option v-for="city in cities" :key="city" :value="city">{{ city }}</option>
          </select>
        </div>
        <div class="mb-3">
          <label class="drawer-label">Кухни</label>
          <MultiSelect
            :modelValue="placesStore.filters.cuisine_type_ids"
            @update:model-value="v => placesStore.filters.cuisine_type_ids = v"
            :options="catalogs.cuisineTypes"
            placeholder="любые"
          />
        </div>
        <div class="mb-3">
          <label class="drawer-label">Категории</label>
          <MultiSelect
            :modelValue="placesStore.filters.category_ids"
            @update:model-value="v => placesStore.filters.category_ids = v"
            :options="catalogs.categories"
            placeholder="любые"
          />
        </div>
        <div class="mb-3">
          <label class="drawer-label">Сортировка</label>
          <select :value="sortBase" @change="onSortChange" class="form-select">
            <option value="">сначала новые</option>
            <option value="rating">по рейтингу ↓</option>
            <option value="rating_asc">по рейтингу ↑</option>
            <option value="name">по названию</option>
            <option value="rating_user">по оценке друга</option>
          </select>
          <select
            v-if="sortBase === 'rating_user'"
            v-model="ratingUserId"
            class="form-select mt-2"
            @change="onRatingUserChange"
          >
            <option value="">— выбери друга —</option>
            <option v-for="u in friends" :key="u.id" :value="u.id">{{ u.username }}</option>
          </select>
        </div>

        <!-- Q5: «кто был» — multi-select аватарок-чипов из друзей круга. -->
        <div class="mb-3">
          <label class="drawer-label">Кто был</label>
          <input
            v-if="friends.length > 10"
            v-model="attendedSearch"
            type="text"
            class="form-control mb-2"
            placeholder="поиск имени…"
          />
          <div class="att-chips">
            <button
              v-for="u in attendedVisible"
              :key="u.id"
              type="button"
              class="att-chip"
              :class="{ on: placesStore.filters.attended_by.includes(u.id) }"
              @click="toggleAttended(u.id)"
            >
              <span
                class="r-tag sb-author-tag"
                :class="[authorColor(u.id), { 'has-photo': !!u.avatar_url }]"
              >
                <img v-if="u.avatar_url" :src="u.avatar_url" alt="" class="r-ph" />
                <template v-else>{{ (u.username || '?').slice(0, 1).toUpperCase() }}</template>
              </span>
              <span class="att-name">{{ u.username }}</span>
            </button>
          </div>
        </div>

        <!-- Q5: «когда» — paper-control date inputs «с/по» + preset-чипы. -->
        <div class="mb-3">
          <label class="drawer-label">Когда</label>
          <div class="date-row">
            <input
              v-model="placesStore.filters.visit_from"
              type="date"
              class="form-control"
              aria-label="с"
            />
            <span class="date-sep">—</span>
            <input
              v-model="placesStore.filters.visit_to"
              type="date"
              class="form-control"
              aria-label="по"
            />
          </div>
          <div class="preset-chips">
            <button
              v-for="p in datePresets"
              :key="p.key"
              type="button"
              class="preset-chip"
              :class="{ on: activePreset === p.key }"
              @click="applyDatePreset(p.key)"
            >
{{ p.label }}
</button>
          </div>
        </div>
        <div class="form-check form-switch mb-3">
          <input
            v-model="placesStore.filters.is_gem"
            class="form-check-input"
            type="checkbox"
            id="gemFilterMobile"
            role="switch"
          />
          <label class="form-check-label" for="gemFilterMobile">только&nbsp;жемчужины&nbsp;♦</label>
        </div>
        <div class="d-grid gap-2 mt-4">
          <button class="btn btn-apply" data-bs-dismiss="offcanvas" @click="fetchFiltered">применить</button>
          <button class="btn btn-link reset-btn" @click="resetFilters">сбросить</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePlacesStore } from '../stores/places'
import { useCatalogsStore } from '../stores/catalogs'
import { useAuthStore } from '../stores/auth'
import MultiSelect from '../components/MultiSelect.vue'
import Polaroid from '../components/scrapbook/Polaroid.vue'
import Tape from '../components/scrapbook/Tape.vue'
import Stamp from '../components/scrapbook/Stamp.vue'
import GemBadge from '../components/scrapbook/GemBadge.vue'
import ResultCard from '../components/scrapbook/ResultCard.vue'
import http from '../api/http'
import { authorColor } from '../composables/useFeed'

const route = useRoute()
const router = useRouter()
const placesStore = usePlacesStore()
const catalogs = useCatalogsStore()
const auth = useAuthStore()

const cities = ref([])
// Список друзей круга для полки «По друзьям» (DESIGN-DECISIONS §F1).
// Сортировка по review_count desc — сверху самые активные. Своего профиля
// не исключаем сознательно: видеть себя в полке полезно.
const friends = ref([])

const activeFilterCount = computed(() => {
  const f = placesStore.filters
  let n = 0
  if (f.city) n++
  if (f.cuisine_type_ids?.length) n++
  if (f.category_ids?.length) n++
  if (f.sort) n++
  if (f.is_gem) n++
  if (f.min_rating) n++
  if (f.attended_by?.length) n++
  if (f.visit_from || f.visit_to) n++
  return n
})

// Q5 — sort распадается на «base» + опционально rating_user-id. Эти computed
// держат UI в синке с filter.sort: если sort = 'rating_user:5', то sortBase
// равно 'rating_user', а ratingUserId — '5'.
const sortBase = computed(() => {
  const s = placesStore.filters.sort || ''
  return s.startsWith('rating_user:') ? 'rating_user' : s
})
const ratingUserId = ref('')
function onSortChange(e) {
  const v = e.target.value
  if (v === 'rating_user') {
    placesStore.filters.sort = ratingUserId.value ? `rating_user:${ratingUserId.value}` : 'rating_user'
  } else {
    placesStore.filters.sort = v
    ratingUserId.value = ''
  }
  fetchFiltered()
}
function onRatingUserChange() {
  if (ratingUserId.value) {
    placesStore.filters.sort = `rating_user:${ratingUserId.value}`
    fetchFiltered()
  }
}
function clearSort() {
  placesStore.filters.sort = ''
  ratingUserId.value = ''
  fetchFiltered()
}

// Q5 — «кто был» multi-select.
const attendedSearch = ref('')
const attendedVisible = computed(() => {
  const q = attendedSearch.value.trim().toLowerCase()
  if (!q) return friends.value
  return friends.value.filter(u => (u.username || '').toLowerCase().includes(q))
})
function toggleAttended(id) {
  const arr = placesStore.filters.attended_by || []
  placesStore.filters.attended_by = arr.includes(id) ? arr.filter(x => x !== id) : [...arr, id]
}
function friendName(id) {
  return friends.value.find(u => u.id === id)?.username || ''
}

// Q5 — preset-чипы для «когда». Конвертируются в visit_from / visit_to.
const datePresets = [
  { key: 'this-year',    label: 'этот год' },
  { key: 'prev-year',    label: 'прошлый год' },
  { key: 'last-30-days', label: 'последние 30 дней' },
]
function isoDate(d) { return d.toISOString().slice(0, 10) }
function applyDatePreset(key) {
  const now = new Date()
  if (key === 'this-year') {
    placesStore.filters.visit_from = `${now.getFullYear()}-01-01`
    placesStore.filters.visit_to   = `${now.getFullYear()}-12-31`
  } else if (key === 'prev-year') {
    const y = now.getFullYear() - 1
    placesStore.filters.visit_from = `${y}-01-01`
    placesStore.filters.visit_to   = `${y}-12-31`
  } else if (key === 'last-30-days') {
    const from = new Date(now)
    from.setDate(from.getDate() - 30)
    placesStore.filters.visit_from = isoDate(from)
    placesStore.filters.visit_to   = isoDate(now)
  }
  fetchFiltered()
}
const activePreset = computed(() => {
  const f = placesStore.filters
  if (!f.visit_from || !f.visit_to) return ''
  const now = new Date()
  const y = now.getFullYear()
  if (f.visit_from === `${y}-01-01` && f.visit_to === `${y}-12-31`) return 'this-year'
  if (f.visit_from === `${y - 1}-01-01` && f.visit_to === `${y - 1}-12-31`) return 'prev-year'
  // last-30-days меряется по разности в днях между to и сегодня
  const to = new Date(f.visit_to)
  const from = new Date(f.visit_from)
  const dayMs = 86400000
  if (Math.abs(now - to) < dayMs && Math.round((to - from) / dayMs) === 30) return 'last-30-days'
  return ''
})
const dateRangeChip = computed(() => {
  const f = placesStore.filters
  if (!f.visit_from && !f.visit_to) return ''
  const fmt = s => s ? s.replace(/-/g, '.').slice(2) : '…' // 26.12.05
  return `${fmt(f.visit_from)}–${fmt(f.visit_to)}`
})
function clearDateRange() {
  placesStore.filters.visit_from = ''
  placesStore.filters.visit_to = ''
  fetchFiltered()
}

const showShelves = computed(() => !placesStore.filters.search && activeFilterCount.value === 0)

// shelves derived from currently-loaded places
const topGems = computed(() => placesStore.places.filter((p) => p.is_gem_place).slice(0, 5))
const cityShelf = computed(() => {
  const map = new Map()
  for (const p of placesStore.places) {
    if (!p.city) continue
    map.set(p.city, (map.get(p.city) || 0) + 1)
  }
  return [...map.entries()].sort((a, b) => b[1] - a[1]).slice(0, 8).map(([name, count]) => ({ name, count }))
})
const cuisineShelf = computed(() => {
  const map = new Map()
  for (const p of placesStore.places) {
    if (!p.cuisine_type_id) continue
    const k = p.cuisine_type_id
    if (!map.has(k)) map.set(k, { id: k, name: p.cuisine_type, count: 0 })
    map.get(k).count++
  }
  return [...map.values()].sort((a, b) => b.count - a.count).slice(0, 8)
})

// Tilt + tape + placeholder helpers
const tilts = ['t-l3', 't-r2', 't-l2', 't-r3', 't-l1']
function tiltFor(id) { return tilts[(id ?? 0) % tilts.length] }
const tapeVariants = ['', 'rose', 'mint', 'blue']
function tapeVariantFor(id) { return tapeVariants[(id ?? 0) % tapeVariants.length] }
const placeholders = ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach', 'sb-photo-brick', 'sb-photo-cream', 'sb-photo-slate', 'sb-photo-indigo']
function placeholderFor(id) { return placeholders[(id ?? 0) % placeholders.length] }

function cuisineName(id) { return catalogs.cuisineTypes.find((c) => c.id === id)?.name || '' }
function categoryName(id) { return catalogs.categories.find((c) => c.id === id)?.name || '' }
function sortLabel(s) {
  if (s.startsWith('rating_user:')) {
    const id = parseInt(s.slice('rating_user:'.length), 10)
    return `по оценке ${friendName(id) || '…'}`
  }
  return ({ rating: 'рейтинг ↓', rating_asc: 'рейтинг ↑', name: 'алфавит' }[s] || 'сортировка')
}

function removeCuisine(id) {
  placesStore.filters.cuisine_type_ids = (placesStore.filters.cuisine_type_ids || []).filter((x) => x !== id)
  fetchFiltered()
}
function removeCategory(id) {
  placesStore.filters.category_ids = (placesStore.filters.category_ids || []).filter((x) => x !== id)
  fetchFiltered()
}

function clearSearch() {
  placesStore.filters.search = ''
  fetchFiltered()
}

function filterByCuisine(id) {
  const arr = placesStore.filters.cuisine_type_ids || []
  if (!arr.includes(id)) {
    placesStore.filters.cuisine_type_ids = [...arr, id]
  }
  fetchFiltered()
}

async function rollDice() {
  // B5 — серверный /api/random с теми же фильтрами + exclude_visited_by=me.
  // На сервере без подходящих 404 → fallback на клиентский рандом из текущей
  // выдачи (чтобы кнопка не висела в пустой ленте).
  const params = {}
  const f = placesStore.filters
  if (f.city) params.city = f.city
  if (f.cuisine_type_ids?.length) params.cuisine_type_id = f.cuisine_type_ids.join(',')
  if (f.is_gem) params.is_gem = 'true'
  params.exclude_visited_by = 'me'
  try {
    const { data } = await http.get('/random', { params })
    if (data?.id) {
      router.push(`/places/${data.id}`)
      return
    }
  } catch (_) {
    // 404 / network error → fallback ниже. Сознательно глушим — кнопка
    // должна работать даже когда сервер ещё не подкрутили.
  }
  // Fallback — клиентский рандом из текущего списка (без exclude_visited_by).
  const list = placesStore.places
  if (!list.length) return
  const p = list[Math.floor(Math.random() * list.length)]
  router.push(`/places/${p.id}`)
}

function resetFilters() {
  placesStore.filters.search = ''
  placesStore.filters.city = ''
  placesStore.filters.cuisine_type_ids = []
  placesStore.filters.category_ids = []
  placesStore.filters.sort = ''
  placesStore.filters.is_gem = false
  placesStore.filters.min_rating = ''
  placesStore.filters.attended_by = []
  placesStore.filters.visit_from = ''
  placesStore.filters.visit_to = ''
  ratingUserId.value = ''
  fetchFiltered()
}

function loadFiltersFromURL() {
  const q = route.query
  placesStore.filters.search = q.search || ''
  placesStore.filters.city = q.city || ''
  placesStore.filters.cuisine_type_ids = q.cuisine_type_ids ? q.cuisine_type_ids.split(',').map(Number) : []
  placesStore.filters.category_ids = q.category_ids ? q.category_ids.split(',').map(Number) : []
  placesStore.filters.sort = q.sort || ''
  placesStore.filters.is_gem = q.is_gem === 'true'
  placesStore.filters.min_rating = q.min_rating || ''
  placesStore.filters.attended_by = q.attended_by ? q.attended_by.split(',').map(Number) : []
  placesStore.filters.visit_from = q.visit_from || ''
  placesStore.filters.visit_to = q.visit_to || ''
  if (placesStore.filters.sort?.startsWith('rating_user:')) {
    ratingUserId.value = placesStore.filters.sort.slice('rating_user:'.length)
  }
  if (q.page) placesStore.page = parseInt(q.page) || 1
}

function syncFiltersToURL() {
  const params = {}
  Object.entries(placesStore.filters).forEach(([key, val]) => {
    if (Array.isArray(val)) {
      if (val.length > 0) params[key] = val.join(',')
    } else if (val !== '' && val !== false && val !== null) {
      params[key] = String(val)
    }
  })
  if (placesStore.page > 1) params.page = String(placesStore.page)
  router.replace({ query: params })
}

const totalPages = computed(() => Math.ceil(placesStore.total / placesStore.limit) || 1)

function goToPage(p) {
  if (p < 1 || p > totalPages.value) return
  placesStore.page = p
  syncFiltersToURL()
  placesStore.fetchPlaces()
}

let debounceTimer = null
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => { placesStore.page = 1; syncFiltersToURL(); placesStore.fetchPlaces() }, 300)
}

function fetchFiltered() {
  placesStore.page = 1
  syncFiltersToURL()
  placesStore.fetchPlaces()
}

async function fetchCities() {
  try {
    const { data } = await http.get('/places/cities')
    cities.value = data || []
  } catch { /* ignore */ }
}

async function fetchFriends() {
  try {
    const { data } = await http.get('/users')
    friends.value = (data || []).filter(u => u.review_count > 0)
  } catch { /* ignore */ }
}

onMounted(async () => {
  loadFiltersFromURL()
  await catalogs.fetchAll()
  await fetchCities()
  await fetchFriends()
  await placesStore.fetchPlaces()
})
</script>

<style scoped lang="scss">
.find {
  // Keep `.sb-screen` bottom padding (BottomTabBar gap) — split shorthand.
  padding-top: calc(18px + var(--aeva-safe-top, 0px));
  padding-inline: 16px;
}

.find-header {
  margin-bottom: 12px;
}

.find-search {
  display: flex;
  gap: 10px;
  align-items: center;
  margin-bottom: 10px;
}

.search-stamp {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  background: oklch(0.94 0.05 85 / 0.6);
  border: 1.4px solid var(--sb-terracotta);
  border-radius: 4px;
  padding: 6px 10px;
  position: relative;
  box-shadow: inset 0 0 0 0.5px rgba(140, 60, 30, 0.2);
}
.search-glyph {
  font-family: var(--sb-hand);
  font-size: 18px;
  color: var(--sb-terracotta);
  line-height: 1;
}
.search-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 16px;
  color: var(--sb-ink);
  min-height: 28px;
  padding: 2px 0;

  &::placeholder { color: var(--sb-ink-mute); font-style: italic; }
}
.clear-x {
  background: transparent;
  border: none;
  font-size: 18px;
  color: var(--sb-ink-mute);
  cursor: pointer;
  line-height: 1;
}

.filter-pin {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: oklch(0.93 0.04 85);
  border: none;
  border-radius: 999px;
  padding: 6px 12px 6px 8px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink);
  cursor: pointer;
  box-shadow: 0 1px 1px rgba(40,30,20,0.1), 0 3px 8px rgba(40,30,20,0.14);
  min-height: 36px;
  position: relative;

  .head {
    width: 14px;
    height: 14px;
    border-radius: 50%;
    background: radial-gradient(circle at 35% 30%, oklch(0.7 0.16 25), oklch(0.42 0.18 25));
    box-shadow: inset 0 -1px 1px rgba(0,0,0,0.3), inset 0 1px 1px rgba(255,255,255,0.3);
  }
  .count {
    background: var(--sb-terracotta);
    color: var(--sb-on-accent);
    font-family: var(--sb-serif);
    font-style: normal;
    font-size: 10px;
    font-weight: 600;
    border-radius: 999px;
    padding: 2px 6px;
    line-height: 1;
  }
  &.on { background: oklch(0.92 0.07 25); color: var(--sb-terracotta); }
}

.find-actions {
  display: flex;
  gap: 14px;
  align-items: center;
  margin-bottom: 12px;
  flex-wrap: wrap;
}
.dice {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--sb-paper-card);
  border: none;
  border-radius: 999px;
  padding: 8px 14px;
  font-family: var(--sb-serif);
  font-size: 14px;
  color: var(--sb-ink);
  cursor: pointer;
  box-shadow: 0 1px 1px rgba(40,30,20,0.06), 0 4px 10px rgba(40,30,20,0.10);
  min-height: 36px;
  &:disabled { opacity: 0.5; cursor: not-allowed; }
  .glyph { font-size: 16px; }
}
.add-link {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink);
  text-decoration: none;
  background: transparent;
  border: none;
  cursor: pointer;
  &:hover { color: var(--sb-terracotta); }
}

.find-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 14px;
}
.chip {
  background: oklch(0.94 0.05 85 / 0.7);
  border: 1px solid rgba(40, 30, 20, 0.18);
  border-radius: 999px;
  padding: 4px 10px;
  font-family: var(--sb-serif);
  font-size: 12px;
  color: var(--sb-ink);
  cursor: pointer;
  &:hover { background: oklch(0.92 0.07 25); }
  &.reset { background: transparent; color: var(--sb-ink-mute); border-style: dashed; }
}

.shelf {
  margin: 0 -16px 18px;
  padding: 4px 16px 0;
}
.shelf-head {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 8px;

  h3 {
    font-family: var(--sb-serif);
    font-style: italic;
    font-weight: 500;
    font-size: 19px;
    color: var(--sb-ink);
    margin: 0;
  }
  .shelf-all {
    font-family: var(--sb-hand);
    font-size: 16px;
    color: var(--sb-ink-mute);
    background: transparent;
    border: none;
    cursor: pointer;
    margin-left: auto;
    &:hover { color: var(--sb-terracotta); }
  }
  .shelf-soon {
    font-family: var(--sb-hand);
    font-size: 14px;
    color: var(--sb-ink-mute);
    margin-left: auto;
  }
}
.shelf-row {
  display: flex;
  gap: 18px;
  overflow-x: auto;
  // Generous top/bottom padding so tape (top: -9px on each polaroid) and
  // shadow (≈8px below) aren't visually clipped by the scroll container.
  padding: 16px 4px 22px;
  scroll-snap-type: x mandatory;
  -webkit-overflow-scrolling: touch;

  // hide scrollbar
  scrollbar-width: none;
  &::-webkit-scrollbar { display: none; }
}
.shelf-gem {
  text-decoration: none;
  color: inherit;
  scroll-snap-align: start;
  flex-shrink: 0;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 130px;

  .gem-corner {
    position: absolute;
    top: 6px;
    right: 6px;
    z-index: 2;
  }
  .shelf-cap {
    margin-top: 6px;
    font-family: var(--sb-serif);
    font-size: 13px;
    color: var(--sb-ink);
    text-align: center;
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.shelf-stamps {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 6px;
}

/* Города — вертикальный список (F1) */
.city-list {
  list-style: none;
  margin: 0;
  padding: 0;
}
.city-list li {
  border-bottom: 1px dashed rgba(40, 30, 20, 0.14);
}
.city-list li:last-child { border-bottom: none; }
.city-row {
  display: flex;
  align-items: baseline;
  width: 100%;
  background: transparent;
  border: none;
  padding: 10px 4px;
  font-family: var(--sb-serif);
  cursor: pointer;
  text-align: left;
  text-decoration: none;
  color: inherit;

  &:hover .city-name { color: var(--sb-terracotta); }
}
.city-name {
  flex: 1;
  font-style: italic;
  font-weight: 500;
  font-size: 18px;
  color: var(--sb-ink);
  line-height: 1.1;
}
.city-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  font-family: var(--sb-serif);
  font-weight: 600;
  font-size: 11px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--sb-ink);
  border: 1.4px solid var(--sb-ink);
  border-radius: 2px;
  padding: 4px 9px 3px;
  box-shadow: inset 0 0 0 0.5px rgba(40, 30, 20, 0.18);
}

:deep(.sb-stamp.clickable) { cursor: pointer; }
:deep(.sb-stamp.clickable:hover) { background: oklch(0.92 0.05 145 / 0.4); }

/* Friends shelf — горизонтальная карусель аватарок (DESIGN-DECISIONS §F1).
   60px кружок, имя серифой и count рукописным под ним. */
.shelf-row.friends {
  gap: 18px;
  padding: 10px 4px 14px;
}
.friend-tile {
  flex: 0 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  text-decoration: none;
  color: inherit;
  width: 76px;
}
.friend-avatar {
  position: relative;
  width: 60px;
  height: 60px;
  font-size: 22px;
  box-shadow: 0 0 0 2px var(--sb-paper-card), 0 2px 6px rgba(40, 30, 20, 0.18);
}
.friend-avatar.has-photo { background: var(--sb-paper-card); overflow: hidden; }
.friend-avatar .r-ph { display: block; width: 100%; height: 100%; object-fit: cover; }
.friend-name {
  font-family: var(--sb-serif);
  font-size: 13px;
  color: var(--sb-ink);
  text-align: center;
  word-break: break-word;
  line-height: 1.1;
}
.friend-count {
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
  line-height: 1;
}
.friend-tile:hover .friend-name { color: var(--sb-terracotta); }

.find-results {
  margin-top: 4px;
}
.results-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.find-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 14px;
  margin: 18px 0 8px;
  font-family: var(--sb-serif);
  color: var(--sb-ink-soft);
}
.page-link {
  background: transparent;
  border: none;
  font-family: var(--sb-serif);
  font-size: 22px;
  color: var(--sb-ink);
  cursor: pointer;
  padding: 4px 12px;
  &:disabled { color: var(--sb-ink-mute); cursor: not-allowed; }
}
.page-meta {
  font-family: var(--sb-hand);
  font-size: 16px;
}

/* Drawer paper override */
.find-drawer {
  background: var(--sb-paper) !important;
  font-family: var(--sb-serif);

  .offcanvas-header {
    border-bottom: 1px dashed rgba(40, 30, 20, 0.3);
    .offcanvas-title {
      font-family: var(--sb-serif);
      font-style: italic;
      font-weight: 500;
    }
  }
}
.drawer-label {
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  margin-bottom: 4px;
  display: block;
}
.btn-apply {
  background: var(--sb-terracotta);
  color: var(--sb-on-accent);
  border: none;
  border-radius: 999px;
  padding: 10px 18px;
  font-family: var(--sb-serif);
  font-style: italic;
  &:hover { background: oklch(0.55 0.14 30); color: var(--sb-on-accent); }
}
.reset-btn {
  font-family: var(--sb-serif);
  font-style: italic;
  color: var(--sb-ink-mute);
  text-decoration: none;
}

/* Q5 — «кто был»: лента аватарок-чипов с многоточным выбором. */
.att-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.att-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px 4px 4px;
  background: transparent;
  border: 1.4px dashed rgba(40, 30, 20, 0.18);
  border-radius: 999px;
  cursor: pointer;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink);

  &.on {
    border-style: solid;
    border-color: var(--sb-terracotta);
    background: oklch(0.92 0.07 25);
    color: var(--sb-terracotta);
  }

  .r-tag {
    position: relative;
    width: 24px;
    height: 24px;
    font-size: 11px;
  }
  .r-tag.has-photo { overflow: hidden; }
  .r-ph { width: 100%; height: 100%; object-fit: cover; display: block; }
  .att-name {
    line-height: 1;
  }
}

/* Q5 — «когда»: бумажные date-инпуты + preset-чипы под ними. */
.date-row {
  display: flex;
  align-items: center;
  gap: 8px;
  .form-control { flex: 1 1 0; }
}
.date-sep {
  font-family: var(--sb-hand);
  font-size: 18px;
  color: var(--sb-ink-mute);
}
.preset-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}
.preset-chip {
  background: transparent;
  border: 1.4px dashed rgba(40, 30, 20, 0.22);
  border-radius: 999px;
  padding: 4px 12px;
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-soft);
  cursor: pointer;

  &.on {
    border-style: solid;
    border-color: var(--sb-terracotta);
    color: var(--sb-terracotta);
    background: oklch(0.94 0.05 30);
  }
}
</style>
