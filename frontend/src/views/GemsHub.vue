<template>
  <div class="sb-paper sb-grain sb-screen gems-page">
    <header class="gems-head">
      <div class="back">
        <router-link to="/" class="back-link">← к Доске</router-link>
      </div>

      <div class="title-row">
        <GemBadge :size="28" />
        <h1 class="title">Жемчужины</h1>
      </div>
      <div v-if="!loading" class="sub">{{ subLabel }}</div>
    </header>

    <section v-if="loading" class="sb-empty">листаем заметки…</section>

    <template v-else>
      <!-- B2 — САМАЯ ПЕРВАЯ найденная жемчужина: крупный полароид + история
           «Имя нашла дата». Эталон v3/03-gems-hub.png. -->
      <section v-if="firstGem" class="first-gem">
        <div class="cap-tag">самая первая</div>
        <router-link :to="`/places/${firstGem.id}`" class="first-card">
          <div class="first-photo">
            <Polaroid
              :src="coverOf(firstGem)"
              :width="140"
              :height="140"
              :caption="`${firstGem.name} · ${firstGem.city || ''}`"
              :gem="true"
              :placeholder="placeholderFor(firstGem)"
            >
              <Tape :variant="tapeFor(firstGem)" :style="tapeStyleFor(firstGem)" />
              <span class="gem-corner"><GemBadge :size="18" /></span>
            </Polaroid>
          </div>
          <div class="first-body">
            <div class="first-name">{{ firstGem.name }}</div>
            <div class="first-meta">{{ firstGem.city }}{{ firstGem.cuisine_type ? ` · ${firstGem.cuisine_type}` : '' }}</div>
            <div v-if="firstGemHistory" class="first-history">
              <div class="hist-line">{{ firstGemHistory.first }}</div>
              <div v-if="firstGemHistory.confirmedExtra" class="hist-line dim">{{ firstGemHistory.confirmedExtra }}</div>
            </div>
          </div>
        </router-link>
      </section>

      <!-- По городам: чипы (НЕ вертикальный список). -->
      <div v-if="hub.by_city?.length" class="sb-section-head" style="padding: 0 18px 8px">
        <h2>По городам</h2>
      </div>
      <div v-if="hub.by_city?.length" class="city-chips">
        <router-link
          v-for="c in hub.by_city"
          :key="c.city"
          :to="`/cities/${encodeURIComponent(c.city)}`"
          class="city-chip"
        >
          <Stamp kind="ink">{{ c.city }}</Stamp>
          <span class="chip-count">{{ c.gem_count }}<span class="diamond">◆</span></span>
        </router-link>
      </div>

      <!-- Кто отмечал: ряд аватарок с инициалами/фото -->
      <div v-if="hub.by_user?.length" class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Кто отмечал</h2>
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
          <span class="person-count">{{ u.count }}<span class="diamond">◆</span></span>
        </router-link>
      </div>

      <!-- Все жемчужины: крупные карточки через ArtifactCard.
           Для безфотных автоматом включится PhotoFreeCard (G-layout). -->
      <div class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Все жемчужины</h2>
      </div>
      <div v-if="restGems.length === 0" class="sb-empty">
        пока ничего не отмечали как жемчужину
      </div>
      <div v-else class="all-grid">
        <div v-for="p in restGems" :key="p.id" class="all-cell" :class="cellTiltFor(p)">
          <ArtifactCard :place="p" />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import http from '../api/http'
import ArtifactCard from '../components/scrapbook/ArtifactCard.vue'
import Polaroid from '../components/scrapbook/Polaroid.vue'
import Tape from '../components/scrapbook/Tape.vue'
import Stamp from '../components/scrapbook/Stamp.vue'
import GemBadge from '../components/scrapbook/GemBadge.vue'
import { authorColor } from '../composables/useFeed'

const hub = ref({ places: [], by_city: [], by_user: [], total: 0 })
const loading = ref(false)

function pluralRu(n, forms) {
  const mod100 = n % 100
  if (mod100 >= 11 && mod100 <= 14) return forms[2]
  const last = n % 10
  if (last === 1) return forms[0]
  if (last >= 2 && last <= 4) return forms[1]
  return forms[2]
}

const subLabel = computed(() => {
  const n = hub.value.total || 0
  const cities = hub.value.by_city?.length || 0
  if (n === 0) return 'пусто'
  return `${n} ${pluralRu(n, ['жемчужина', 'жемчужины', 'жемчужин'])} в ${cities} ${pluralRu(cities, ['городе', 'городах', 'городах'])}`
})

