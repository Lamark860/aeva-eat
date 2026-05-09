<template>
  <div class="location-picker">
    <div class="mb-2 position-relative">
      <input
        v-model="searchQuery"
        type="text"
        class="form-control"
        placeholder="Название заведения или адрес..."
        autocomplete="off"
        @input="onInput"
        @keydown.enter.prevent="pickFirst"
        @keydown.down.prevent="moveDown"
        @keydown.up.prevent="moveUp"
        @blur="onBlur"
        @focus="showDropdown = suggestions.length > 0"
      />
      <ul v-if="showDropdown && suggestions.length" class="suggest-dropdown">
        <li
          v-for="(item, idx) in suggestions"
          :key="idx"
          class="suggest-item"
          :class="{ active: idx === activeIdx }"
          @mousedown.prevent="pickItem(item)"
        >
          <span class="suggest-name">{{ item.name }}</span>
          <span v-if="item.address" class="suggest-desc">{{ item.address }}</span>
          <span v-else-if="item.description" class="suggest-desc">{{ item.description }}</span>
        </li>
      </ul>
      <small class="text-muted d-block mt-1">Начните вводить — появятся заведения и адреса. Или кликните по карте.</small>
      <div v-if="suggestions.length === 0 && showDropdown" class="suggest-dropdown">
        <li class="suggest-item text-muted">Ничего не найдено</li>
      </div>
    </div>

    <div ref="mapEl" class="location-picker-map"></div>

    <div v-if="modelValue.lat && modelValue.lng" class="mt-2 d-flex align-items-center gap-2">
      <span class="badge bg-light text-dark border">
        📍 {{ modelValue.lat.toFixed(5) }}, {{ modelValue.lng.toFixed(5) }}
      </span>
      <span v-if="resolvedAddress" class="text-muted small text-truncate" style="max-width:300px">{{ resolvedAddress }}</span>
      <button type="button" class="btn btn-outline-secondary btn-sm ms-auto" @click="clear">Сбросить</button>
    </div>
  </div>
</template>

<script setup>
/* global ymaps */
import { ref, onMounted, onUnmounted, watch } from 'vue'
import http from '../api/http'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({ lat: null, lng: null })
  },
  city: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'address-found', 'place-found'])

const mapEl = ref(null)
const searchQuery = ref('')
const resolvedAddress = ref('')
const suggestions = ref([])
const showDropdown = ref(false)
const activeIdx = ref(-1)

let map = null
let marker = null
let debounceTimer = null

function setMarker(lat, lng) {
  try {
    if (marker) {
      marker.geometry.setCoordinates([lat, lng])
    } else if (map) {
      marker = new ymaps.Placemark([lat, lng], {}, {
        draggable: true,
        preset: 'islands#redDotIcon'
      })
      marker.events.add('dragend', () => {
        const coords = marker.geometry.getCoordinates()
        emit('update:modelValue', { lat: coords[0], lng: coords[1] })
        reverseGeocode(coords[0], coords[1])
      })
      map.geoObjects.add(marker)
    }
    map?.setCenter([lat, lng], Math.max(map.getZoom(), 15), { duration: 300 })
    emit('update:modelValue', { lat, lng })
  } catch (e) {
    console.error('[LocationPicker] setMarker error:', e)
  }
}

function clear() {
  if (marker && map) {
    map.geoObjects.remove(marker)
    marker = null
  }
  resolvedAddress.value = ''
  searchQuery.value = ''
  suggestions.value = []
  showDropdown.value = false
  emit('update:modelValue', { lat: null, lng: null })
}

function buildQuery(q) {
  const prefix = props.city ? props.city + ', ' : ''
  return q.startsWith(prefix) ? q : prefix + q
}

function onInput() {
  activeIdx.value = -1
  clearTimeout(debounceTimer)
  const q = searchQuery.value.trim()
  if (q.length < 2) {
    suggestions.value = []
    showDropdown.value = false
    return
  }
  debounceTimer = setTimeout(() => fetchSuggestions(q), 400)
}

async function fetchSuggestions(q) {
  try {
    const fullQ = buildQuery(q)
    const center = map ? map.getCenter() : [55.7558, 37.6173]
    // Use the auth-aware http client — /api/suggest is behind JWT middleware,
    // bare fetch() returns 401 silently and the dropdown stays empty.
    const { data } = await http.get('/suggest', {
      params: { text: fullQ, ll: `${center[1]},${center[0]}` },
    })
    const items = (data.results || []).map(r => {
      const name = r.title?.text || ''
      const sub = r.subtitle?.text || ''
      // subtitle format: "Бар · улица Карла Маркса, 5" or just "улица Карла Маркса, 5"
      const addressPart = sub.includes('·') ? sub.split('·').slice(1).join('·').trim() : sub
      return {
        displayName: name + (sub ? ', ' + sub : ''),
        value: name + (addressPart ? ', ' + addressPart : ''),
        name,
        description: sub,
        address: addressPart,
        tags: r.tags || [],
        uri: r.uri || null
      }
    })
    suggestions.value = items
    showDropdown.value = items.length > 0
  } catch {
    suggestions.value = []
    showDropdown.value = false
  }
}

function pickItem(item) {
  searchQuery.value = item.displayName || ''
  showDropdown.value = false
  suggestions.value = []
  // Build search query with city prefix for ymaps.search
  searchAndPlace(buildQuery(item.value || item.name), item.name || item.value || '')
}

