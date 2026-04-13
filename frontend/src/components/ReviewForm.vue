<template>
  <div class="card mb-4">
    <div class="card-header">
      <h5 class="mb-0">{{ isEdit ? 'Редактировать отзыв' : 'Оставить отзыв' }}</h5>
    </div>
    <div class="card-body">
      <div v-if="error" class="alert alert-danger">{{ error }}</div>

      <form @submit.prevent="handleSubmit">
        <RatingInput v-model="form.food_rating" label="🍴 Кухня" color="#f59e0b" class="mb-2" />
        <RatingInput v-model="form.service_rating" label="🤝 Сервис" color="#3b82f6" class="mb-2" />
        <RatingInput v-model="form.vibe_rating" label="✨ Вайб" color="#2eb872" class="mb-2" />

        <div v-if="totalRating !== null" class="d-flex align-items-center gap-2 mb-3 px-1">
          <span class="text-muted small">Общий рейтинг:</span>
          <span class="fw-bold fs-5" :style="{ color: totalColor }">{{ totalRating }}</span>
          <span class="text-muted small">/ 10</span>
        </div>

        <div class="mb-3">
          <div class="form-check">
            <input v-model="form.is_gem" class="form-check-input" type="checkbox" id="gemCheck" />
            <label class="form-check-label" for="gemCheck">💎 Жемчужина</label>
          </div>
        </div>

        <div class="mb-3">
          <label class="form-label">Комментарий</label>
          <textarea v-model="form.comment" class="form-control" rows="3" placeholder="Что понравилось / не понравилось?"></textarea>
        </div>

        <div class="mb-3">
          <label class="form-label">Дата визита</label>
          <input v-model="form.visited_at" type="date" class="form-control" style="max-width: 200px;" />
        </div>

        <div class="mb-3">
          <label class="form-label">Фото</label>
          <div
            class="photo-drop-zone"
            :class="{ 'has-preview': photoPreview }"
            @dragover.prevent
            @drop.prevent="onPhotoDrop"
            @click="$refs.photoInput.click()"
          >
            <img v-if="photoPreview" :src="photoPreview" class="photo-preview" />
            <span v-else class="text-muted small">📷 Перетащите фото или нажмите</span>
          </div>
          <input ref="photoInput" type="file" accept="image/*" class="d-none" @change="onPhotoSelect" />
        </div>

        <div class="d-flex gap-2">
          <button type="submit" class="btn btn-primary" :disabled="loading || !isValid">
            <span v-if="loading" class="spinner-border spinner-border-sm me-1"></span>
            {{ isEdit ? 'Сохранить' : 'Отправить' }}
          </button>
          <button v-if="isEdit" type="button" class="btn btn-outline-secondary" @click="$emit('cancel')">
            Отмена
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import RatingInput from './RatingInput.vue'

const props = defineProps({
  review: { type: Object, default: null },
  placeId: { type: Number, required: true }
})

const emit = defineEmits(['submitted', 'cancel'])

const isEdit = computed(() => !!props.review)

const form = ref({
  food_rating: props.review?.food_rating || 0,
  service_rating: props.review?.service_rating || 0,
  vibe_rating: props.review?.vibe_rating || 0,
  is_gem: props.review?.is_gem || false,
  comment: props.review?.comment || '',
  visited_at: props.review?.visited_at || new Date().toISOString().split('T')[0]
})

const error = ref('')
const loading = ref(false)
const photoFile = ref(null)
const photoPreview = ref(null)

const isValid = computed(() =>
  form.value.food_rating >= 0 && form.value.service_rating >= 0 && form.value.vibe_rating >= 0
)

const totalRating = computed(() => {
  const { food_rating, service_rating, vibe_rating } = form.value
  if (food_rating === 0 && service_rating === 0 && vibe_rating === 0) return null
  return ((food_rating + service_rating + vibe_rating) / 3).toFixed(1)
})

const totalColor = computed(() => {
  const v = parseFloat(totalRating.value)
  if (v >= 8) return '#2eb872'
  if (v >= 5) return '#f59e0b'
  return '#ef4444'
})

function onPhotoSelect(e) {
  const file = e.target.files?.[0]
  if (file) { photoFile.value = file; photoPreview.value = URL.createObjectURL(file) }
}

function onPhotoDrop(e) {
  const file = e.dataTransfer.files?.[0]
  if (file && file.type.startsWith('image/')) { photoFile.value = file; photoPreview.value = URL.createObjectURL(file) }
}

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    emit('submitted', {
      ...form.value,
      comment: form.value.comment || undefined,
      visited_at: form.value.visited_at || undefined,
      _photoFile: photoFile.value || undefined
    })
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.photo-drop-zone {
  border: 2px dashed #ddd;
  border-radius: 0.5rem;
  padding: 1rem;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s;
  min-height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.photo-drop-zone:hover { border-color: var(--bs-primary); }
.photo-drop-zone.has-preview { padding: 0.25rem; }
.photo-preview {
  max-height: 150px;
  max-width: 100%;
  border-radius: 0.4rem;
  object-fit: cover;
}
</style>
