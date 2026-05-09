<template>
  <router-link
    :to="`/places/${place.id}`"
    class="sb-artifact"
    :class="{ featured, 'no-photo': isTicketOnly }"
    :aria-label="place.name"
  >
    <!-- A3 — билетик-only артефакт (NEXT.md §A3). Если у места нет ни cover'а,
         ни review-фото, но есть рейтинги — рендерим компактный билетик
         вместо пустого полароида. Ритм ленты оживает чередованием. -->
    <template v-if="isTicketOnly">
      <div class="art-ticket-card">
        <header class="tc-head">
          <h3 class="tc-name">{{ place.name }}</h3>
          <Stamp v-if="place.city" kind="ink" class="tc-stamp">{{ place.city }}</Stamp>
          <Stamp v-if="place.is_gem_place" kind="gem" class="tc-stamp">жемчужина</Stamp>
        </header>

        <Ticket
          :compact="!featured"
          :food="place.avg_food"
          :service="place.avg_service"
          :vibe="place.avg_vibe"
        />

        <div v-if="caption" class="tc-cap">{{ caption }}</div>

        <div v-if="reviewers.length" class="art-people inline">
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

      <div v-if="metaLine" class="sb-meta">{{ metaLine }}</div>
    </template>

    <!-- Фото-вариант: полароид одиночный или стопка из 3-х -->
    <template v-else>
      <div class="art-photo">
        <Polaroid
          v-if="!hasStack"
          :src="place.image_url || ''"
          :caption="caption"
          :gem="!!place.is_gem_place"
          :placeholder="placeholderClass"
        >
          <Tape :variant="tapeVariant" :style="tapeStyle" />
          <span v-if="place.is_gem_place" class="sb-gem-corner">
            <GemBadge :size="22" />
          </span>
        </Polaroid>

        <PolaroidStack
          v-else
          :photos="stackPhotos"
          :placeholder="placeholderClass"
          :gem="!!place.is_gem_place"
        >
          <Tape :variant="tapeVariant" :style="tapeStyle" />
          <span v-if="place.is_gem_place" class="sb-gem-corner">
            <GemBadge :size="22" />
          </span>
        </PolaroidStack>

        <div v-if="reviewers.length" class="art-people">
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

      <div v-if="hasStack && caption" class="sb-stack-caption">{{ caption }}</div>

      <div v-if="hasRatings" class="sb-ratings sb-t-l1">
        <Ticket
          :compact="!featured"
          :food="place.avg_food"
          :service="place.avg_service"
          :vibe="place.avg_vibe"
        />
      </div>

      <div v-if="metaLine" class="sb-meta">{{ metaLine }}</div>
    </template>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import Polaroid from './Polaroid.vue'
import PolaroidStack from './PolaroidStack.vue'
import Tape from './Tape.vue'
import Ticket from './Ticket.vue'
import Stamp from './Stamp.vue'
import GemBadge from './GemBadge.vue'
import { authorColor, formatVisitCaption } from '../../composables/useFeed'

const props = defineProps({
  place:    { type: Object, required: true },
  // A1 — full-width «звезда» (NEXT.md §A1). Полароид крупнее, текст
  // и билетик чуть жирнее, AuthorTag-стек крупнее. Раскладку (grid-column:
  // 1/-1) определяет родитель; здесь только визуальный апгрейд карточки.
  featured: { type: Boolean, default: false },
})

const placeholderPalette = ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach', 'sb-photo-brick', 'sb-photo-cream', 'sb-photo-slate', 'sb-photo-indigo']
const placeholderClass = computed(() => placeholderPalette[(props.place.id ?? 0) % placeholderPalette.length])

const tapeVariants = ['', 'rose', 'mint', 'blue']
const tapeVariant = computed(() => tapeVariants[(props.place.id ?? 0) % tapeVariants.length])

// Tape position varies per item for a less-uniform feel.
const tapeStyle = computed(() => {
  const variants = [
    { top: '-10px', left: '32px', transform: 'rotate(-12deg)' },
    { top: '-8px',  left: '36px', transform: 'rotate(6deg)' },
    { top: '-8px',  right: '16px', transform: 'rotate(8deg)' },
    { top: '-9px',  right: '28px', transform: 'rotate(-6deg)' },
  ]
  return variants[(props.place.id ?? 0) % variants.length]
})

const caption = computed(() => {
  const d = props.place.created_at || props.place.updated_at
  if (!d) return props.place.name
  return formatVisitCaption(props.place.name, d)
})

// Показываем до 4-х аватарок круга. Остальные сворачиваются в caveat-«+N».
const MAX_REVIEWERS = 4
const allReviewers = computed(() => props.place.reviewers || [])
const reviewers = computed(() => allReviewers.value.slice(0, MAX_REVIEWERS))
const extraReviewers = computed(() => Math.max(0, allReviewers.value.length - MAX_REVIEWERS))

