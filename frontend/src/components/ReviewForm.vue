<template>
  <div class="card mb-4">
    <div class="card-header">
      <h5 class="mb-0">{{ isEdit ? 'Редактировать отзыв' : 'Оставить отзыв' }}</h5>
    </div>
    <div class="card-body">
      <div v-if="error" class="alert alert-danger">{{ error }}</div>

      <form @submit.prevent="handleSubmit">
        <RatingInput v-model="form.food_rating" label="Кухня" color="warning" class="mb-2" />
        <RatingInput v-model="form.service_rating" label="Сервис" color="info" class="mb-2" />
        <RatingInput v-model="form.vibe_rating" label="Вайб" color="success" class="mb-3" />

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

const isValid = computed(() =>
  form.value.food_rating >= 1 && form.value.service_rating >= 1 && form.value.vibe_rating >= 1
)

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    emit('submitted', {
      ...form.value,
      comment: form.value.comment || undefined,
      visited_at: form.value.visited_at || undefined
    })
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка'
  } finally {
    loading.value = false
  }
}
</script>
