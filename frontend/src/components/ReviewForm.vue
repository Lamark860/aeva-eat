<template>
  <div class="sb-review-form">
    <header class="head">
      <h3 class="title">{{ isEdit ? 'Редактировать отзыв' : 'Оставить отзыв' }}</h3>
      <span class="sub">прикнопь впечатления</span>
    </header>

    <div v-if="error" class="alert alert-danger">{{ error }}</div>

    <form @submit.prevent="handleSubmit">
      <RatingInput v-model="form.food_rating" label="кухня" tone="ochre" />
      <RatingInput v-model="form.service_rating" label="сервис" tone="plum" />
      <RatingInput v-model="form.vibe_rating" label="вайб" tone="moss" />

      <div v-if="totalRating !== null" class="overall">
        <span class="o-lbl">общая</span>
        <span class="o-val">{{ totalRating }}</span>
        <span class="o-of">/ 10</span>
      </div>

      <div class="field">
        <button
          type="button"
          class="gem-toggle"
          :class="{ on: form.is_gem }"
          @click="form.is_gem = !form.is_gem"
        >
          <GemBadge v-if="form.is_gem" :size="22" class="gem-icon" />
          <span v-else class="gem-empty" aria-hidden="true">♢</span>
          <span class="gem-lbl">{{ form.is_gem ? 'жемчужина' : 'отметить как жемчужину' }}</span>
        </button>
      </div>

      <div class="field">
        <label class="lbl">комментарий</label>
        <textarea
          v-model="form.comment"
          class="form-control paper-control"
          rows="3"
          placeholder="что понравилось / не понравилось?"
        ></textarea>
      </div>

      <div class="field">
        <label class="lbl">дата визита</label>
        <input v-model="form.visited_at" type="date" class="form-control paper-control date-input" />
      </div>

      <div class="field">
        <label class="lbl">фото <span class="lbl-meta">до {{ MAX_PHOTOS }}</span></label>
        <div class="photo-stack">
          <!-- Уже сохранённые фото отзыва (при edit) -->
          <div
            v-for="ph in existingPhotos"
            :key="`ex-${ph.id}`"
            class="photo-slot filled"
            :class="slotTilt(ph.id)"
          >
            <img :src="ph.url" class="photo-img" alt="" />
            <button
              type="button"
              class="photo-x"
              :aria-label="'удалить фото'"
              @click.stop="removeExistingPhoto(ph)"
            >
              ×
            </button>
            <span class="photo-cap">фото</span>
          </div>

          <!-- Новые превью (ещё не загружены на бэк) -->
          <div
            v-for="(p, i) in newPreviews"
            :key="`new-${i}`"
            class="photo-slot filled"
            :class="slotTilt(i + 100)"
          >
            <img :src="p.url" class="photo-img" alt="" />
            <button
              type="button"
              class="photo-x"
              :aria-label="'удалить фото'"
              @click.stop="removeNewPhoto(i)"
            >
              ×
            </button>
            <span class="photo-cap">фото</span>
          </div>

          <!-- Следующий пустой слот — лист в стопке. Не показываем после лимита. -->
          <div
            v-if="totalPhotos < MAX_PHOTOS"
            class="photo-slot empty"
            :class="slotTilt(totalPhotos + 200)"
            @dragover.prevent
            @drop.prevent="onPhotoDrop"
            @click="$refs.photoInput.click()"
          >
            <div class="dashed-corners"></div>
            <span class="hint">{{ totalPhotos === 0 ? 'тапни фото' : '+ ещё' }}</span>
          </div>
        </div>
        <input ref="photoInput" type="file" accept="image/*" multiple hidden @change="onPhotoSelect" />
      </div>

      <div class="field">
        <label class="lbl">видеосообщение</label>
        <div v-if="existingVideoUrl && !videoFile" class="existing-video">
          <div class="kruzhok-frame">
            <video :src="existingVideoUrl" class="kruzhok-vid" playsinline loop muted autoplay></video>
          </div>
          <button type="button" class="cancel-link" @click="existingVideoUrl = null">удалить</button>
        </div>
        <VideoRecorder v-else :uploading="loading" @recorded="onVideoRecorded" />
      </div>

      <div class="cta">
        <button type="submit" class="btn-apply" :disabled="loading || !isValid">
          <span v-if="loading" class="spinner-border spinner-border-sm me-1"></span>
          {{ isEdit ? 'сохранить' : 'отправить' }}
        </button>
        <button v-if="isEdit" type="button" class="cancel-link" @click="$emit('cancel')">отмена</button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import RatingInput from './RatingInput.vue'
import VideoRecorder from './VideoRecorder.vue'
import GemBadge from './scrapbook/GemBadge.vue'

const props = defineProps({
  review:  { type: Object, default: null },
  placeId: { type: Number, required: true },
})

const emit = defineEmits(['submitted', 'cancel'])

const isEdit = computed(() => !!props.review)

