<template>
  <div class="sb-paper sb-grain sb-screen city-page">
    <header class="city-head">
      <div class="back">
        <router-link to="/places" class="back-link">← к Найти</router-link>
      </div>
      <h1 class="city-name">{{ cityName }}</h1>
      <div v-if="agg" class="meta-line">
        {{ placesLabel }}
        <template v-if="agg.gem_count > 0"> · {{ gemsLabel }}</template>
        <template v-if="agg.contributor_count > 0"> · {{ contributorsLabel }}</template>
      </div>
    </header>

    <section v-if="loading" class="sb-empty">листаем заметки…</section>

    <template v-else>
      <!-- B1 — «Жемчужины Казани» горизонтальный shelf из 3 крупных curated полароидов.
           Не все — кураторский срез. См. v3/02-city-guide.png. -->
      <section v-if="gems.length" class="gems-section">
        <div class="sb-section-head">
          <h2>Жемчужины {{ cityName }}</h2>
          <span class="sub">сюда возвращаются</span>
        </div>
        <div class="gems-shelf">
          <router-link
            v-for="p in featuredGems"
            :key="`gem-${p.id}`"
            :to="`/places/${p.id}`"
            class="gem-card"
          >
            <Polaroid
              :src="coverOf(p)"
              :width="148"
              :height="148"
              :caption="p.name"
              :gem="true"
              :placeholder="placeholderFor(p)"
            >
              <Tape :variant="tapeFor(p)" :style="tapeStyleFor(p)" />
              <span class="gem-corner">
                <GemBadge :size="18" />
              </span>
            </Polaroid>
          </router-link>
        </div>
      </section>

      <!-- B1 — Цитата от круга. Самый длинный комментарий из обзоров мест города. -->
      <section v-if="cityQuote" class="city-quote">
        <div class="cq-cap">цитата от круга</div>
        <blockquote class="cq-text">
          <span class="cq-mark">«</span>{{ cityQuote.text }}<span class="cq-mark">»</span>
        </blockquote>
        <router-link :to="`/places/${cityQuote.placeId}`" class="cq-source">
          — {{ cityQuote.author }} · {{ cityQuote.placeName }}
        </router-link>
      </section>

      <!-- Все места — компактным списком ниже -->
      <div class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Все {{ agg?.count || places.length }}</h2>
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
import Polaroid from '../components/scrapbook/Polaroid.vue'
import Tape from '../components/scrapbook/Tape.vue'
import GemBadge from '../components/scrapbook/GemBadge.vue'

const route = useRoute()
const agg = ref(null)
const places = ref([])
const gems = ref([])
const loading = ref(false)

const cityName = computed(() => decodeURIComponent(route.params.name || ''))

function pluralRu(n, forms) {
  const mod100 = n % 100
  if (mod100 >= 11 && mod100 <= 14) return forms[2]
  const last = n % 10
  if (last === 1) return forms[0]
  if (last >= 2 && last <= 4) return forms[1]
  return forms[2]
}

const placesLabel = computed(() => {
  const n = places.value.length || agg.value?.count || 0
  if (n === 0) return 'ничего нет'
  return `${n} ${pluralRu(n, ['место', 'места', 'мест'])}`
})
const gemsLabel = computed(() => {
  const n = agg.value?.gem_count || 0
  return `${n} ${pluralRu(n, ['жемчужина', 'жемчужины', 'жемчужин'])}`
})
const contributorsLabel = computed(() => {
  const n = agg.value?.contributor_count || 0
  return `из круга ${n}`
})

const coverOf = (p) => p.image_url || p.feed_photos?.[0]?.url || ''

// Берём первые 3 жемчужины города. Если у некоторых нет cover — пускай рисуют
// плейсхолдер с tape — это всё ещё лучше, чем 9 строк подряд.
const featuredGems = computed(() => gems.value.slice(0, 3))

