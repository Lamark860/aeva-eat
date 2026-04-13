<template>
  <div class="card mb-3">
    <div v-if="review.image_url" class="review-img-wrap">
      <img :src="review.image_url" :alt="'Фото отзыва'" class="review-img" />
    </div>
    <div class="card-body">
      <div class="d-flex justify-content-between align-items-start">
        <div>
          <div class="d-flex gap-2 align-items-center mb-1">
            <span class="fw-bold fs-5" style="color: var(--bs-primary)">{{ overallRating }}</span>
            <span class="badge bg-warning text-dark" title="Кухня">🍴 {{ review.food_rating }}</span>
            <span class="badge bg-info text-dark" title="Сервис">🤝 {{ review.service_rating }}</span>
            <span class="badge bg-success" title="Вайб">✨ {{ review.vibe_rating }}</span>
            <span v-if="review.is_gem" class="badge bg-primary" title="Жемчужина">💎</span>
          </div>
          <div class="d-flex align-items-center gap-1 text-muted small mb-1">
            <span
              v-for="(author, i) in review.authors"
              :key="author.id"
              class="d-inline-flex align-items-center gap-1"
            >
              <span v-if="i > 0" class="me-1">+</span>
              <span class="review-author-avatar">{{ author.username.charAt(0).toUpperCase() }}</span>
              <span>{{ author.username }}</span>
            </span>
            <span v-if="review.authors && review.authors.length > 1" class="badge bg-light text-muted ms-1">совместный</span>
            <span v-if="review.visited_at"> · {{ formatDate(review.visited_at) }}</span>
          </div>
        </div>
        <div v-if="canEdit" class="d-flex gap-1">
          <button class="btn btn-outline-primary btn-sm" @click="$emit('edit', review)">✏️</button>
          <button class="btn btn-outline-danger btn-sm" @click="$emit('delete', review.id)">🗑</button>
        </div>
      </div>
      <p v-if="review.comment" class="mb-0 mt-2">{{ review.comment }}</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  review: { type: Object, required: true },
  canEdit: { type: Boolean, default: false }
})

defineEmits(['edit', 'delete'])

const overallRating = computed(() => {
  const r = props.review
  return ((r.food_rating + r.service_rating + r.vibe_rating) / 3).toFixed(1)
})

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
}
</script>

<style scoped>
.review-img-wrap {
  overflow: hidden;
  border-radius: 0.75rem 0.75rem 0 0;
}
.review-img {
  width: 100%;
  max-height: 250px;
  object-fit: cover;
}
.review-author-avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--bs-primary);
  color: #fff;
  font-size: 0.6rem;
  font-weight: 700;
}
</style>