const form = ref({
  food_rating:    props.review?.food_rating    || 0,
  service_rating: props.review?.service_rating || 0,
  vibe_rating:    props.review?.vibe_rating    || 0,
  is_gem:         props.review?.is_gem         || false,
  comment:        props.review?.comment        || '',
  visited_at:     props.review?.visited_at     || new Date().toISOString().split('T')[0],
})

const error = ref('')
const loading = ref(false)

const MAX_PHOTOS = 5
// Фото уже на бэке (edit-flow): { id, url }. На submit удалённые из этого
// списка дозаявляются через store.deleteReviewPhoto.
const existingPhotos = ref([...(props.review?.photos || [])])
const removedPhotoIds = ref([])
// Новые превью — { url: blob-url, file: File }. На submit улетают одним
// multipart-запросом через uploadReviewPhotos.
const newPreviews = ref([])

const totalPhotos = computed(() => existingPhotos.value.length + newPreviews.value.length)

const videoFile = ref(null)
const existingVideoUrl = ref(props.review?.video_url || null)

const isValid = computed(
  () => form.value.food_rating >= 0 && form.value.service_rating >= 0 && form.value.vibe_rating >= 0,
)

const totalRating = computed(() => {
  const { food_rating, service_rating, vibe_rating } = form.value
  if (food_rating === 0 && service_rating === 0 && vibe_rating === 0) return null
  return ((food_rating + service_rating + vibe_rating) / 3).toFixed(1)
})

function addFiles(files) {
  for (const file of files) {
    if (totalPhotos.value >= MAX_PHOTOS) break
    if (!file.type.startsWith('image/')) continue
    newPreviews.value.push({ url: URL.createObjectURL(file), file })
  }
}

function onPhotoSelect(e) {
  const files = Array.from(e.target.files || [])
  addFiles(files)
  // Сбрасываем input, чтобы повторный выбор того же файла снова сработал.
  e.target.value = ''
}

function onPhotoDrop(e) {
  const files = Array.from(e.dataTransfer.files || [])
  addFiles(files)
}

function removeNewPhoto(idx) {
  const [removed] = newPreviews.value.splice(idx, 1)
  if (removed?.url) URL.revokeObjectURL(removed.url)
}

function removeExistingPhoto(ph) {
  existingPhotos.value = existingPhotos.value.filter(p => p.id !== ph.id)
  removedPhotoIds.value.push(ph.id)
}

// Детерминистический tilt по ключу — стопка не «прыгает» между рендерами.
function slotTilt(key) {
  const variants = ['tilt-l', 'tilt-r', 'tilt-l2', 'tilt-r2']
  return variants[Math.abs(key) % variants.length]
}

function onVideoRecorded(blob) {
  videoFile.value = blob
}

async function handleSubmit() {
  error.value = ''
  loading.value = true
  emit('submitted', {
    ...form.value,
    comment:    form.value.comment || undefined,
    visited_at: form.value.visited_at || undefined,
    _newPhotoFiles:    newPreviews.value.map(p => p.file),
    _removedPhotoIds:  [...removedPhotoIds.value],
    _videoFile:        videoFile.value || undefined,
  })
}
</script>

<style scoped lang="scss">
.sb-review-form {
  /* PlaceDetail wraps this in its own .detail-form paper card; Profile keeps it inline.
     We don't repeat the box-shadow here to avoid stacked shadows. */
  font-family: var(--sb-serif);
}

.head {
  margin-bottom: 14px;
  .title {
    font-family: var(--sb-serif);
    font-style: italic;
    font-weight: 500;
    font-size: 19px;
    color: var(--sb-ink);
    margin: 0;
  }
  .sub {
    font-family: var(--sb-hand);
    font-size: 15px;
    color: var(--sb-ink-mute);
  }
}

.alert-danger {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  background: oklch(0.92 0.07 25);
  color: var(--sb-terracotta);
  border: 1px dashed oklch(0.55 0.14 30 / 0.5);
  border-radius: 4px;
  padding: 8px 12px;
  margin-bottom: 12px;
}

.overall {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin: 8px 0 14px;
  padding: 6px 4px 0;

  .o-lbl {
    font-family: var(--sb-serif);
    font-style: italic;
    font-size: 14px;
    color: var(--sb-ink-mute);
  }
  .o-val {
    font-family: var(--sb-hand);
    font-size: 28px;
    color: var(--sb-terracotta);
    line-height: 1;
  }
  .o-of {
    font-family: var(--sb-hand);
    font-size: 16px;
    color: var(--sb-ink-mute);
  }
}

.field {
  margin-bottom: 14px;
}

.lbl {
  display: block;
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  margin-bottom: 4px;
}

/* Override Bootstrap's form-control on this form — paper background */
.paper-control {
  font-family: var(--sb-serif);
  font-size: 16px;
  color: var(--sb-ink);
  background: oklch(0.97 0.018 82);
  border: 1.4px solid rgba(40, 30, 20, 0.18);
  border-radius: 3px;
  padding: 8px 10px;
  box-shadow: inset 0 1px 2px rgba(40, 30, 20, 0.04);

  &:focus {
    border-color: var(--sb-terracotta);
    box-shadow: 0 0 0 2px oklch(0.55 0.14 30 / 0.15);
    background: #fdfcf7;
  }
  &::placeholder {
    color: var(--sb-ink-mute);
    font-style: italic;
  }
}

