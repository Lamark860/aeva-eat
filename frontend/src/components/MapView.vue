<template>
  <div ref="mapContainer" :style="{ height: height, width: '100%', borderRadius: '12px', overflow: 'hidden' }"></div>
</template>

<script setup>
/* global ymaps */
import { ref, onMounted, onUnmounted, watch } from 'vue'

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
let placemarks = []

function ratingColor(rating) {
  if (rating >= 8) return '#2eb872'
  if (rating >= 5) return '#f59e0b'
  return '#ef4444'
}

function buildBalloonContent(place, avgRating) {
  return `<div style="min-width:180px;max-width:250px;">
    ${place.image_url ? `<img src="${place.image_url}" style="width:100%;height:100px;object-fit:cover;border-radius:4px 4px 0 0;display:block;" />` : ''}
    <div style="padding:8px 10px;">
      <strong style="font-size:0.9rem;display:block;margin-bottom:2px;">${place.name}</strong>
      ${place.cuisine_type ? `<div style="font-size:0.75rem;color:#888;margin-bottom:4px;">${place.cuisine_type}</div>` : ''}
      <div style="display:flex;align-items:center;gap:8px;margin-bottom:6px;">
        ${avgRating ? `<span style="font-weight:700;font-size:0.95rem;color:#e85d30;">${avgRating}</span>` : ''}
        <span style="font-size:0.75rem;color:#999;">${place.review_count || 0} отз.</span>
      </div>
      <a href="/places/${place.id}" style="font-size:0.8rem;color:#e85d30;text-decoration:none;font-weight:600;">Подробнее →</a>
    </div>
  </div>`
}

function updateMarkers() {
  if (!map) return

  placemarks.forEach(pm => map.geoObjects.remove(pm))
  placemarks = []

  const coords = []
  props.places.forEach(place => {
    if (!place.lat || !place.lng) return

    const avgRating = place.avg_food && place.avg_service && place.avg_vibe
      ? ((Number(place.avg_food) + Number(place.avg_service) + Number(place.avg_vibe)) / 3).toFixed(1)
      : null

    const hasGem = place.is_gem_place
    const color = avgRating ? ratingColor(parseFloat(avgRating)) : '#e85d30'
    const label = avgRating || '?'
    const gemBadge = hasGem
      ? '<circle cx="34" cy="6" r="8" fill="#fff"/><text x="34" y="10" text-anchor="middle" font-size="10">💎</text>'
      : ''

    const iconSvg = `<svg width="42" height="52" viewBox="0 0 42 52" xmlns="http://www.w3.org/2000/svg">
      <defs><filter id="s${place.id}" x="-20%" y="-10%" width="140%" height="130%">
        <feDropShadow dx="0" dy="2" stdDeviation="2" flood-opacity="0.35"/>
      </filter></defs>
      <path d="M21 2C11 2 3 10 3 20c0 14 18 29 18 29s18-15 18-29C39 10 31 2 21 2z"
            fill="${color}" stroke="#fff" stroke-width="2.5" filter="url(#s${place.id})"/>
      <text x="21" y="24" text-anchor="middle" fill="#fff" font-size="12" font-weight="700"
            font-family="Inter,system-ui,sans-serif" style="text-shadow:0 1px 2px rgba(0,0,0,.3)">${label}</text>
      ${gemBadge}
    </svg>`

    const MarkerLayout = ymaps.templateLayoutFactory.createClass(
      '<div style="position:relative;width:42px;height:52px;">' + iconSvg + '</div>'
    )

    const pm = new ymaps.Placemark(
      [place.lat, place.lng],
      {
        hintContent: place.name,
        balloonContentBody: buildBalloonContent(place, avgRating)
      },
      {
        iconLayout: MarkerLayout,
        iconShape: { type: 'Circle', coordinates: [21, 20], radius: 20 },
        iconOffset: [-21, -52],
        hideIconOnBalloonOpen: false
      }
    )

    pm.events.add('click', () => emit('marker-click', place))
    map.geoObjects.add(pm)
    placemarks.push(pm)
    coords.push([place.lat, place.lng])
  })

  if (coords.length > 0 && !props.singleMarker) {
    map.setBounds(
      ymaps.util.bounds.fromPoints(coords),
      { checkZoomRange: true, zoomMargin: 40, duration: 300 }
    ).then(() => {
      if (map.getZoom() > 15) map.setZoom(15)
    })
  } else if (coords.length === 1) {
    map.setCenter(coords[0], 15, { duration: 800 })
  }
}

onMounted(() => {
  ymaps.ready(() => {
    map = new ymaps.Map(mapContainer.value, {
      center: props.center,
      zoom: props.zoom,
      controls: ['zoomControl']
    })

    updateMarkers()
  })
})

watch(() => props.places, updateMarkers, { deep: true })
watch(() => props.center, (newCenter) => {
  if (map && newCenter) map.setCenter(newCenter, props.zoom, { duration: 800 })
})

onUnmounted(() => { if (map) { map.destroy(); map = null } })
</script>