const stackPhotos = computed(() => props.place.feed_photos || [])
const hasStack = computed(() => stackPhotos.value.length >= 2)

const hasRatings = computed(() => {
  const p = props.place
  return [p.avg_food, p.avg_service, p.avg_vibe].some((v) => v !== null && v !== undefined)
})

// A3: ни cover'а, ни review-photos, но место «настоящее» (есть рейтинги или
// хотя бы reviewers). Делаем билетик-only, чтобы не было пустых полароидов.
const isTicketOnly = computed(() => {
  const p = props.place
  const hasAnyPhoto = !!p.image_url || stackPhotos.value.length > 0
  if (hasAnyPhoto) return false
  return hasRatings.value || (p.reviewers || []).length > 0
})

const metaLine = computed(() => {
  const parts = []
  if (props.place.city) parts.push(props.place.city)
  if (props.place.cuisine_type) parts.push(props.place.cuisine_type)
  return parts.join(' · ')
})
</script>

<style scoped lang="scss">
.sb-artifact {
  display: block;
  text-decoration: none;
  color: inherit;
  position: relative;
  width: 100%;
}

.sb-gem-corner {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
}

/* Обёртка вокруг полароида/стопки, чтобы стек авторов мог липнуть к её
   нижнему-левому углу независимо от внутренних transform'ов слоёв. */
.art-photo {
  position: relative;
}

.art-people {
  position: absolute;
  bottom: -10px;
  left: -6px;
  display: inline-flex;
  align-items: center;
  z-index: 4;
  pointer-events: none;
}
/* Сбрасываем absolute из .sb-author-tag — внутри .art-people они в потоке. */
.art-people .r-tag {
  position: relative;
  width: 24px;
  height: 24px;
  margin-left: -8px;
  pointer-events: auto;
  box-shadow:
    0 0 0 2px #fdfcf7,
    0 1px 2px rgba(40, 30, 20, 0.18);
}
.art-people .r-tag:first-child { margin-left: 0; }
.art-people .r-tag.has-photo {
  background: #fdfcf7;
  overflow: hidden;
}
.art-people .r-ph {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.art-people .r-extra {
  margin-left: 6px;
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
  /* «каракулями, не круглый бейдж» — без подложки */
  transform: rotate(-3deg);
}

.sb-ratings {
  display: block;
  margin-top: 10px;
  margin-left: 4px;
}

.sb-meta {
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
  margin-top: 4px;
  padding-left: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* A1 — featured-вариант. Карточка во всю ширину бакета, полароид крупнее,
   текст увереннее, AuthorTag-стек чуть больше. Раскладку колонки задаёт
   родитель (.doska-cell.featured { grid-column: 1 / -1 }). */
.sb-artifact.featured {
  // Явная max-width не нужна — растягивается до ширины ячейки.
  // Внутренний полароид — fluid (width: 100%), он сам займёт всё место.

  .sb-meta {
    font-size: 16px;
    color: var(--sb-ink);
    white-space: normal;
  }

  // Билетик без compact (на узкой колонке зажимали; на full-width дышит).
  .sb-ratings {
    margin-top: 14px;
    margin-left: 8px;
  }

  .sb-stack-caption {
    font-size: 18px;
    color: var(--sb-ink);
  }

  // Аватарки в стеке авторов крупнее — на full-width они теряются 24-pix.
  .art-people {
    bottom: -14px;
    left: -2px;
  }
  .art-people .r-tag {
    width: 30px;
    height: 30px;
    margin-left: -10px;
  }
  .art-people .r-tag:first-child { margin-left: 0; }
  .art-people .r-extra {
    font-size: 16px;
    margin-left: 8px;
  }
}

/* A3 — билетик-only артефакт. Маленькая бумажная карточка с серифа-именем,
   штампиком, билетиком-рейтингом и стеком авторов. Без полароида. */
.art-ticket-card {
  position: relative;
  background: #fdfcf7;
  padding: 14px 14px 24px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 12px rgba(40, 30, 20, 0.08);
  border-radius: 1px;
}
.tc-head {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px 8px;
  margin-bottom: 10px;
}
.tc-name {
  flex: 1 1 100%;
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 18px;
  color: var(--sb-ink);
  line-height: 1.15;
  margin: 0;
  word-break: break-word;
}
.tc-stamp {
  flex: 0 0 auto;
  transform: rotate(-2deg);
}
.tc-cap {
  font-family: var(--sb-hand);
  font-size: 15px;
  color: var(--sb-ink-soft);
  margin-top: 8px;
}

/* Inline-вариант стека авторов внутри билетик-карточки — без абсолюта. */
.art-people.inline {
  position: static;
  margin-top: 10px;
}

/* В stack-режиме caption-полоса полароида не рендерится верхним слоем —
   рисуем подпись отдельно под стопкой, чтобы не ломать наклоны. */
.sb-stack-caption {
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-soft);
  margin-top: 18px;
  padding-left: 4px;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
