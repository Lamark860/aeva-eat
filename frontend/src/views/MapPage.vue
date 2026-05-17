<template>
  <div class="sb-paper sb-grain sb-screen map-page">
    <header class="map-head">
      <div class="sb-section-head" style="padding: 0">
        <h2>Карта</h2>
        <span class="sub">точки на бумаге</span>
      </div>
    </header>

    <div class="map-search">
      <div class="search-stamp">
        <span class="search-glyph">⌕</span>
        <input
          v-model="placesStore.filters.search"
          type="text"
          class="search-input"
          placeholder="место, кухня, город…"
          @input="debouncedFetch"
        />
        <button v-if="placesStore.filters.search" class="clear-x" aria-label="Очистить" @click="clearSearch">×</button>
      </div>
      <button
        class="filter-pin"
        :class="{ on: activeFilterCount > 0 }"
        type="button"
        data-bs-toggle="offcanvas"
        data-bs-target="#mapFilterDrawer"
        aria-label="Фильтры"
      >
        <span class="head" aria-hidden="true"></span>
        <span class="lbl">фильтры</span>
        <span v-if="activeFilterCount > 0" class="count">{{ activeFilterCount }}</span>
      </button>
    </div>

    <MapView
      class="map-canvas"
      :places="placesWithCoords"
      :height="mapHeight"
      @marker-click="onMarkerClick"
    />

    <div class="map-stats">
      {{ placesWithCoords.length }} из {{ placesStore.places.length }} мест на карте
      <span v-if="placesStore.places.length > placesWithCoords.length">
        · у {{ placesStore.places.length - placesWithCoords.length }} нет координат
      </span>
    </div>

    <!-- Filter drawer (offcanvas, paper-styled) -->
    <div
      class="offcanvas offcanvas-bottom find-drawer"
      tabindex="-1"
      id="mapFilterDrawer"
      aria-labelledby="mapFilterDrawerLabel"
      style="height: auto; max-height: 85vh"
    >
      <div class="offcanvas-header">
        <h5 class="offcanvas-title" id="mapFilterDrawerLabel">Фильтры</h5>
        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Закрыть"></button>
      </div>
      <div class="offcanvas-body">
        <div class="mb-3">
          <label class="drawer-label">Город</label>
          <select v-model="draft.city" class="form-select">
            <option value="">все</option>
            <option v-for="city in cities" :key="city" :value="city">{{ city }}</option>
          </select>
        </div>
        <div class="mb-3">
          <label class="drawer-label">Кухни</label>
          <MultiSelect
            :modelValue="draft.cuisine_type_ids"
            @update:model-value="v => draft.cuisine_type_ids = v"
            :options="catalogs.cuisineTypes"
            placeholder="любые"
          />
        </div>
        <div class="mb-3">
          <label class="drawer-label">Категории</label>
          <MultiSelect
            :modelValue="draft.category_ids"
            @update:model-value="v => draft.category_ids = v"
            :options="catalogs.categories"
            placeholder="любые"
          />
        </div>
        <div class="form-check form-switch mb-3">
          <input v-model="draft.is_gem" class="form-check-input" type="checkbox" id="gemFilterMap" role="switch" />
          <label class="form-check-label" for="gemFilterMap">только&nbsp;жемчужины&nbsp;♦</label>
        </div>
        <div class="d-grid gap-2 mt-4">
          <button class="btn btn-apply" data-bs-dismiss="offcanvas" @click="applyDraft">применить</button>
          <button class="btn btn-link reset-btn" @click="resetDraft">сбросить</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { usePlacesStore } from '../stores/places'
import { useCatalogsStore } from '../stores/catalogs'
import MapView from '../components/MapView.vue'
import MultiSelect from '../components/MultiSelect.vue'
import http from '../api/http'

const router = useRouter()
const placesStore = usePlacesStore()
const catalogs = useCatalogsStore()
const cities = ref([])