function pickFirst() {
  if (activeIdx.value >= 0 && activeIdx.value < suggestions.value.length) {
    pickItem(suggestions.value[activeIdx.value])
  } else if (suggestions.value.length > 0) {
    pickItem(suggestions.value[0])
  } else {
    searchAndPlace(buildQuery(searchQuery.value))
  }
}

function moveDown() {
  if (suggestions.value.length === 0) return
  activeIdx.value = (activeIdx.value + 1) % suggestions.value.length
}

function moveUp() {
  if (suggestions.value.length === 0) return
  activeIdx.value = activeIdx.value <= 0 ? suggestions.value.length - 1 : activeIdx.value - 1
}

function onBlur() {
  setTimeout(() => { showDropdown.value = false }, 150)
}

function searchAndPlace(query, originalName) {
  if (!query || !map) return
  const searchQ = query
  console.log('[LocationPicker] searchAndPlace:', searchQ)
  ymaps.search(searchQ, {
    results: 1,
    boundedBy: map.getBounds(),
    strictBounds: false
  }).then(res => {
    const obj = res.geoObjects.get(0)
    if (obj) {
      const coords = obj.geometry.getCoordinates()
      const props = obj.properties.getAll()
      const name = props.name || ''
      const description = props.description || ''
      const full = name + (description ? ', ' + description : '')

      setMarker(coords[0], coords[1])
      resolvedAddress.value = full
      emit('address-found', { address: full, lat: coords[0], lng: coords[1] })

      // Parse city from description: "ул. Татарстан, 3/2, Вахитовский район, Казань"
      let city = ''
      if (description) {
        const parts = description.split(', ')
        city = parts[parts.length - 1] || ''
      }

      // Business results have companyMetaData with structured address
      const company = props.companyMetaData
      const isBiz = props.type === 'business'
      const address = isBiz
        ? description.replace(/, [^,]+$/, '') // remove city from end
        : (description || '')

      const placeName = isBiz ? name : (originalName || name)

      emit('place-found', {
        name: placeName,
        city,
        address,
        categories: props.categories || [],
        rating: props.rating || null,
        url: props.url || null,
        lat: coords[0],
        lng: coords[1]
      })
    }
  })
}

async function reverseGeocode(lat, lng) {
  try {
    const result = await ymaps.geocode([lat, lng], { results: 1 })
    const firstGeoObject = result.geoObjects.get(0)
    if (firstGeoObject) {
      const address = firstGeoObject.getAddressLine()
      resolvedAddress.value = address
      emit('address-found', { address, lat, lng })
    }
  } catch { /* ignore */ }
}

function centerOnCity(city) {
  if (!map || !city) return
  ymaps.geocode(city, { results: 1 }).then(res => {
    const obj = res.geoObjects.get(0)
    if (obj) {
      const coords = obj.geometry.getCoordinates()
      if (!marker) {
        map.setCenter(coords, 12, { duration: 500 })
      }
    }
  })
}

onMounted(() => {
  ymaps.ready(() => {
    const initLat = props.modelValue.lat || 55.7558
    const initLng = props.modelValue.lng || 37.6173
    const initZoom = props.modelValue.lat ? 15 : 12

    map = new ymaps.Map(mapEl.value, {
      center: [initLat, initLng],
      zoom: initZoom,
      controls: ['zoomControl']
    })

    if (props.modelValue.lat && props.modelValue.lng) {
      setMarker(props.modelValue.lat, props.modelValue.lng)
    } else if (props.city) {
      centerOnCity(props.city)
    }

    map.events.add('click', (e) => {
      const coords = e.get('coords')
      setMarker(coords[0], coords[1])
      reverseGeocode(coords[0], coords[1])
    })
  })
})

watch(() => props.city, (newCity) => {
  if (newCity && map && !marker) {
    centerOnCity(newCity)
  }
})

watch(() => props.modelValue, (val) => {
  if (val.lat && val.lng && marker) {
    const coords = marker.geometry.getCoordinates()
    if (Math.abs(coords[0] - val.lat) > 0.0001 || Math.abs(coords[1] - val.lng) > 0.0001) {
      setMarker(val.lat, val.lng)
    }
  }
}, { deep: true })

onUnmounted(() => {
  clearTimeout(debounceTimer)
  if (map) { map.destroy(); map = null }
})
</script>

<style scoped>
.location-picker-map {
  height: 300px;
  width: 100%;
  border-radius: 0.75rem;
  border: 2px solid #e0ddd9;
  z-index: 0;
}

/* Drop map height on mobile so the search input + suggestions fit above the fold */
@media (max-width: 767.98px) {
  .location-picker-map {
    height: 220px;
  }
}
.suggest-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  z-index: 1050;
  background: #fff;
  border: 1px solid #e0ddd9;
  border-radius: 0 0 0.5rem 0.5rem;
  margin-top: -1px;
  max-height: 240px;
  overflow-y: auto;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
  list-style: none;
  padding: 4px 0;
}
.suggest-item {
  padding: 8px 12px;
  cursor: pointer;
  font-size: 0.875rem;
  line-height: 1.3;
  transition: background 0.15s;
  display: flex;
  flex-direction: column;
  gap: 1px;
}
.suggest-item:hover,
.suggest-item.active {
  background: #fef7f4;
}
.suggest-name {
  color: #333;
  font-weight: 500;
}
.suggest-desc {
  color: #888;
  font-size: 0.78rem;
}
</style>
