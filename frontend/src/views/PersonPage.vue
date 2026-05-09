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
    </header>

    <section v-if="loading" class="sb-empty">листаем заметки…</section>

    <template v-else>
      <div v-if="cities.length" class="sb-section-head" style="padding: 0 18px 8px">
        <h2>Города</h2>
        <span class="sub">где был(а)</span>
      </div>
      <div v-if="cities.length" class="city-list">
        <router-link
          v-for="c in cities"
          :key="c.city"
          :to="`/cities/${encodeURIComponent(c.city)}`"
          class="city-item"
        >
          <span class="city-name">{{ c.city }}</span>
          <Stamp kind="ink">{{ c.count }}</Stamp>
        </router-link>
      </div>

      <div v-if="gems.length" class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Жемчужины</h2>
      </div>
      <div v-if="gems.length" class="shelf">
        <ResultCard v-for="p in gems" :key="`gem-${p.id}`" :place="p" />
      </div>

      <div class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Визиты</h2>
      </div>
      <div v-if="places.length === 0" class="sb-empty">
        пока ничего не прикноплено
      </div>
      <div v-else class="shelf">
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
import Stamp from '../components/scrapbook/Stamp.vue'
import { authorColor } from '../composables/useFeed'

const route = useRoute()
const profile = ref(null)
const places = ref([])
const gems = ref([])
const cities = ref([])
const loading = ref(false)

const initial = computed(() =>
  (profile.value?.username || '?').slice(0, 1).toUpperCase()
)

async function load() {
  const id = route.params.id
  if (!id) return
  loading.value = true
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
  box-shadow: 0 0 0 2px #fdfcf7, 0 2px 6px rgba(40, 30, 20, 0.18);
}
.r-tag.has-photo { background: #fdfcf7; overflow: hidden; }
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

.city-list {
  display: flex;
  flex-direction: column;
  padding: 0 16px;
}
.city-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 4px;
  border-bottom: 1px dashed rgba(40, 30, 20, 0.18);
  text-decoration: none;
  color: var(--sb-ink);
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 17px;

  &:hover .city-name { color: var(--sb-terracotta); }
}
.city-name {
  flex: 1;
  word-break: break-word;
}

.shelf {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 0 16px;
}
</style>
