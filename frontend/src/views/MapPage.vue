<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h2>Карта заведений</h2>
      <router-link to="/places" class="btn btn-outline-secondary btn-sm">← Список</router-link>
    </div>

    <!-- Filters -->
    <div class="card mb-3">
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
            <select v-model="placesStore.filters.cuisine_type_id" class="form-select form-select-sm" @change="fetchFiltered">
              <option value="">Все кухни</option>
              <option v-for="ct in catalogs.cuisineTypes" :key="ct.id" :value="ct.id">{{ ct.name }}</option>
            </select>
          </div>
          <div class="col-md-2">
            <select v-model="placesStore.filters.category_id" class="form-select form-select-sm" @change="fetchFiltered">
              <option value="">Все категории</option>
              <option v-for="cat in catalogs.categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
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

    <MapView
      :places="placesWithCoords"
      :height="'calc(100vh - 280px)'"
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
import { onMounted, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { usePlacesStore } from '../stores/places'
import { useCatalogsStore } from '../stores/catalogs'
import MapView from '../components/MapView.vue'
import http from '../api/http'

const router = useRouter()
const placesStore = usePlacesStore()
const catalogs = useCatalogsStore()
const cities = ref([])

const placesWithCoords = computed(() =>
  placesStore.places.filter(p => p.lat && p.lng)
)

let debounceTimer = null
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => placesStore.fetchPlaces(), 300)
}

function fetchFiltered() {
  placesStore.fetchPlaces()
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
  await placesStore.fetchPlaces()
})
</script>