textarea.paper-control {
  resize: vertical;
  min-height: 72px;
}

.date-input { max-width: 200px; }
@media (max-width: 767.98px) {
  .date-input { max-width: 100%; }
}

/* Gem toggle — stamp-styled chip */
.gem-toggle {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: transparent;
  border: 1.4px dashed rgba(40, 30, 20, 0.3);
  border-radius: 999px;
  padding: 6px 14px 6px 10px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink-mute);
  cursor: pointer;
  min-height: 36px;

  .gem-icon { display: inline-flex; }
  .gem-empty {
    color: var(--sb-ink-mute);
    font-size: 18px;
    line-height: 1;
  }
  .gem-lbl { line-height: 1; }

  &.on {
    background: oklch(0.94 0.05 85 / 0.7);
    border-style: solid;
    border-color: var(--sb-terracotta);
    color: var(--sb-terracotta);
    // DESIGN-DECISIONS R3: «прихлопнул печатью» — короткий thump на
    // включении. Каждый клик в `on` перезапускает анимацию.
    animation: sb-gem-stamp 220ms ease-out;
  }
}
@keyframes sb-gem-stamp {
  0%   { transform: scale(1); }
  35%  { transform: scale(0.85); }
  70%  { transform: scale(1.05); }
  100% { transform: scale(1); }
}

.lbl-meta {
  font-family: var(--sb-hand);
  font-size: 13px;
  color: var(--sb-ink-soft);
  margin-left: 4px;
}

/* Photo stack — стопка полароидов «следующий лист» (DESIGN-DECISIONS §R1).
   Каждый слот — мини-полароид с лёгким наклоном. Заполненные показывают
   фото и крестик-удаление; пустой — пунктирные уголки и подпись. */
.photo-stack {
  display: flex;
  flex-wrap: wrap;
  gap: 14px 18px;
  padding: 6px 4px 12px;
}

.photo-slot {
  position: relative;
  flex: 0 0 auto;
  width: 120px;
  height: 138px;          /* квадрат фото + caption-полоса снизу */
  background: #fdfcf7;
  padding: 6px 6px 22px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.08),
    0 3px 8px rgba(40, 30, 20, 0.10);
  border-radius: 1px;
  cursor: pointer;
  transition: transform 180ms ease;

  &.tilt-l  { transform: rotate(-2.5deg); }
  &.tilt-r  { transform: rotate(2.5deg); }
  &.tilt-l2 { transform: rotate(-1.2deg); }
  &.tilt-r2 { transform: rotate(1.5deg); }

  &.empty {
    background: rgba(253, 252, 247, 0.5);
    box-shadow:
      0 1px 1px rgba(40, 30, 20, 0.04),
      0 2px 5px rgba(40, 30, 20, 0.06);

    .dashed-corners {
      position: absolute;
      inset: 6px 6px 22px;
      border: 1.5px dashed rgba(40, 30, 20, 0.26);
      border-radius: 2px;
    }
    .hint {
      position: absolute;
      inset: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      font-family: var(--sb-hand);
      font-size: 16px;
      color: var(--sb-ink-mute);
      padding-bottom: 16px;
      pointer-events: none;
    }
  }
}

.photo-img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 1px;
}

.photo-cap {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 4px;
  text-align: center;
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-soft);
  line-height: 1;
  pointer-events: none;
}

.photo-x {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: var(--sb-ink);
  color: #fdfcf7;
  border: 2px solid #fdfcf7;
  font-family: var(--sb-serif);
  font-size: 14px;
  line-height: 1;
  cursor: pointer;
  z-index: 2;
  box-shadow: 0 2px 4px rgba(40, 30, 20, 0.2);

  &:hover { background: var(--sb-terracotta); }
}

.existing-video {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;
}
.kruzhok-frame {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  background: #000;
  box-shadow:
    0 0 0 3px #fdfcf7,
    0 2px 4px rgba(40, 30, 20, 0.18),
    0 8px 18px rgba(40, 30, 20, 0.18);
}
.kruzhok-vid {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cta {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 8px;
}
.btn-apply {
  background: var(--sb-terracotta);
  color: #fff;
  border: none;
  border-radius: 999px;
  padding: 12px 22px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 15px;
  cursor: pointer;
  flex: 1;
  &:disabled { opacity: 0.5; cursor: not-allowed; }
  &:hover:not(:disabled) { background: oklch(0.55 0.14 30); color: #fff; }
}
.cancel-link {
  background: transparent;
  border: none;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 15px;
  color: var(--sb-ink-mute);
  cursor: pointer;
  padding: 12px 8px;
  &:hover { color: var(--sb-ink); }
}
</style>
