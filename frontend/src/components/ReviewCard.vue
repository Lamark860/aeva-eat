<template>
  <article class="sb-review">
    <router-link
      v-if="showPlace && review.place_id"
      :to="`/places/${review.place_id}`"
      class="place-link"
    >
      <span class="place-name">{{ review.place_name || `место #${review.place_id}` }}</span>
      <span class="arr">→</span>
    </router-link>

    <header class="sb-review-head">
      <div class="authors">
        <router-link
          v-for="(author, i) in review.authors"
          :key="author.id"
          :to="`/people/${author.id}`"
          class="author-line"
        >
          <AuthorTag
            :name="author.username"
            :color="authorColor(author.id)"
            :src="author.avatar_url || ''"
            class="static"
            :title="author.username"
          />
          <span class="username">{{ author.username }}</span>
          <span v-if="i < review.authors.length - 1" class="sep">·</span>
        </router-link>
        <span v-if="review.authors && review.authors.length > 1" class="joint">вместе</span>
      </div>
      <div v-if="review.visited_at" class="date">{{ formatDate(review.visited_at) }}</div>

      <div v-if="canEdit" class="actions">
        <button class="link-btn" @click="$emit('edit', review)" title="Редактировать">ред</button>
        <button class="link-btn danger" @click="$emit('delete', review.id)" title="Удалить">×</button>
      </div>
    </header>

    <!-- Photo polaroids — на детали места все фото в линию,
         горизонтальный скролл (DESIGN-DECISIONS §L3). Если фото нет, но
         есть legacy review.image_url — показываем его одиночным полароидом. -->
    <div v-if="photoList.length > 0" class="sb-review-photos">
      <div class="photos-row">
        <Polaroid
          v-for="(p, i) in photoList"
          :key="p.id ?? `legacy-${i}`"
          :src="p.url"
          :width="220"
          :height="220"
          :tilt="photoTiltAt(i)"
          :gem="!!review.is_gem && i === 0"
        >
          <Tape :variant="tapeVariantAt(i)" :style="tapeStyleAt(i)" />
          <span v-if="review.is_gem && i === 0" class="gem-corner">
            <GemBadge :size="22" />
          </span>
        </Polaroid>
      </div>
    </div>

    <!-- Ratings ticket -->
    <div class="sb-review-ratings sb-t-r1">
      <Ticket
        :food="review.food_rating"
        :service="review.service_rating"
        :vibe="review.vibe_rating"
      />
      <Stamp v-if="review.is_gem" kind="gem" class="gem-stamp">жемчужина</Stamp>
    </div>

    <!-- Video kruzhok (real player) -->
    <div v-if="review.video_url" class="sb-review-video">
      <div class="kruzhok-wrap" @click="toggleVideo">
        <video
          ref="videoEl"
          :src="review.video_url"
          class="kruzhok-player"
          playsinline
          loop
          muted
          preload="metadata"
        ></video>
        <div v-if="!playing" class="kruzhok-overlay">▶</div>
      </div>
    </div>

    <p v-if="review.comment" class="sb-review-text">{{ review.comment }}</p>

    <div class="sb-tear" aria-hidden="true"></div>
  </article>
</template>

<script setup>
import { ref, computed } from 'vue'
import Polaroid from './scrapbook/Polaroid.vue'
import Tape from './scrapbook/Tape.vue'
import Ticket from './scrapbook/Ticket.vue'
import Stamp from './scrapbook/Stamp.vue'
import GemBadge from './scrapbook/GemBadge.vue'
import AuthorTag from './scrapbook/AuthorTag.vue'
import { authorColor } from '../composables/useFeed'

const props = defineProps({
  review:    { type: Object, required: true },
  canEdit:   { type: Boolean, default: false },
  showPlace: { type: Boolean, default: false },
})

defineEmits(['edit', 'delete'])

const videoEl = ref(null)
const playing = ref(false)

// Photo list нормализуется: при наличии photos[] (новый API) — берём оттуда,
// иначе legacy single image_url. Так старые отзывы продолжают рендериться.
const photoList = computed(() => {
  const photos = props.review.photos || []
  if (photos.length > 0) return photos
  if (props.review.image_url) return [{ url: props.review.image_url }]
  return []
})

// Per-photo tilt / tape variant — deterministic by review id + index, чтобы
// в горизонтальном ряду полароиды смотрелись нерегулярно, но стабильно.
const tilts = ['t-l3', 't-r2', 't-l2', 't-r3', 't-l1']
function photoTiltAt(i) {
  return tilts[((props.review.id ?? 0) + i) % tilts.length]
}