const placesWithCoords = computed(() => placesStore.places.filter((p) => p.lat && p.lng))

const isMobile = ref(false)
function checkMobile() { isMobile.value = window.innerWidth < 768 }
onMounted(checkMobile)
window.addEventListener('resize', checkMobile)
onBeforeUnmount(() => window.removeEventListener('resize', checkMobile))

// Reserve space for: top wordmark area (~84) + search row (~58) + tabbar+safe (~80)
const mapHeight = computed(() => (isMobile.value ? 'calc(100vh - 240px)' : 'calc(100vh - 260px)'))

const activeFilterCount = computed(() => {
  const f = placesStore.filters
  let n = 0
  if (f.city) n++
  if (f.cuisine_type_ids?.length) n++
  if (f.category_ids?.length) n++
  if (f.is_gem) n++
  return n
})

function clearSearch() {
  placesStore.filters.search = ''
  fetchFiltered()
}

let debounceTimer = null
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => placesStore.fetchAllPlaces(), 300)
}

function fetchFiltered() {
  placesStore.fetchAllPlaces()
}

function onMarkerClick(place) {
  router.push(`/places/${place.id}`)
}

async function fetchCities() {
  try {
    const { data } = await http.get('/places/cities')
    cities.value = data || []
  } catch { /* ignore */ }
}

// ---- Draft state для drawer (как в Places.vue) ----
// Drawer пишет в локальный draft, apply копирует в store + fetch. Иначе
// каждый тап чекбокса в drawer'е дёргал бы карту через watch-deep.
const draft = ref({ city: '', cuisine_type_ids: [], category_ids: [], is_gem: false })
function seedDraft() {
  const f = placesStore.filters
  draft.value = {
    city: f.city || '',
    cuisine_type_ids: [...(f.cuisine_type_ids || [])],
    category_ids: [...(f.category_ids || [])],
    is_gem: !!f.is_gem,
  }
}
function applyDraft() {
  const f = placesStore.filters
  const d = draft.value
  f.city = d.city
  f.cuisine_type_ids = [...d.cuisine_type_ids]
  f.category_ids = [...d.category_ids]
  f.is_gem = d.is_gem
  fetchFiltered()
}
function resetDraft() {
  draft.value = { city: '', cuisine_type_ids: [], category_ids: [], is_gem: false }
}
function onDrawerHidden() {
  if (!document.querySelector('.offcanvas.show')) {
    document.querySelectorAll('.offcanvas-backdrop').forEach(b => b.remove())
    document.body.style.overflow = ''
    document.body.style.paddingRight = ''
  }
}

onMounted(async () => {
  await catalogs.fetchAll()
  await fetchCities()
  await placesStore.fetchAllPlaces()
  const drawer = document.getElementById('mapFilterDrawer')
  if (drawer) {
    drawer.addEventListener('show.bs.offcanvas', seedDraft)
    drawer.addEventListener('hidden.bs.offcanvas', onDrawerHidden)
  }
})
</script>

<style scoped lang="scss">
.map-page {
  padding-top: calc(18px + var(--aeva-safe-top, 0px));
  padding-inline: 16px;
}

.map-head { margin-bottom: 12px; }

.map-search {
  display: flex;
  gap: 10px;
  align-items: center;
  margin-bottom: 12px;
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
  box-shadow: 0 1px 1px rgba(40, 30, 20, 0.1), 0 3px 8px rgba(40, 30, 20, 0.14);
  min-height: 36px;

  .head {
    width: 14px;
    height: 14px;
    border-radius: 50%;
    background: radial-gradient(circle at 35% 30%, oklch(0.7 0.16 25), oklch(0.42 0.18 25));
    box-shadow: inset 0 -1px 1px rgba(0, 0, 0, 0.3), inset 0 1px 1px rgba(255, 255, 255, 0.3);
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

.map-canvas { display: block; }

.map-stats {
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
  margin: 8px 0 10px;
  text-align: center;
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
</style>
