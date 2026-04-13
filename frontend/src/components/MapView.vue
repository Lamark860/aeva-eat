<template>
  <div ref="mapContainer" :style="{ height: height, width: '100%', borderRadius: '8px' }"></div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

// Fix Leaflet default marker icons in bundlers
delete L.Icon.Default.prototype._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png'
})

const props = defineProps({
  places: { type: Array, default: () => [] },
  center: { type: Array, default: () => [55.7558, 37.6173] },
  zoom: { type: Number, default: 12 },
  height: { type: String, default: '500px' },
  singleMarker: { type: Boolean, default: false }
})

const emit = defineEmits(['marker-click'])

const mapContainer = ref(null)
let map = null
let markersLayer = null

const gemIcon = L.divIcon({
  className: 'gem-marker',
  html: '<div style="background:#0d6efd;color:white;border-radius:50%;width:28px;height:28px;display:flex;align-items:center;justify-content:center;font-size:14px;box-shadow:0 2px 6px rgba(0,0,0,.3);">💎</div>',
  iconSize: [28, 28],
  iconAnchor: [14, 14]
})

function updateMarkers() {
  if (!map || !markersLayer) return
  markersLayer.clearLayers()

  const bounds = []
  props.places.forEach(place => {
    if (!place.lat || !place.lng) return

    const avgRating = place.avg_food && place.avg_service && place.avg_vibe
      ? ((Number(place.avg_food) + Number(place.avg_service) + Number(place.avg_vibe)) / 3).toFixed(1)
      : null

    const hasGem = place.is_gem_place

    const marker = L.marker([place.lat, place.lng], {
      icon: hasGem ? gemIcon : new L.Icon.Default()
    })

    const popupContent = `
      <div style="min-width:150px">
        <strong>${place.name}</strong><br/>
        ${place.cuisine_type ? '<small>🍽 ' + place.cuisine_type + '</small><br/>' : ''}
        ${avgRating ? '<small>⭐ ' + avgRating + '/10</small> · ' : ''}
        <small>${place.review_count || 0} отз.</small><br/>
        <a href="/places/${place.id}">Подробнее →</a>
      </div>
    `
    marker.bindPopup(popupContent)
    marker.on('click', () => emit('marker-click', place))
    markersLayer.addLayer(marker)
    bounds.push([place.lat, place.lng])
  })

  if (bounds.length > 0 && !props.singleMarker) {
    map.fitBounds(bounds, { padding: [30, 30], maxZoom: 15 })
  } else if (bounds.length === 1) {
    map.setView(bounds[0], 15)
  }
}

onMounted(() => {
  map = L.map(mapContainer.value).setView(props.center, props.zoom)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
    maxZoom: 19
  }).addTo(map)

  markersLayer = L.layerGroup().addTo(map)
  updateMarkers()
})

watch(() => props.places, updateMarkers, { deep: true })
watch(() => props.center, (newCenter) => {
  if (map && newCenter) map.setView(newCenter, props.zoom)
})

onUnmounted(() => {
  if (map) {
    map.remove()
    map = null
  }
})
</script>

<style>
.gem-marker {
  background: none !important;
  border: none !important;
}
</style>
