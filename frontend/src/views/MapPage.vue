<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h2>Карта заведений</h2>
      <router-link to="/places" class="btn btn-outline-secondary btn-sm">← Список</router-link>
    </div>

    <!-- Desktop filter row (md+) -->
    <div class="card mb-3 d-none d-md-block">
      <div class="card-body py-2">
        <div class="row g-2">
          <div class="col-md-3">
            <input
              v-model="placesStore.filters.search"
              type="text"
              class="form-control form-control-sm"
              placeholder="🔍 Поиск"
              @input="debouncedFetch"
            />
          </div>
          <div class="col-md-2">
            <select v-model="placesStore.filters.city" class="form-select form-select-sm" @change="fetchFiltered">
              <option value="">Все города</option>
              <option v-for="city in cities" :key="city" :value="city">{{ city }}</option>
            </select>
          </div>
          <div class="col-md-2">
            <MultiSelect
              :modelValue="placesStore.filters.cuisine_type_ids"
              @update:model-value="v => { placesStore.filters.cuisine_type_ids = v; fetchFiltered() }"
              :options="catalogs.cuisineTypes"
              placeholder="Кухни"
            />
          </div>
          <div class="col-md-2">
            <MultiSelect
              :modelValue="placesStore.filters.category_ids"
              @update:model-value="v => { placesStore.filters.category_ids = v; fetchFiltered() }"
              :options="catalogs.categories"
              placeholder="Категории"
            />
          </div>
          <div class="col-md-1 d-flex align-items-center">
            <div class="form-check">
              <input v-model="placesStore.filters.is_gem" class="form-check-input" type="checkbox" @change="fetchFiltered" />
              <label class="form-check-label">💎</label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile filter bar -->
    <div class="d-md-none mb-2">
      <div class="d-flex gap-2">
        <input
          v-model="placesStore.filters.search"
          type="text"
          class="form-control flex-grow-1"
          placeholder="🔍 Поиск"
          @input="debouncedFetch"
        />
        <button
          class="btn btn-outline-secondary position-relative px-3"
          type="button"
          data-bs-toggle="offcanvas"
          data-bs-target="#mapFilterDrawer"
          aria-label="Открыть фильтры"
        >
          🎛
          <span
            v-if="activeFilterCount > 0"
            class="position-absolute top-0 start-100 translate-middle badge rounded-pill bg-primary"
            style="font-size: 0.65rem"
          >{{ activeFilterCount }}</span>
        </button>
      </div>
    </div>

    <!-- Mobile filter drawer -->
    <div
      class="offcanvas offcanvas-bottom d-md-none"
      tabindex="-1"
      id="mapFilterDrawer"
      aria-labelledby="mapFilterDrawerLabel"
      style="height: auto; max-height: 85vh"
    >
      <div class="offcanvas-header border-bottom">
        <h5 class="offcanvas-title" id="mapFilterDrawerLabel">Фильтры</h5>
        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Закрыть"></button>
      </div>
      <div class="offcanvas-body">
        <div class="mb-3">
          <label class="form-label small text-muted mb-1">Город</label>
          <select v-model="placesStore.filters.city" class="form-select">
            <option value="">Все города</option>
            <option v-for="city in cities" :key="city" :value="city">{{ city }}</option>
          </select>
        </div>
        <div class="mb-3">
          <label class="form-label small text-muted mb-1">Кухни</label>
          <MultiSelect
            :modelValue="placesStore.filters.cuisine_type_ids"
            @update:model-value="v => placesStore.filters.cuisine_type_ids = v"
            :options="catalogs.cuisineTypes"
            placeholder="Любые"
          />
        </div>
        <div class="mb-3">
          <label class="form-label small text-muted mb-1">Категории</label>
          <MultiSelect
            :modelValue="placesStore.filters.category_ids"
            @update:model-value="v => placesStore.filters.category_ids = v"
            :options="catalogs.categories"
            placeholder="Любые"
          />
        </div>
        <div class="form-check form-switch mb-3">
          <input
            v-model="placesStore.filters.is_gem"
            class="form-check-input"
            type="checkbox"
            id="gemFilterMapMobile"
            role="switch"
          />
          <label class="form-check-label" for="gemFilterMapMobile">💎 Только жемчужины</label>
        </div>
        <div class="d-grid gap-2 mt-4">
          <button class="btn btn-primary" data-bs-dismiss="offcanvas" @click="fetchFiltered">Применить</button>
          <button class="btn btn-link" @click="resetFilters">Сбросить</button>
        </div>
      </div>
    </div>

    <MapView
      :places="placesWithCoords"
      :height="mapHeight"
      @marker-click="onMarkerClick"
    />

    <div class="mt-2 text-muted small">
      {{ placesWithCoords.length }} из {{ placesStore.places.length }} заведений на карте
      <span v-if="placesStore.places.length > placesWithCoords.length">
        (у {{ placesStore.places.length - placesWithCoords.length }} нет координат)
      </span>
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

const placesWithCoords = computed(() =>
  placesStore.places.filter(p => p.lat && p.lng)
)

// Map height: tall on mobile (use most of viewport minus header/filter/tab-bar);
// fixed-ish on desktop (legacy behavior).
const isMobile = ref(false)
function checkMobile() { isMobile.value = window.innerWidth < 768 }
onMounted(checkMobile)
window.addEventListener('resize', checkMobile)
onBeforeUnmount(() => window.removeEventListener('resize', checkMobile))

const mapHeight = computed(() =>
  isMobile.value
    // Mobile: top navbar (~56) + h2 (~40) + filter row (~52) + bottom tab bar (~56 + safe area) ≈ 220
    ? 'calc(100vh - 220px)'
    : 'calc(100vh - 280px)'
)

const activeFilterCount = computed(() => {
  const f = placesStore.filters
  let n = 0
  if (f.city) n++
  if (f.cuisine_type_ids?.length) n++
  if (f.category_ids?.length) n++
  if (f.is_gem) n++
  return n
})

function resetFilters() {
  placesStore.filters.city = ''
  placesStore.filters.cuisine_type_ids = []
  placesStore.filters.category_ids = []
  placesStore.filters.is_gem = false
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

onMounted(async () => {
  await catalogs.fetchAll()
  await fetchCities()
  await placesStore.fetchAllPlaces()
})
</script>