const tapeVariants = ['', 'rose', 'mint', 'blue']
function tapeVariantAt(i) {
  return tapeVariants[((props.review.id ?? 0) + i) % tapeVariants.length]
}

const tapeStyleVariants = [
  { top: '-10px', left: '38px', transform: 'rotate(-12deg)' },
  { top: '-9px',  left: '52px', transform: 'rotate(6deg)' },
  { top: '-9px',  right: '24px', transform: 'rotate(8deg)' },
]
function tapeStyleAt(i) {
  return tapeStyleVariants[((props.review.id ?? 0) + i) % tapeStyleVariants.length]
}

function toggleVideo() {
  if (!videoEl.value) return
  if (videoEl.value.paused) {
    videoEl.value.muted = false
    videoEl.value.play()
    playing.value = true
  } else {
    videoEl.value.pause()
    playing.value = false
  }
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
}
</script>

<style scoped lang="scss">
.sb-review {
  background: #fdfcf7;
  margin: 0 0 22px;
  padding: 14px 14px 18px;
  position: relative;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 14px rgba(40, 30, 20, 0.07);
  border-radius: 1px;
}

.place-link {
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
  color: var(--sb-ink);
  margin-bottom: 6px;
  padding-bottom: 8px;
  border-bottom: 1px dashed rgba(40, 30, 20, 0.18);

  .place-name {
    font-family: var(--sb-serif);
    font-style: italic;
    font-weight: 500;
    font-size: 18px;
    line-height: 1.15;
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .arr {
    font-family: var(--sb-hand);
    font-size: 18px;
    color: var(--sb-ink-mute);
  }

  &:hover .place-name { color: var(--sb-terracotta); }
}

.sb-review-head {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: start;
  gap: 4px 8px;
  margin-bottom: 12px;
}

.authors {
  grid-column: 1;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  font-family: var(--sb-serif);
  font-size: 14px;
  color: var(--sb-ink);
}
.author-line {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
  color: inherit;
  &:hover .username { color: var(--sb-terracotta); }
}
.username { font-weight: 500; }
.sep { color: var(--sb-ink-mute); }
.joint {
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
}

/* AuthorTag here is laid out inline, not absolute — flatten its absolute positioning. */
.author-line :deep(.sb-author-tag) {
  position: relative;
  width: 22px;
  height: 22px;
}

.date {
  grid-column: 1;
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-soft);
}

.actions {
  grid-column: 2;
  grid-row: 1 / span 2;
  display: flex;
  gap: 6px;
  align-self: start;
}
.link-btn {
  background: transparent;
  border: none;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 13px;
  color: var(--sb-ink-mute);
  cursor: pointer;
  padding: 4px 6px;
  &:hover { color: var(--sb-ink); }
  &.danger:hover { color: var(--sb-terracotta); }
}

.sb-review-photos {
  margin: 6px 0 14px;
  position: relative;
}
.photos-row {
  display: flex;
  align-items: center;
  gap: 18px;
  overflow-x: auto;
  /* Запас сверху и снизу для tape-полосок и тени наклонённых полароидов. */
  padding: 14px 4px 18px;
  scroll-snap-type: x mandatory;
  -webkit-overflow-scrolling: touch;

  /* Скрываем родной скроллбар — скрапбук-вид. */
  scrollbar-width: none;
  &::-webkit-scrollbar { display: none; }

  & > :deep(.sb-polaroid) {
    flex: 0 0 auto;
    scroll-snap-align: start;
  }
}
.gem-corner {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
}

.sb-review-ratings {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 0 12px 4px;
  flex-wrap: wrap;
}
.gem-stamp { transform: rotate(-2deg); }

.sb-review-video {
  display: flex;
  margin-bottom: 12px;
}
.kruzhok-wrap {
  width: 132px;
  height: 132px;
  border-radius: 50%;
  overflow: hidden;
  position: relative;
  background: #000;
  cursor: pointer;
  box-shadow:
    0 0 0 3px #fdfcf7,
    0 2px 4px rgba(40, 30, 20, 0.18),
    0 8px 18px rgba(40, 30, 20, 0.18);
}
.kruzhok-player {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.kruzhok-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.3);
  color: #fff;
  font-size: 2rem;
}

.sb-review-text {
  font-family: var(--sb-serif);
  font-size: 15px;
  line-height: 1.45;
  color: var(--sb-ink);
  margin: 8px 0 4px;
}

.sb-tear {
  margin-top: 12px;
}
</style>
