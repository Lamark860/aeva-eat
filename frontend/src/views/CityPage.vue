<template>
  <div class="sb-paper sb-grain sb-screen city-page">
    <header class="city-head">
      <div class="back">
        <router-link to="/places" class="back-link">← к Найти</router-link>
      </div>
      <h1 class="city-name">{{ cityName }}</h1>
      <div v-if="agg" class="stats sb-t-l1">
        <Ticket
          :food="agg.count"
          :service="agg.gem_count"
          :vibe="agg.contributor_count"
          :labels="['мест', 'жемчужин', 'из&nbsp;круга']"
        />
      </div>
    </header>

    <section v-if="loading" class="sb-empty">листаем заметки…</section>

    <template v-else>
      <div v-if="gems.length" class="sb-section-head" style="padding: 0 18px 8px">
        <h2>Жемчужины</h2>
        <span class="sub">сюда возвращаются</span>
      </div>
      <div v-if="gems.length" class="city-shelf">
        <ResultCard v-for="p in gems" :key="`gem-${p.id}`" :place="p" />
      </div>

      <div class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Все места</h2>
        <span class="sub">{{ placesLabel }}</span>
      </div>

      <div v-if="places.length === 0" class="sb-empty">
        в этом городе пока ничего не прикноплено
      </div>
      <div v-else class="city-shelf">
        <ResultCard v-for="p in places" :key="p.id" :place="p" />
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import http from '../api/http'
import ResultCard from '../components/scrapbook/ResultCard.vue'
import Ticket from '../components/scrapbook/Ticket.vue'

const route = useRoute()
const agg = ref(null)
const places = ref([])
const gems = ref([])
const loading = ref(false)

const cityName = computed(() => decodeURIComponent(route.params.name || ''))

const placesLabel = computed(() => {
  const n = places.value.length
  if (n === 0) return 'ничего нет'
  if (n === 1) return '1 место'
  if (n >= 11 && n <= 14) return `${n} мест`
  const last = n % 10
  if (last === 1) return `${n} место`
  if (last >= 2 && last <= 4) return `${n} места`
  return `${n} мест`
})

async function load() {
  loading.value = true
  try {
    const name = encodeURIComponent(cityName.value)
    const [a, p, g] = await Promise.all([
      http.get(`/cities/${name}`),
      http.get(`/cities/${name}/places`),
      http.get(`/cities/${name}/gems`),
    ])
    agg.value = a.data
    places.value = p.data?.places || []
    gems.value = g.data?.places || []
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(() => route.params.name, () => { if (route.params.name) load() })
</script>

<style scoped lang="scss">
.city-page {
  padding-top: calc(18px + var(--aeva-safe-top, 0px));
  padding-inline: 0;
}

.city-head {
  padding: 0 18px 14px;
}
.back {
  margin-bottom: 8px;
}
.back-link {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink-mute);
  text-decoration: none;
  &:hover { color: var(--sb-ink); }
}

.city-name {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 32px;
  color: var(--sb-ink);
  margin: 0 0 12px;
  line-height: 1.05;
  word-break: break-word;
}

.stats { display: inline-block; }

.city-shelf {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 0 16px;
}
</style>
