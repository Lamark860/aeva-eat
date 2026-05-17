<template>
  <div class="sb-paper sb-grain sb-screen person-page">
    <header class="person-head">
      <div class="back">
        <router-link to="/" class="back-link">← к Доске</router-link>
      </div>

      <div v-if="profile" class="person-id">
        <span
          class="r-tag sb-author-tag big"
          :class="[authorColor(profile.id), { 'has-photo': !!profile.avatar_url }]"
          :title="profile.username"
        >
          <img v-if="profile.avatar_url" :src="profile.avatar_url" alt="" class="r-ph" />
          <template v-else>{{ initial }}</template>
        </span>
        <h1 class="name">{{ profile.username }}</h1>
      </div>

      <div v-if="profile" class="stats sb-t-r1">
        <Ticket
          :food="profile.place_count"
          :service="profile.gem_count"
          :vibe="profile.city_count"
          :labels="['мест', 'жемчужин', 'городов']"
        />
      </div>

      <div v-if="cuisineLine" class="cuisine-line">{{ cuisineLine }}</div>
    </header>

    <section v-if="loading" class="sb-empty">листаем заметки…</section>

    <template v-else>
      <!-- B8 — Города чипами (не вертикальным списком на 30+ строк).
           Одинаковый словарь со страницами B2 (Gems Hub) и B3 (Place header). -->
      <div v-if="cities.length" class="sb-section-head" style="padding: 0 18px 8px">
        <h2>Города</h2>
        <span class="sub">где был(а)</span>
      </div>
      <div v-if="cities.length" class="city-chips">
        <router-link
          v-for="c in cities"
          :key="c.city"
          :to="`/cities/${encodeURIComponent(c.city)}`"
          class="city-chip"
        >
          <Stamp kind="ink">{{ c.city }}</Stamp>
          <span class="chip-count">{{ c.count }}</span>
        </router-link>
      </div>

      <!-- B8 — Жемчужины 2-колоночной сеткой через ArtifactCard
           (для безфотных автоматом PhotoFreeCard G-layout). -->
      <div v-if="gems.length" class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Жемчужины</h2>
        <span class="sub">{{ gemsLabel }}</span>
      </div>
      <div v-if="gems.length" class="all-grid">
        <div v-for="p in gems" :key="`gem-${p.id}`" class="all-cell" :class="cellTiltFor(p)">
          <ArtifactCard :place="p" />
        </div>
      </div>

      <!-- B8 — Визиты: лимит 6, остаток через «↓ ещё N». При 144 местах
           страница больше не 32k px, а ~3k px пока юзер не раскроет. -->
      <div v-if="places.length" class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Визиты</h2>
        <span class="sub">{{ visitsLabel }}</span>
      </div>
      <div v-if="places.length === 0" class="sb-empty">
        пока ничего не прикноплено
      </div>
      <div v-else class="visit-shelf">
        <ResultCard v-for="p in visibleVisits" :key="p.id" :place="p" />
      </div>

      <div v-if="hiddenCount > 0" class="more-row">
        <button class="more-btn" type="button" @click="visibleVisitCount += 12">
          ↓ ещё {{ Math.min(hiddenCount, 12) }} <span v-if="hiddenCount > 12">из {{ hiddenCount }}</span>
        </button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import http from '../api/http'
import ArtifactCard from '../components/scrapbook/ArtifactCard.vue'
import ResultCard from '../components/scrapbook/ResultCard.vue'
import Ticket from '../components/scrapbook/Ticket.vue'
import Stamp from '../components/scrapbook/Stamp.vue'
import { authorColor } from '../composables/useFeed'
import { favoriteCuisinePhrase } from '../composables/useCuisine'

const route = useRoute()
const profile = ref(null)
const cuisineLine = computed(() => favoriteCuisinePhrase(profile.value))
const places = ref([])
const gems = ref([])
const cities = ref([])
const loading = ref(false)
const visibleVisitCount = ref(6)

const initial = computed(() =>
  (profile.value?.username || '?').slice(0, 1).toUpperCase()
)

function pluralRu(n, forms) {
  const mod100 = n % 100
  if (mod100 >= 11 && mod100 <= 14) return forms[2]
  const last = n % 10
  if (last === 1) return forms[0]
  if (last >= 2 && last <= 4) return forms[1]
  return forms[2]
}

const gemsLabel = computed(() => {
  const n = gems.value.length
  return `${n} ${pluralRu(n, ['находка', 'находки', 'находок'])}`
})
const visitsLabel = computed(() => {
  const n = places.value.length
  return `всего ${n}`
})

const visibleVisits = computed(() => places.value.slice(0, visibleVisitCount.value))
const hiddenCount = computed(() => Math.max(0, places.value.length - visibleVisitCount.value))

const tilts = ['sb-t-l2', 'sb-t-r1', 'sb-t-l1', 'sb-t-r2']
const cellTiltFor = (p) => tilts[(p.id ?? 0) % tilts.length]

async function load() {
  const id = route.params.id
  if (!id) return
  loading.value = true
  visibleVisitCount.value = 6
  try {
    const [u, p, g, c] = await Promise.all([
      http.get(`/users/${id}`),
      http.get(`/users/${id}/places`),
      http.get(`/users/${id}/gems`),
      http.get(`/users/${id}/cities`),
    ])
    profile.value = u.data
    places.value = p.data || []
    gems.value = g.data || []
    cities.value = c.data || []
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(() => route.params.id, () => { if (route.params.id) load() })
</script>

<style scoped lang="scss">
.person-page {
  padding-top: calc(18px + var(--aeva-safe-top, 0px));
  padding-inline: 0;
}

.person-head {
  padding: 0 18px 14px;
}
.back { margin-bottom: 8px; }
.back-link {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink-mute);
  text-decoration: none;
  &:hover { color: var(--sb-ink); }
}

.person-id {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}
.r-tag.big {
  position: relative;
  width: 56px;
  height: 56px;
  font-size: 22px;
  box-shadow: 0 0 0 2px var(--sb-paper-card), 0 2px 6px rgba(40, 30, 20, 0.18);
}
.r-tag.has-photo { background: var(--sb-paper-card); overflow: hidden; }
.r-ph { width: 100%; height: 100%; object-fit: cover; display: block; }

.name {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 30px;
  color: var(--sb-ink);
  margin: 0;
  word-break: break-word;
}

.stats { display: inline-block; }

.cuisine-line {
  margin-top: 10px;
  font-family: var(--sb-hand);
  font-size: 18px;
  color: var(--sb-ink-soft);
  line-height: 1.2;
}

/* B8 — города чипами */
.city-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 10px;
  padding: 0 18px 4px;
}
.city-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  text-decoration: none;
  color: inherit;
}
.chip-count {
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
}

/* B8 — 2-колоночная сетка для жемчужин */
.all-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 28px 14px;
  padding: 14px 16px 12px;
}
.all-cell {
  min-width: 0;
}

.visit-shelf {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 0 16px;
}

.more-row {
  text-align: center;
  padding: 14px 0 24px;
}
.more-btn {
  background: none;
  border: none;
  font-family: var(--sb-hand);
  font-size: 17px;
  color: var(--sb-ink-soft);
  cursor: pointer;
  padding: 6px 10px;
  &:hover { color: var(--sb-terracotta); }
}
</style>
