<template>
  <div>
    <!-- Hero -->
    <section class="text-center py-5 mb-4">
      <h1 class="display-4 fw-bold mb-3">
        <span class="text-gradient">AEVA Eat</span>
      </h1>
      <p class="lead text-muted mx-auto" style="max-width: 500px">
        Храни впечатления о заведениях, делись с друзьями, находи жемчужины на карте.
      </p>
      <div class="d-flex gap-2 justify-content-center mt-4">
        <template v-if="auth.isAuthenticated">
          <router-link to="/places/new" class="btn btn-primary btn-lg">+ Добавить заведение</router-link>
          <router-link to="/map" class="btn btn-outline-secondary btn-lg">Открыть карту</router-link>
        </template>
        <template v-else>
          <router-link to="/register" class="btn btn-primary btn-lg">Начать</router-link>
          <router-link to="/login" class="btn btn-outline-secondary btn-lg">Уже есть аккаунт</router-link>
        </template>
      </div>
    </section>

    <!-- Stats strip -->
    <section v-if="stats" class="row g-3 mb-5">
      <div class="col-4">
        <div class="card text-center py-3">
          <div class="fs-2 fw-bold text-gradient">{{ stats.placeCount }}</div>
          <small class="text-muted">Заведений</small>
        </div>
      </div>
      <div class="col-4">
        <div class="card text-center py-3">
          <div class="fs-2 fw-bold text-gradient">{{ stats.reviewCount }}</div>
          <small class="text-muted">Отзывов</small>
        </div>
      </div>
      <div class="col-4">
        <div class="card text-center py-3">
          <div class="fs-2 fw-bold text-gradient">{{ stats.gemCount }}</div>
          <small class="text-muted">💎 Жемчужин</small>
        </div>
      </div>
    </section>

    <!-- Recent places -->
    <section v-if="recentPlaces.length">
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h4 class="mb-0">Последние заведения</h4>
        <router-link to="/places" class="text-decoration-none">Все →</router-link>
      </div>
      <div class="row g-3">
        <div v-for="place in recentPlaces" :key="place.id" class="col-12 col-md-6 col-lg-4">
          <PlaceCard :place="place" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import PlaceCard from '../components/PlaceCard.vue'
import http from '../api/http'

const auth = useAuthStore()
const recentPlaces = ref([])
const stats = ref(null)

onMounted(async () => {
  try {
    const { data } = await http.get('/places', { params: { sort: '', limit: 6 } })
    recentPlaces.value = data || []

    // Calculate basic stats from places list
    const allPlaces = (await http.get('/places')).data || []
    const totalReviews = allPlaces.reduce((sum, p) => sum + (p.review_count || 0), 0)
    const gems = allPlaces.filter(p => p.is_gem_place).length
    stats.value = {
      placeCount: allPlaces.length,
      reviewCount: totalReviews,
      gemCount: gems
    }
  } catch { /* ignore */ }
})
</script>
