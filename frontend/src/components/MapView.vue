<template>
  <div ref="mapContainer" class="sb-map" :style="{ height: height }"></div>
</template>

<script setup>
/* global ymaps */
import { ref, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps({
  places:       { type: Array,   default: () => [] },
  center:       { type: Array,   default: () => [55.7558, 37.6173] },
  zoom:         { type: Number,  default: 12 },
  height:       { type: String,  default: '500px' },
  singleMarker: { type: Boolean, default: false },
})

const emit = defineEmits(['marker-click'])

const mapContainer = ref(null)
let map = null
let placemarks = []

// Scrapbook palette as hex (OKLCH not universally rendered inside Yandex placemarks)
const C = {
  terra:      '#c66b3e',
  terraDark:  '#6e2210',
  ochre:      '#caa44a',
  ochreLight: '#e8c97a',
  ochreDark:  '#7c5f1c',
  moss:       '#7c9b6a',
  mossDark:   '#3f5c30',
  ink:        '#2d231b',
  paper:      'var(--sb-paper-card)',
  highlight:  'rgba(255,255,255,0.45)',
}

// Traffic-light coloring for rated markers — green good / ochre mid / terra bad.
// No rating yet (wishlist or freshly added) → ochre, NOT terra: terra means
// "tried it, was bad" in the light system. Painting unrated places terra would
// make wishlist look like a junkyard. Per DESIGN-DECISIONS.md M1.
function ratingTone(rating) {
  if (rating == null) return { light: C.ochre, dark: C.ochreDark, shade: C.ochreDark }
  const v = Number(rating)
  if (v >= 8.0) return { light: C.moss,  dark: C.mossDark,  shade: C.mossDark  }
  if (v >= 5.0) return { light: C.ochre, dark: C.ochreDark, shade: C.ochreDark }
  return { light: C.terra, dark: C.terraDark, shade: C.terraDark }
}

function ratingLabel(rating) {
  if (rating == null) return ''
  const v = Number(rating)
  // Show whole number when 8/9/10, one decimal otherwise — keeps glyphs readable on a 28px head.
  if (Number.isInteger(v)) return String(v)
  if (v >= 9.95) return '10'
  return v.toFixed(1)
}

function pushpinSVG(id, rating) {
  // 40x52, pin head sits at top, line goes down to point
  const tone = ratingTone(rating)
  const label = ratingLabel(rating)
  const fontSize = label.length > 2 ? 9.5 : 11.5

  return `<svg width="40" height="52" viewBox="0 0 40 52" xmlns="http://www.w3.org/2000/svg">
    <defs>
      <radialGradient id="head-${id}" cx="35%" cy="28%">
        <stop offset="0%" stop-color="${tone.light}"/>
        <stop offset="100%" stop-color="${tone.dark}"/>
      </radialGradient>
      <filter id="shadow-${id}" x="-30%" y="-10%" width="160%" height="160%">
        <feDropShadow dx="0" dy="3" stdDeviation="2.5" flood-opacity="0.4"/>
      </filter>
    </defs>
    <!-- thin needle from head down to point -->
    <line x1="20" y1="20" x2="20" y2="50" stroke="${C.ink}" stroke-width="1.6" stroke-linecap="round" opacity="0.85"/>
    <circle cx="20" cy="50" r="2" fill="${C.ink}" opacity="0.8"/>
    <!-- pushpin head -->
    <circle cx="20" cy="18" r="14" fill="url(#head-${id})" filter="url(#shadow-${id})"/>
    <!-- bottom shading -->
    <ellipse cx="20" cy="22" rx="9" ry="4" fill="${tone.shade}" opacity="0.25"/>
    <!-- highlight -->
    <circle cx="16" cy="13" r="3.5" fill="${C.highlight}"/>
    ${label ? `<text x="20" y="${label.length > 2 ? 21.5 : 22}" text-anchor="middle" fill="#fff" font-family="Lora, Georgia, serif" font-weight="600" font-size="${fontSize}" style="text-shadow:0 1px 1px rgba(0,0,0,0.35)" paint-order="stroke" stroke="rgba(40,30,20,0.35)" stroke-width="0.4">${label}</text>` : ''}
  </svg>`
}

function gemSVG(id, rating) {
  // B6 — на ромбе жемчужины НЕ показываем рейтинг. Ромб сам по себе вердикт
  // «топ» — число (особенно низкое, типа 3.5) вступает с ним в конфликт.
  // Рейтинг видно в балуне при тапе; на маркере — только сияние ромба.
  // 44x54: gem diamond on top with subtle glow + needle line below.
  void rating

  return `<svg width="44" height="54" viewBox="0 0 44 54" xmlns="http://www.w3.org/2000/svg">
    <defs>
      <linearGradient id="gem-${id}" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0%" stop-color="${C.ochreLight}"/>
        <stop offset="100%" stop-color="${C.ochre}"/>
      </linearGradient>
      <filter id="g-shadow-${id}" x="-40%" y="-20%" width="180%" height="160%">
        <feDropShadow dx="0" dy="2" stdDeviation="2.2" flood-opacity="0.32"/>
      </filter>
      <radialGradient id="halo-${id}" cx="50%" cy="40%">
        <stop offset="0%" stop-color="${C.ochreLight}" stop-opacity="0.6"/>
        <stop offset="60%" stop-color="${C.ochreLight}" stop-opacity="0.0"/>
      </radialGradient>
    </defs>
    <!-- soft halo -->
    <circle cx="22" cy="18" r="22" fill="url(#halo-${id})"/>
    <!-- needle -->
    <line x1="22" y1="32" x2="22" y2="52" stroke="${C.ink}" stroke-width="1.6" stroke-linecap="round" opacity="0.85"/>
    <circle cx="22" cy="52" r="2" fill="${C.ink}" opacity="0.8"/>
    <!-- gem diamond -->
    <path d="M22 4 L36 18 L22 32 L8 18 Z"
          fill="url(#gem-${id})"
          stroke="${C.terraDark}"
          stroke-width="1.6"
          stroke-linejoin="round"
          filter="url(#g-shadow-${id})"/>
    <!-- facets -->
    <path d="M22 4 L16 18 L22 32 M22 4 L28 18 L22 32 M8 18 L36 18"
          stroke="${C.terraDark}" stroke-width="0.9" fill="none" opacity="0.55"/>
    <!-- highlight -->
    <path d="M14 9 L17 9 L11 19 L8 19 Z" fill="#fff" opacity="0.6"/>
  </svg>`
}

function buildBalloonContent(place, avgRating) {
  // All inline styles — Yandex balloon strips classes, so we can't rely on stylesheet.
  const cover = place.image_url
    ? `<div style="background:var(--sb-paper-card);padding:6px 6px 22px;margin:0 0 8px;box-shadow:0 1px 1px rgba(40,30,20,.08),0 4px 10px rgba(40,30,20,.10);border-radius:1px;display:inline-block;transform:rotate(-1.5deg);position:relative">
         <img src="${place.image_url}" style="display:block;width:160px;height:110px;object-fit:cover;border-radius:1px" />
         <div style="position:absolute;left:0;right:0;bottom:4px;text-align:center;font-family:Caveat,cursive;font-size:14px;color:#5a4a3c;line-height:1">${escapeHtml(place.name)}</div>
       </div>`
    : ''

  const meta = []
  if (place.city) meta.push(escapeHtml(place.city))
  if (place.cuisine_type) meta.push(escapeHtml(place.cuisine_type))
  const metaHtml = meta.length
    ? `<div style="font-family:Caveat,cursive;font-size:15px;color:#7a6a5c;margin-top:4px">${meta.join(' · ')}</div>`
    : ''

  const ticket = avgRating
    ? `<div style="display:inline-flex;align-items:center;background:#ead8a3;padding:4px 9px;border-radius:2px;font-family:Lora,Georgia,serif;margin-top:8px;box-shadow:0 1px 1px rgba(40,30,20,.08)">
         <span style="font-size:8px;letter-spacing:.18em;text-transform:uppercase;color:#7a6a5c;margin-right:6px">общая</span>
         <span style="font-family:Caveat,cursive;font-size:20px;color:#2d231b;line-height:1">${avgRating}</span>
       </div>`
    : ''

  const gem = place.is_gem_place
    ? `<span style="display:inline-block;font-family:Lora,Georgia,serif;font-weight:600;font-size:9.5px;letter-spacing:.18em;text-transform:uppercase;color:#c66b3e;border:1.4px solid #c66b3e;padding:3px 7px 2px;border-radius:2px;background:rgba(232,201,122,.5);margin-left:8px">жемчужина</span>`
    : ''

  // Author stack — DESIGN-DECISIONS M2. 22px circles, max 4, no «с нами:» label,
  // only a colon as visual hint that what follows is people. Avatars when set,
  // colored initials otherwise.
  const authorPalette = {
    terra: { bg: '#dcb19c', ink: '#5a2b1c' },
    ochre: { bg: '#e8d29a', ink: '#5d4a14' },
    moss:  { bg: '#bdcfb1', ink: '#33472a' },
    plum:  { bg: '#cdb1be', ink: '#4a2438' },
  }
  const palette = ['terra', 'ochre', 'moss', 'plum']
  const reviewers = (place.reviewers || []).slice(0, 4)
  const peopleHtml = reviewers.length
    ? `<div style="display:inline-flex;align-items:center;margin-top:10px;font-family:Caveat,cursive;font-size:15px;color:#7a6a5c">
         <span style="margin-right:6px">:</span>
         ${reviewers.map((r, i) => {
           const tone = authorPalette[palette[Math.abs(r.id || 0) % palette.length]]
           const initial = escapeHtml((r.username || '?').slice(0, 1).toUpperCase())
           const inner = r.avatar_url
             ? `<img src="${r.avatar_url}" alt="" style="width:100%;height:100%;object-fit:cover;display:block" />`
             : initial
           const bg = r.avatar_url ? 'var(--sb-paper-card)' : tone.bg
           return `<span title="${escapeHtml(r.username || '')}" style="display:inline-flex;align-items:center;justify-content:center;width:22px;height:22px;border-radius:50%;background:${bg};color:${tone.ink};font-family:Lora,Georgia,serif;font-weight:600;font-size:10px;box-shadow:0 0 0 2px var(--sb-paper-card),0 1px 2px rgba(40,30,20,.25);margin-left:${i === 0 ? '0' : '-6px'};overflow:hidden">${inner}</span>`
         }).join('')}
       </div>`
    : ''

  return `<div style="min-width:200px;max-width:260px;font-family:Lora,Georgia,serif;color:#2d231b">
    ${cover}
    <div style="font-family:Lora,Georgia,serif;font-style:italic;font-weight:500;font-size:17px;line-height:1.15;color:#2d231b">${escapeHtml(place.name)}${gem}</div>
    ${metaHtml}
    <div>${ticket}</div>
    ${peopleHtml}
    <div><a href="/places/${place.id}" style="display:inline-block;font-family:Lora,Georgia,serif;font-style:italic;font-size:14px;color:#c66b3e;text-decoration:none;margin-top:10px">подробнее →</a></div>
  </div>`
}

function escapeHtml(s) {
  return String(s ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function updateMarkers() {
  if (!map) return

  placemarks.forEach((pm) => map.geoObjects.remove(pm))
  placemarks = []

  const coords = []
  props.places.forEach((place) => {
    if (!place.lat || !place.lng) return

    const avgRating =
      place.avg_food && place.avg_service && place.avg_vibe
        ? ((Number(place.avg_food) + Number(place.avg_service) + Number(place.avg_vibe)) / 3).toFixed(1)
        : null

    const id = place.id
    const isGem = place.is_gem_place
    const svg = isGem ? gemSVG(id, avgRating) : pushpinSVG(id, avgRating)
    const w = isGem ? 44 : 40
    const h = isGem ? 54 : 52

    const MarkerLayout = ymaps.templateLayoutFactory.createClass(
      `<div style="position:relative;width:${w}px;height:${h}px;">${svg}</div>`,
    )

    const pm = new ymaps.Placemark(
      [place.lat, place.lng],
      {
        hintContent: place.name,
        balloonContentBody: buildBalloonContent(place, avgRating),
      },
      {
        iconLayout: MarkerLayout,
        // anchor at the needle tip (bottom-center)
        iconShape: { type: 'Rectangle', coordinates: [[-w / 2, -h], [w / 2, 0]] },
        iconOffset: [-(w / 2), -h],
        hideIconOnBalloonOpen: false,
      },
    )

    pm.events.add('click', () => emit('marker-click', place))
    map.geoObjects.add(pm)
    placemarks.push(pm)
    coords.push([place.lat, place.lng])
  })

  if (coords.length > 0 && !props.singleMarker) {
    map
      .setBounds(ymaps.util.bounds.fromPoints(coords), {
        checkZoomRange: true,
        zoomMargin: 40,
        duration: 300,
      })
      .then(() => {
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
      controls: ['zoomControl'],
    })
    updateMarkers()
  })
})

watch(() => props.places, updateMarkers, { deep: true })
watch(
  () => props.center,
  (newCenter) => {
    if (map && newCenter) map.setCenter(newCenter, props.zoom, { duration: 800 })
  },
)

onUnmounted(() => {
  if (map) {
    map.destroy()
    map = null
  }
})
</script>

<style scoped lang="scss">
.sb-map {
  width: 100%;
  border-radius: 4px;
  overflow: hidden;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 14px rgba(40, 30, 20, 0.07);
}
</style>
