<template>
  <div>
    <h2 class="mb-4">Мой профиль</h2>

    <div v-if="auth.user" class="card mb-4">
      <div class="card-body">
        <h5>{{ auth.user.username }}</h5>
        <p class="text-muted mb-0">{{ auth.user.email }}</p>
        <small class="text-muted">Зарегистрирован: {{ new Date(auth.user.created_at).toLocaleDateString('ru') }}</small>
      </div>
    </div>

    <h4 class="mb-3">Мои отзывы</h4>

    <div v-if="reviewsStore.loading" class="text-center py-3">
      <div class="spinner-border spinner-border-sm"></div>
    </div>

    <div v-else-if="reviewsStore.reviews.length === 0" class="text-muted text-center py-3">
      У вас пока нет отзывов.
    </div>

    <div v-else>
      <div v-for="rv in reviewsStore.reviews" :key="rv.id" class="card mb-3">
        <div class="card-body">
          <div class="d-flex justify-content-between">
            <div>
              <router-link :to="`/places/${rv.place_id}`" class="fw-bold text-decoration-none">
                Заведение #{{ rv.place_id }}
              </router-link>
              <span v-if="rv.visited_at" class="text-muted ms-2">· {{ rv.visited_at }}</span>
            </div>
            <div class="d-flex gap-2">
              <span class="badge bg-warning text-dark">🍴 {{ rv.food_rating }}</span>
              <span class="badge bg-info text-dark">🤝 {{ rv.service_rating }}</span>
              <span class="badge bg-success">✨ {{ rv.vibe_rating }}</span>
              <span v-if="rv.is_gem" class="badge bg-primary">💎</span>
            </div>
          </div>
          <p v-if="rv.comment" class="mb-0 mt-2 text-muted">{{ rv.comment }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useReviewsStore } from '../stores/reviews'

const auth = useAuthStore()
const reviewsStore = useReviewsStore()

onMounted(() => {
  if (auth.user) {
    reviewsStore.fetchByUser(auth.user.id)
  }
})
</script>
