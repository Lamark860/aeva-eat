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
      <button v-if="placesStore.filters.sort" class="chip" @click="placesStore.filters.sort = ''; fetchFiltered()">
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
          <button class="shelf-all" @click="placesStore.filters.is_gem = true; fetchFiltered()">все →</button>
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
            <button class="city-row" @click="filterByCity(c.name)">
              <span class="city-name">{{ c.name }}</span>
              <span class="city-count">{{ c.count }}</span>
            </button>
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

      <!-- friends shelf — placeholder, awaits backend -->
      <section class="shelf shelf-stub">
        <div class="shelf-head">
          <h3>По друзьям</h3>
          <span class="shelf-soon">скоро</span>
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
          <select v-model="placesStore.filters.sort" class="form-select">
            <option value="">сначала новые</option>
            <option value="rating">по рейтингу ↓</option>
            <option value="rating_asc">по рейтингу ↑</option>
            <option value="name">по названию</option>
          </select>
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

const route = useRoute()
const router = useRouter()
const placesStore = usePlacesStore()
const catalogs = useCatalogsStore()
const auth = useAuthStore()

const cities = ref([])

const activeFilterCount = computed(() => {
  const f = placesStore.filters
  let n = 0
  if (f.city) n++
  if (f.cuisine_type_ids?.length) n++
  if (f.category_ids?.length) n++
  if (f.sort) n++
  if (f.is_gem) n++
  if (f.min_rating) n++
  return n
})

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

function filterByCity(name) {
  placesStore.filters.city = name
  fetchFiltered()
}
function filterByCuisine(id) {
  const arr = placesStore.filters.cuisine_type_ids || []
  if (!arr.includes(id)) {
    placesStore.filters.cuisine_type_ids = [...arr, id]
  }
  fetchFiltered()
}

function rollDice() {
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

onMounted(async () => {
  loadFiltersFromURL()
  await catalogs.fetchAll()
  await fetchCities()
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
    color: #fff;
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
  background: #fdfcf7;
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

.shelf-stub {
  opacity: 0.55;
}

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
  color: #fff;
  border: none;
  border-radius: 999px;
  padding: 10px 18px;
  font-family: var(--sb-serif);
  font-style: italic;
  &:hover { background: oklch(0.55 0.14 30); color: #fff; }
}
.reset-btn {
  font-family: var(--sb-serif);
  font-style: italic;
  color: var(--sb-ink-mute);
  text-decoration: none;
}
</style>
