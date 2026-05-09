<template>
  <div class="sb-paper sb-grain sb-screen gems-page">
    <header class="gems-head">
      <div class="back">
        <router-link to="/" class="back-link">← к Доске</router-link>
      </div>

      <div class="title-row">
        <GemBadge :size="36" />
        <h1 class="title">Жемчужины круга</h1>
      </div>
      <div class="sub">{{ totalLabel }}</div>
    </header>

    <section v-if="loading" class="sb-empty">листаем заметки…</section>

    <template v-else>
      <div v-if="hub.by_city?.length" class="sb-section-head" style="padding: 0 18px 8px">
        <h2>По городам</h2>
      </div>
      <div v-if="hub.by_city?.length" class="city-list">
        <router-link
          v-for="c in hub.by_city"
          :key="c.city"
          :to="`/cities/${encodeURIComponent(c.city)}`"
          class="city-item"
        >
          <span class="city-name">{{ c.city }}</span>
          <Stamp kind="gem">{{ c.gem_count }}</Stamp>
        </router-link>
      </div>

      <div v-if="hub.by_user?.length" class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Кто отметил</h2>
      </div>
      <div v-if="hub.by_user?.length" class="people-row">
        <router-link
          v-for="u in hub.by_user"
          :key="u.user_id"
          :to="`/people/${u.user_id}`"
          class="person-tile"
        >
          <span
            class="r-tag sb-author-tag"
            :class="[authorColor(u.user_id), { 'has-photo': !!u.avatar_url }]"
            :title="u.username"
          >
            <img v-if="u.avatar_url" :src="u.avatar_url" alt="" class="r-ph" />
            <template v-else>{{ (u.username || '?').slice(0,1).toUpperCase() }}</template>
          </span>
          <span class="person-name">{{ u.username }}</span>
          <span class="person-count">{{ u.count }}</span>
        </router-link>
      </div>

      <div class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Все жемчужины</h2>
      </div>
      <div v-if="(hub.places || []).length === 0" class="sb-empty">
        пока ничего не отмечали как жемчужину
      </div>
      <div v-else class="shelf">
        <ResultCard v-for="p in hub.places" :key="p.id" :place="p" />
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import http from '../api/http'
import ResultCard from '../components/scrapbook/ResultCard.vue'
import Stamp from '../components/scrapbook/Stamp.vue'
import GemBadge from '../components/scrapbook/GemBadge.vue'
import { authorColor } from '../composables/useFeed'

const hub = ref({ places: [], by_city: [], by_user: [], total: 0 })
const loading = ref(false)

const totalLabel = computed(() => {
  const n = hub.value?.total || 0
  if (n === 0) return 'пусто'
  if (n === 1) return '1 жемчужина'
  if (n >= 11 && n <= 14) return `${n} жемчужин`
  const last = n % 10
  if (last === 1) return `${n} жемчужина`
  if (last >= 2 && last <= 4) return `${n} жемчужины`
  return `${n} жемчужин`
})

async function load() {
  loading.value = true
  try {
    const { data } = await http.get('/gems')
    hub.value = data || { places: [], by_city: [], by_user: [], total: 0 }
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped lang="scss">
.gems-page {
  padding-top: calc(18px + var(--aeva-safe-top, 0px));
  padding-inline: 0;
}

.gems-head {
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

.title-row {
  display: flex;
  align-items: center;
  gap: 10px;
}
.title {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 28px;
  color: var(--sb-ink);
  margin: 0;
}
.sub {
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  margin-top: 4px;
}

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
.city-name { flex: 1; word-break: break-word; }

.people-row {
  display: flex;
  flex-wrap: wrap;
  gap: 18px 22px;
  padding: 0 18px;
}
.person-tile {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  text-decoration: none;
  color: inherit;
  width: 72px;
}
.person-tile .r-tag {
  position: relative;
  width: 48px;
  height: 48px;
  font-size: 18px;
  box-shadow: 0 0 0 2px #fdfcf7, 0 2px 4px rgba(40, 30, 20, 0.18);
}
.person-tile .r-tag.has-photo { background: #fdfcf7; overflow: hidden; }
.person-tile .r-ph { width: 100%; height: 100%; object-fit: cover; display: block; }
.person-name {
  font-family: var(--sb-serif);
  font-size: 13px;
  color: var(--sb-ink);
  text-align: center;
  word-break: break-word;
}
.person-count {
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
}

.shelf {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 0 16px;
}
</style>
