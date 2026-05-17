<template>
  <router-link
    :to="`/places/${place._placeId ?? place.id}`"
    class="sb-artifact"
    :class="{ featured, 'no-photo': isTicketOnly, 'has-kruzhok': allVideos.length > 0 }"
    :aria-label="place.name"
  >
    <!-- A3 — безфотный артефакт. PhotoFreeCard сам выбирает раскладку Q/T/G
         по данным: жемчужина → штамп-доминанта, длинная цитата → цитата-
         доминанта, иначе → билетик-доминанта. Эталон v3/01-photofree-card.png.
         Кружочки видео могут быть сбоку — выносим в общий flex-row. -->
    <template v-if="isTicketOnly">
      <div class="art-photo">
        <div class="art-photo-main">
          <PhotoFreeCard
            :place="place"
            :featured="featured"
            :attendees="attendees"
          />
        </div>
        <KruzhokStack :videos="allVideos" :size="featured ? 'large' : 'normal'" />
      </div>
    </template>

    <!-- Фото-вариант: полароид (одиночный или стопка) + кружочки видео в столбце справа -->
    <template v-else>
      <div class="art-photo">
        <div class="art-photo-main">
          <Polaroid
            v-if="!hasStack"
            :src="coverPhoto"
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

        <KruzhokStack :videos="allVideos" :size="featured ? 'large' : 'normal'" />
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

      <!-- Comment на фото-карточке. Раньше показывался только в PFC
           (ticket-only), на фото-варианте текст молча терялся. Если
           короткий — мелкий caveat под meta, если длинный — обрезаем
           многоточием (полный читается на странице места). -->
      <p v-if="commentText" class="sb-comment">«{{ commentText }}»</p>
    </template>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import Polaroid from './Polaroid.vue'
import PolaroidStack from './PolaroidStack.vue'
import PhotoFreeCard from './PhotoFreeCard.vue'
import KruzhokStack from './KruzhokStack.vue'
import Tape from './Tape.vue'
import Ticket from './Ticket.vue'
import GemBadge from './GemBadge.vue'
// authorColor сохраняем — нужен для аватарок в фото-варианте (template v-else)
import { authorColor, formatVisitCaption } from '../../composables/useFeed'

const props = defineProps({
  place:    { type: Object, required: true },
  // A1 — full-width «звезда» (NEXT.md §A1). Полароид крупнее, текст
  // и билетик чуть жирнее, AuthorTag-стек крупнее. Раскладку (grid-column:
  // 1/-1) определяет родитель; здесь только визуальный апгрейд карточки.
  featured: { type: Boolean, default: false },
  // Q4 — если передан, перебивает place.reviewers: показываем только тех,
  // кто был в этом конкретном визите (event-grouped feed). null = не задано.
  attendees: { type: Array, default: null },
})

const placeholderPalette = ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach', 'sb-photo-brick', 'sb-photo-cream', 'sb-photo-slate', 'sb-photo-indigo']
const placeholderClass = computed(() => placeholderPalette[(props.place._placeId ?? props.place.id ?? 0) % placeholderPalette.length])

const tapeVariants = ['', 'rose', 'mint', 'blue']
const tapeVariant = computed(() => tapeVariants[(props.place._placeId ?? props.place.id ?? 0) % tapeVariants.length])

// Tape position varies per item for a less-uniform feel.
const tapeStyle = computed(() => {
  const variants = [
    { top: '-10px', left: '32px', transform: 'rotate(-12deg)' },
    { top: '-8px',  left: '36px', transform: 'rotate(6deg)' },
    { top: '-8px',  right: '16px', transform: 'rotate(8deg)' },
    { top: '-9px',  right: '28px', transform: 'rotate(-6deg)' },
  ]
  return variants[(props.place._placeId ?? props.place.id ?? 0) % variants.length]
})

const caption = computed(() => {
  const d = props.place.created_at || props.place.updated_at
  if (!d) return props.place.name
  // Если визитов больше 1 за неделю — показываем «Винни · ×3 за неделю»
  // вместо точного времени. Время одного из визитов вводило бы в
  // заблуждение: «во сколько ходили» — а ходили несколько раз.
  if (props.place.visits_count > 1) {
    return `${props.place.name} · ×${props.place.visits_count} за неделю`
  }
  return formatVisitCaption(props.place.name, d)
})