const coverOf = (p) => p.image_url || p.feed_photos?.[0]?.url || ''

// Первая найденная — самая ранняя first_marked_at.
const firstGem = computed(() => {
  const list = [...(hub.value.places || [])]
  list.sort((a, b) => {
    const aT = a.gem_status?.first_marked_at || ''
    const bT = b.gem_status?.first_marked_at || ''
    return aT.localeCompare(bT)
  })
  return list[0] || null
})

const restGems = computed(() => {
  if (!firstGem.value) return hub.value.places || []
  return (hub.value.places || []).filter(p => p.id !== firstGem.value.id)
})

const monthsRu = ['января', 'февраля', 'марта', 'апреля', 'мая', 'июня',
                  'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря']
function shortDateRu(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  return `${d.getDate()} ${monthsRu[d.getMonth()]}`
}

// «Алина нашла 5 апреля». Без глагола «нашла/нашёл» (см. R5-Q1) —
// формат как для шапки: «Имя · дата». Если есть соавторы, добавляем
// «подтвердили: имена».
const firstGemHistory = computed(() => {
  const gs = firstGem.value?.gem_status
  if (!gs || !gs.marked_by?.length) return null
  const [first, ...rest] = gs.marked_by
  const dt = shortDateRu(gs.first_marked_at)
  const out = {
    first: `${first.username}${dt ? ` · ${dt}` : ''}`,
    confirmedExtra: '',
  }
  if (rest.length) {
    out.confirmedExtra = `подтвердили: ${rest.map(u => u.username).join(', ')}`
  }
  return out
})

const placeholders = ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach', 'sb-photo-brick']
const placeholderFor = (p) => placeholders[(p?.id ?? 0) % placeholders.length]

const tapeVariants = ['', 'rose', 'mint', 'blue']
const tapeFor = (p) => tapeVariants[(p?.id ?? 0) % tapeVariants.length]
const tapeStyleFor = (p) => {
  const variants = [
    { top: '-9px', left: '50%', transform: 'translateX(-50%) rotate(-8deg)', width: '52px' },
    { top: '-9px', left: '16px', transform: 'rotate(-12deg)', width: '46px' },
  ]
  return variants[(p?.id ?? 0) % variants.length]
}

const tilts = ['sb-t-l2', 'sb-t-r1', 'sb-t-l1', 'sb-t-r2']
const cellTiltFor = (p) => tilts[(p.id ?? 0) % tilts.length]

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
  align-items: baseline;
  gap: 10px;
}
.title {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 36px;
  color: var(--sb-ink);
  margin: 0;
}
.sub {
  font-family: var(--sb-hand);
  font-size: 17px;
  color: var(--sb-ink-mute);
  margin-top: 4px;
}

/* B2 — секция «САМАЯ ПЕРВАЯ» с крупным полароидом и историей справа. */
.first-gem {
  padding: 12px 18px 18px;
}
.cap-tag {
  font-family: var(--sb-hand);
  font-size: 13px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--sb-ink-mute);
  margin-bottom: 8px;
}
.first-card {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  text-decoration: none;
  color: inherit;
}
.first-photo {
  flex: 0 0 auto;
  transform: rotate(-3deg);
}
.gem-corner {
  position: absolute;
  top: 6px;
  right: 6px;
  z-index: 2;
}
.first-body {
  flex: 1 1 auto;
  min-width: 0;
  padding-top: 6px;
}
.first-name {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 22px;
  color: var(--sb-ink);
  line-height: 1.1;
}
.first-meta {
  font-family: var(--sb-hand);
  font-size: 15px;
  color: var(--sb-ink-mute);
  margin-top: 2px;
}
.first-history {
  margin-top: 10px;
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-soft);
  line-height: 1.3;
}
.hist-line.dim {
  color: var(--sb-ink-mute);
  font-size: 14px;
}

/* B2 — города чипами, не вертикальным списком. */
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
.diamond {
  color: var(--sb-terracotta);
  margin-left: 2px;
  font-size: 11px;
}

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
  box-shadow: 0 0 0 2px var(--sb-paper-card), 0 2px 4px rgba(40, 30, 20, 0.18);
}
.person-tile .r-tag.has-photo { background: var(--sb-paper-card); overflow: hidden; }
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

/* Все жемчужины — 2-колоночная сетка как на Доске.
   ArtifactCard сам выберет Polaroid / PhotoFreeCard. */
.all-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 28px 14px;
  padding: 14px 16px 26px;
}
.all-cell {
  min-width: 0;
}
</style>
