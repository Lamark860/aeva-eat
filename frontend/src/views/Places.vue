<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2>Заведения</h2>
      <router-link v-if="auth.isAuthenticated" to="/places/new" class="btn btn-primary">
        + Добавить
      </router-link>
    </div>

    <!-- Desktop filter row (md+) -->
    <div class="card mb-4 d-none d-md-block">
      <div class="card-body">
        <div class="row g-2">
          <div class="col-md-3">
            <input
              v-model="placesStore.filters.search"
              type="text"
              class="form-control form-control-sm"
              placeholder="🔍 Поиск по названию"
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
          <div class="col-md-2">
            <select v-model="placesStore.filters.sort" class="form-select form-select-sm" @change="fetchFiltered">
              <option value="">Сначала новые</option>
              <option value="rating">По рейтингу ↓</option>
              <option value="rating_asc">По рейтингу ↑</option>
              <option value="name">По названию</option>
            </select>
          </div>
          <div class="col-md-1 d-flex align-items-center gap-2">
            <div class="form-check">
              <input
                v-model="placesStore.filters.is_gem"
                class="form-check-input"
                type="checkbox"
                id="gemFilter"
                @change="fetchFiltered"
              />
              <label class="form-check-label" for="gemFilter">💎</label>
            </div>
          </div>
        </div>
        <div class="row g-2 mt-1" v-if="placesStore.filters.min_rating">
          <div class="col-auto">
            <span class="badge bg-secondary">Мин. рейтинг: {{ placesStore.filters.min_rating }}+</span>
            <button class="btn btn-sm btn-link p-0 ms-1" @click="placesStore.filters.min_rating = ''; fetchFiltered()">✕</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile filter bar (sticky search + filter button) -->
    <div class="d-md-none mb-3 mobile-filter-bar">
      <div class="d-flex gap-2">
        <input
          v-model="placesStore.filters.search"
          type="text"
          class="form-control flex-grow-1"
          placeholder="🔍 Поиск по названию"
          @input="debouncedFetch"
        />
        <button
          class="btn btn-outline-secondary position-relative px-3"
          type="button"
          data-bs-toggle="offcanvas"
          data-bs-target="#placesFilterDrawer"
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

    <!-- Mobile filter drawer (offcanvas-bottom) -->
    <div
      class="offcanvas offcanvas-bottom d-md-none"
      tabindex="-1"
      id="placesFilterDrawer"
      aria-labelledby="placesFilterDrawerLabel"
      style="height: auto; max-height: 85vh"
    >
      <div class="offcanvas-header border-bottom">
        <h5 class="offcanvas-title" id="placesFilterDrawerLabel">Фильтры</h5>
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
        <div class="mb-3">
          <label class="form-label small text-muted mb-1">Сортировка</label>
          <select v-model="placesStore.filters.sort" class="form-select">
            <option value="">Сначала новые</option>
            <option value="rating">По рейтингу ↓</option>
            <option value="rating_asc">По рейтингу ↑</option>
            <option value="name">По названию</option>
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
          <label class="form-check-label" for="gemFilterMobile">💎 Только жемчужины</label>
        </div>
        <div v-if="placesStore.filters.min_rating" class="mb-3">
          <span class="badge bg-secondary">Мин. рейтинг: {{ placesStore.filters.min_rating }}+</span>
          <button class="btn btn-sm btn-link p-0 ms-1" @click="placesStore.filters.min_rating = ''">✕</button>
        </div>
        <div class="d-grid gap-2 mt-4">
          <button class="btn btn-primary" data-bs-dismiss="offcanvas" @click="fetchFiltered">Применить</button>
          <button class="btn btn-link" @click="resetFilters">Сбросить</button>
        </div>
      </div>
    </div>

    <div v-if="placesStore.loading" class="text-center py-5">
      <div class="spinner-border"></div>
    </div>

    <div v-else-if="placesStore.places.length === 0" class="text-center py-5 text-muted">
      <p>Заведений пока нет. Добавьте первое!</p>
    </div>

    <div v-else class="row g-3">
      <div v-for="place in placesStore.places" :key="place.id" class="col-12 col-md-6 col-lg-4">
        <PlaceCard :place="place" />
      </div>
    </div>

    <!-- Pagination -->
    <nav v-if="totalPages > 1" class="mt-4 d-flex justify-content-center">
      <ul class="pagination">
        <li class="page-item" :class="{ disabled: placesStore.page <= 1 }">
          <a class="page-link" href="#" @click.prevent="goToPage(placesStore.page - 1)">«</a>
        </li>
        <li
          v-for="p in visiblePages"
          :key="p"
          class="page-item"
          :class="{ active: p === placesStore.page }"
        >
          <a class="page-link" href="#" @click.prevent="goToPage(p)">{{ p }}</a>
        </li>
        <li class="page-item" :class="{ disabled: placesStore.page >= totalPages }">
          <a class="page-link" href="#" @click.prevent="goToPage(placesStore.page + 1)">»</a>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePlacesStore } from '../stores/places'
import { useCatalogsStore } from '../stores/catalogs'
import { useAuthStore } from '../stores/auth'
import PlaceCard from '../components/PlaceCard.vue'
import MultiSelect from '../components/MultiSelect.vue'
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

function resetFilters() {
  placesStore.filters.city = ''
  placesStore.filters.cuisine_type_ids = []
  placesStore.filters.category_ids = []
  placesStore.filters.sort = ''
  placesStore.filters.is_gem = false
  placesStore.filters.min_rating = ''
  fetchFiltered()
}

// Sync filters from URL on load
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

const visiblePages = computed(() => {
  const pages = []
  const tp = totalPages.value
  const cp = placesStore.page
  let start = Math.max(1, cp - 2)
  let end = Math.min(tp, cp + 2)
  if (end - start < 4) {
    if (start === 1) end = Math.min(tp, start + 4)
    else start = Math.max(1, end - 4)
  }
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

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
