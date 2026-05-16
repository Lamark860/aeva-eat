<template>
  <!-- A3 — безфотный артефакт в трёх раскладках. Эталон v3/01-photofree-card.png.
       Layout выбирается автоматически: жемчужина → G, длинная цитата → Q,
       иначе → T. Всё на бумажной плашке с tape, как полароид-замена. -->
  <div class="pfc" :class="[`pfc--${layout}`, { 'pfc--featured': featured }]">
    <Tape :variant="tapeVariant" :style="tapeStyle" />

    <!-- G: штамп-доминанта (жемчужина без фото) -->
    <template v-if="layout === 'G'">
      <div class="pfc-gem-stamp">
        <Stamp kind="gem" class="pfc-stamp-gem">жемчужина</Stamp>
      </div>
      <h3 class="pfc-name pfc-name--lg">{{ place.name }}</h3>
      <div v-if="metaLine" class="pfc-meta">{{ metaLine }}</div>
      <Ticket
        v-if="hasRatings"
        :compact="!featured"
        class="pfc-ticket"
        :food="place.avg_food"
        :service="place.avg_service"
        :vibe="place.avg_vibe"
      />
    </template>

    <!-- Q: цитата-доминанта (есть длинный коммент) -->
    <template v-else-if="layout === 'Q'">
      <blockquote class="pfc-quote">
        <span class="pfc-quote-mark">«</span>{{ quoteText }}<span class="pfc-quote-mark">»</span>
      </blockquote>
      <h3 class="pfc-name pfc-name--sm">{{ place.name }}</h3>
      <div v-if="metaLine" class="pfc-meta">{{ metaLine }}</div>
      <Ticket
        v-if="hasRatings"
        compact
        class="pfc-ticket pfc-ticket--side"
        :food="place.avg_food"
        :service="place.avg_service"
        :vibe="place.avg_vibe"
      />
    </template>

    <!-- T: билетик-доминанта (нет цитаты, нет жемчужины) -->
    <template v-else>
      <h3 class="pfc-name pfc-name--lg">{{ place.name }}</h3>
      <div v-if="metaLine" class="pfc-meta">{{ metaLine }}</div>
      <Ticket
        v-if="hasRatings"
        :compact="!featured"
        class="pfc-ticket pfc-ticket--big"
        :food="place.avg_food"
        :service="place.avg_service"
        :vibe="place.avg_vibe"
      />
      <div v-if="visitCaption" class="pfc-cap">{{ visitCaption }}</div>
    </template>

    <div v-if="reviewers.length" class="pfc-people">
      <span
        v-for="r in reviewers"
        :key="r.id"
        class="r-tag sb-author-tag"
        :class="[authorColor(r.id), { 'has-photo': !!r.avatar_url }]"
        :title="r.username"
      >
        <img v-if="r.avatar_url" :src="r.avatar_url" alt="" class="r-ph" />
        <template v-else>{{ (r.username || '?').slice(0, 1).toUpperCase() }}</template>
      </span>
      <span v-if="extraReviewers > 0" class="r-extra">+{{ extraReviewers }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import Tape from './Tape.vue'
import Ticket from './Ticket.vue'
import Stamp from './Stamp.vue'
import { authorColor, formatVisitCaption } from '../../composables/useFeed'

const props = defineProps({
  place:    { type: Object, required: true },
  featured: { type: Boolean, default: false },
  attendees: { type: Array, default: null },
})

const MAX_REVIEWERS = 4
const QUOTE_MAX = 140

const allReviewers = computed(() => {
  if (Array.isArray(props.attendees) && props.attendees.length > 0) return props.attendees
  return props.place.reviewers || []
})
const reviewers = computed(() => allReviewers.value.slice(0, MAX_REVIEWERS))
const extraReviewers = computed(() => Math.max(0, allReviewers.value.length - MAX_REVIEWERS))

const hasRatings = computed(() => {
  const p = props.place
  return [p.avg_food, p.avg_service, p.avg_vibe].some(v => v !== null && v !== undefined)
})

const metaLine = computed(() => {
  const parts = []
  if (props.place.city) parts.push(props.place.city)
  if (props.place.cuisine_type) parts.push(props.place.cuisine_type)
  return parts.join(' · ')
})

const visitCaption = computed(() => {
  const d = props.place.created_at || props.place.updated_at
  if (!d) return ''
  return formatVisitCaption('', d).trim()
})

const quoteText = computed(() => {
  const c = (props.place.top_review_comment || '').trim()
  if (c.length <= QUOTE_MAX) return c
  return c.slice(0, QUOTE_MAX - 1).trimEnd() + '…'
})

// G > Q > T — приоритет: если место gem, оно всегда штамп-доминанта,
// даже если есть цитата (штамп визуально сильнее). Дальше Q если есть текст,
// иначе fallback на билетик-доминанту.
const layout = computed(() => {
  if (props.place.is_gem_place) return 'G'
  if ((props.place.top_review_comment || '').trim().length >= 30) return 'Q'
  return 'T'
})

const tapeVariants = ['', 'rose', 'mint', 'blue']
const tapeVariant = computed(() => tapeVariants[(props.place.id ?? 0) % tapeVariants.length])
const tapeStyle = computed(() => {
  const variants = [
    { top: '-9px', left: '24px', transform: 'rotate(-10deg)' },
    { top: '-8px', left: '32px', transform: 'rotate(6deg)' },
    { top: '-8px', right: '18px', transform: 'rotate(8deg)' },
    { top: '-9px', right: '30px', transform: 'rotate(-6deg)' },
  ]
  return variants[(props.place.id ?? 0) % variants.length]
})
</script>

<style scoped lang="scss">
.pfc {
  position: relative;
  padding: 22px 18px 18px;
  background: var(--sb-paper-card);
  border-radius: 1px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 14px rgba(40, 30, 20, 0.08);
  display: flex;
  flex-direction: column;
  gap: 6px;
  color: inherit;
  text-decoration: none;
  min-height: 156px;
}
.pfc--featured {
  padding: 28px 22px 22px;
  min-height: 200px;
}

.pfc-name {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  color: var(--sb-ink);
  margin: 0;
  line-height: 1.1;
}
.pfc-name--lg {
  font-size: 26px;
}
.pfc--featured .pfc-name--lg {
  font-size: 32px;
}
.pfc-name--sm {
  font-size: 16px;
  font-style: italic;
  color: var(--sb-ink-mute);
  margin-top: 4px;
}

.pfc-meta {
  font-family: var(--sb-hand);
  font-size: 15px;
  color: var(--sb-ink-mute);
  margin-top: -2px;
}

.pfc-ticket {
  margin-top: 8px;
  align-self: flex-start;
}
.pfc-ticket--big {
  margin-top: 10px;
  transform: scale(1.1);
  transform-origin: left top;
}
.pfc-ticket--side {
  position: absolute;
  right: 14px;
  bottom: 14px;
  transform: scale(0.9) rotate(-3deg);
  transform-origin: right bottom;
}

.pfc-cap {
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
  margin-top: 6px;
}

/* Q — крупная рукописная цитата с кавычками-ёлочками */
.pfc-quote {
  font-family: var(--sb-hand);
  font-size: 22px;
  line-height: 1.25;
  color: var(--sb-ink);
  margin: 4px 0 8px;
  padding: 0;
  position: relative;
  hyphens: auto;
  word-break: break-word;
}
.pfc--featured .pfc-quote {
  font-size: 26px;
}
.pfc-quote-mark {
  font-family: var(--sb-serif);
  color: var(--sb-terracotta);
  font-size: 1.2em;
  opacity: 0.85;
  margin: 0 1px;
}

/* G — штамп ЖЕМЧУЖИНА сверху, как заголовок */
.pfc-gem-stamp {
  align-self: flex-start;
  margin-bottom: 4px;
}
.pfc-stamp-gem {
  font-size: 13px;
}

.pfc-people {
  display: inline-flex;
  align-items: center;
  margin-top: 6px;
  align-self: flex-start;
}
.pfc-people .r-tag {
  position: relative;
  width: 24px;
  height: 24px;
  margin-left: -8px;
}
.pfc-people .r-tag:first-child {
  margin-left: 0;
}
.pfc-people .r-tag.has-photo {
  background: var(--sb-paper-card);
  overflow: hidden;
}
.pfc-people .r-ph {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.r-extra {
  margin-left: 4px;
  font-family: var(--sb-hand);
  font-size: 13px;
  color: var(--sb-ink-mute);
}
</style>