// Показываем до 4-х аватарок круга. Остальные сворачиваются в caveat-«+N».
// Q4: если attendees переданы и непустые — рендерим именно их (тех, кто был
// в этом конкретном визите). Иначе фоллбэк на полный place.reviewers —
// чтобы при провале резолва /api/users аватарки не исчезли совсем.
const MAX_REVIEWERS = 4
const allReviewers = computed(() => {
  if (Array.isArray(props.attendees) && props.attendees.length > 0) return props.attendees
  return props.place.reviewers || []
})
const reviewers = computed(() => allReviewers.value.slice(0, MAX_REVIEWERS))
const extraReviewers = computed(() => Math.max(0, allReviewers.value.length - MAX_REVIEWERS))

const stackPhotos = computed(() => props.place.feed_photos || [])
const hasStack = computed(() => stackPhotos.value.length >= 2)

// Q-video: список всех video_url для места. Бэк отдаёт `videos` (свежие
// сверху) или fallback `video_url` для старого ответа. Сам рендер кружочков
// и inline-play вынесены в KruzhokStack — здесь только источник данных.
const allVideos = computed(() => {
  if (Array.isArray(props.place.videos) && props.place.videos.length > 0) {
    return props.place.videos
  }
  if (props.place.video_url) return [props.place.video_url]
  return []
})

const hasRatings = computed(() => {
  const p = props.place
  return [p.avg_food, p.avg_service, p.avg_vibe].some((v) => v !== null && v !== undefined)
})

// Унифицированный cover: либо место имеет своё фото, либо первое из стопки
// review-фото. Если ни того ни другого — это «безфотный» артефакт.
const coverPhoto = computed(() => {
  const p = props.place
  return p.image_url || stackPhotos.value[0]?.url || ''
})

// A3: ни cover'а, ни review-photos, но место «настоящее» (есть рейтинги или
// хотя бы reviewers). Делаем билетик-only, чтобы не было пустых полароидов.
const isTicketOnly = computed(() => {
  if (coverPhoto.value) return false
  return hasRatings.value || (props.place.reviewers || []).length > 0
})

const metaLine = computed(() => {
  const parts = []
  if (props.place.city) parts.push(props.place.city)
  if (props.place.cuisine_type) parts.push(props.place.cuisine_type)
  return parts.join(' · ')
})

const COMMENT_MAX = 120
const commentText = computed(() => {
  const c = (props.place.top_review_comment || '').trim()
  if (!c) return ''
  if (c.length <= COMMENT_MAX) return c
  return c.slice(0, COMMENT_MAX - 1).trimEnd() + '…'
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

/* .art-photo — flex-row контейнер: полароид слева, стопка кружочков справа.
   Описан ниже в Q-video-блоке. */

/* Стек авторов цепляется к низу-левому углу полароидной части (.art-photo-main),
   а не к flex-контейнеру — иначе его сдвигает рост стопки кружочков. */
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
    0 0 0 2px var(--sb-paper-card),
    0 1px 2px rgba(40, 30, 20, 0.18);
}
.art-people .r-tag:first-child { margin-left: 0; }
.art-people .r-tag.has-photo {
  background: var(--sb-paper-card);
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

/* .art-photo — flex-row контейнер: основной контент (полароид или PFC)
   слева, KruzhokStack справа. На узких колонках стопка кружочков мешает
   тексту — переносим её под main через flex-wrap (см. min-width ниже). */
.art-photo {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  flex-wrap: wrap;
}
.art-photo-main {
  position: relative;
  flex: 1 1 auto;
  // Минимальная ширина main-колонки — иначе текст PFC (название, цитата,
  // короткий коммент) превращается в «чек» по одному слову на строку.
  // Если карточка узкая (узкая 2-колоночная сетка + kruzhok сбоку),
  // flex-wrap переносит kruzhok вниз — это нормально для скрапбука.
  min-width: 180px;
}
/* Ticket-only вариант: бумажная плашка не должна растягиваться на всю
   full-width-ячейку — иначе торчит широкий белый фон. */
.sb-artifact.no-photo .art-photo-main {
  flex: 0 1 auto;
  max-width: 360px;
}

/* Comment мелким caveat'ом под meta — впечатление визита, не цитата. */
.sb-comment {
  font-family: var(--sb-hand);
  font-size: 15px;
  line-height: 1.3;
  color: var(--sb-ink-mute);
  margin: 6px 0 0 4px;
  padding: 0;
  word-break: break-word;
}
.sb-artifact.featured .sb-comment {
  font-size: 17px;
  color: var(--sb-ink-soft);
}
</style>