// B1 — цитата от круга. Ищем место с самым длинным top_review_comment.
// Автор — первый из place.reviewers (бэк отдаёт их с /api/cities/:name/places).
const cityQuote = computed(() => {
  const candidates = [...gems.value, ...places.value]
  let best = null
  for (const p of candidates) {
    const c = (p.top_review_comment || '').trim()
    if (!c) continue
    if (!best || c.length > best.c.length) best = { c, p }
  }
  if (!best) return null
  const text = best.c.length > 180 ? best.c.slice(0, 179).trimEnd() + '…' : best.c
  const firstReviewer = (best.p.reviewers || [])[0]
  const author = firstReviewer?.username || 'круг'
  return { text, placeId: best.p.id, placeName: best.p.name, author }
})

const placeholders = ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach', 'sb-photo-brick']
const placeholderFor = (p) => placeholders[(p.id ?? 0) % placeholders.length]

const tapeVariants = ['', 'rose', 'mint', 'blue']
const tapeFor = (p) => tapeVariants[(p.id ?? 0) % tapeVariants.length]
const tapeStyleFor = (p) => {
  const variants = [
    { top: '-9px', left: '40%', transform: 'translateX(-50%) rotate(-8deg)', width: '52px' },
    { top: '-9px', left: '20px', transform: 'rotate(-12deg)', width: '46px' },
    { top: '-9px', right: '20px', transform: 'rotate(10deg)', width: '46px' },
  ]
  return variants[(p.id ?? 0) % variants.length]
}

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
  font-size: 36px;
  color: var(--sb-ink);
  margin: 0 0 6px;
  line-height: 1.05;
  word-break: break-word;
}

/* B1 — meta-строка рукописно вместо билетика «мест/жемчужин/из круга».
   В эталоне это «22 места · 9 жемчужин · вы вN'ом». */
.meta-line {
  font-family: var(--sb-hand);
  font-size: 17px;
  color: var(--sb-ink-mute);
}

.gems-section {
  padding: 12px 0 6px;
}

/* Горизонтальный shelf — 3 кураторских полароида в ряду, скролл если шире. */
.gems-shelf {
  display: flex;
  gap: 14px;
  overflow-x: auto;
  padding: 14px 18px 22px;
  scroll-snap-type: x mandatory;
  scrollbar-width: none;
  &::-webkit-scrollbar { display: none; }
}
.gem-card {
  flex: 0 0 auto;
  text-decoration: none;
  color: inherit;
  scroll-snap-align: start;
  /* Чуть встречный наклон чередуется чтобы было «живо». */
  &:nth-child(1) { transform: rotate(-3deg); }
  &:nth-child(2) { transform: rotate(2deg);  }
  &:nth-child(3) { transform: rotate(-1.5deg); }
}
.gem-corner {
  position: absolute;
  top: 6px;
  right: 6px;
  z-index: 2;
}

/* B1 — цитата от круга. Бумажная плашка, italic-цитата с кавычками-ёлочками,
   подпись «— Имя · место» снизу. */
.city-quote {
  margin: 8px 18px 16px;
  padding: 16px 20px 14px;
  background: var(--sb-paper-card);
  border-radius: 1px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 10px rgba(40, 30, 20, 0.08);
  position: relative;
}
.cq-cap {
  font-family: var(--sb-hand);
  font-size: 13px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--sb-ink-mute);
  margin-bottom: 8px;
}
.cq-text {
  font-family: var(--sb-hand);
  font-size: 22px;
  line-height: 1.25;
  color: var(--sb-ink);
  margin: 0 0 10px;
  hyphens: auto;
  word-break: break-word;
}
.cq-mark {
  font-family: var(--sb-serif);
  color: var(--sb-terracotta);
  opacity: 0.85;
  margin: 0 2px;
}
.cq-source {
  font-family: var(--sb-hand);
  font-size: 15px;
  color: var(--sb-ink-soft);
  text-decoration: none;
  display: block;
  &:hover { color: var(--sb-terracotta); }
}

.city-shelf {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 0 16px;
}
</style>
