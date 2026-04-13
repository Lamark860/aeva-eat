<template>
  <div v-if="placesStore.loading" class="text-center py-5">
    <div class="spinner-border"></div>
  </div>

  <div v-else-if="place">
    <div v-if="place.image_url" class="mb-4">
      <img :src="place.image_url" :alt="place.name" class="img-fluid rounded" style="max-height: 350px; width: 100%; object-fit: cover" />
    </div>

    <div class="d-flex justify-content-between align-items-start mb-3">
      <div>
        <h2>{{ place.name }}</h2>
        <p class="text-muted mb-1" v-if="place.cuisine_type">🍽 {{ place.cuisine_type }}</p>
        <p class="text-muted mb-1" v-if="place.city">
          📍 {{ place.city }}<span v-if="place.address">, {{ place.address }}</span>
        </p>
        <p class="mb-1" v-if="place.website">
          🌐 <a :href="place.website" target="_blank">{{ place.website }}</a>
        </p>
        <div v-if="place.categories && place.categories.length" class="mt-2">
          <span v-for="cat in place.categories" :key="cat" class="badge bg-light text-dark me-1">{{ cat }}</span>
        </div>
      </div>

      <div v-if="isOwner" class="d-flex gap-2">
        <router-link :to="`/places/${place.id}/edit`" class="btn btn-outline-primary btn-sm">Ред.</router-link>
        <button class="btn btn-outline-danger btn-sm" @click="handleDelete">Удалить</button>
      </div>
    </div>

    <!-- Ratings summary -->
    <div v-if="place.review_count > 0" class="card mb-4">
      <div class="card-body">
        <div class="row text-center">
          <div class="col">
            <div class="fs-3 fw-bold text-warning">{{ avgRound(place.avg_food) }}</div>
            <small class="text-muted">Кухня</small>
          </div>
          <div class="col">
            <div class="fs-3 fw-bold text-info">{{ avgRound(place.avg_service) }}</div>
            <small class="text-muted">Сервис</small>
          </div>
          <div class="col">
            <div class="fs-3 fw-bold text-success">{{ avgRound(place.avg_vibe) }}</div>
            <small class="text-muted">Вайб</small>
          </div>
          <div class="col">
            <div class="fs-3 fw-bold">{{ place.review_count }}</div>
            <small class="text-muted">Отзывов</small>
          </div>
        </div>
      </div>
    </div>

    <!-- Mini map -->
    <div v-if="place.lat && place.lng" class="mb-4">
      <MapView
        :places="[place]"
        :center="[place.lat, place.lng]"
        :zoom="15"
        height="250px"
        :single-marker="true"
      />
    </div>

    <!-- Reviews section -->
    <h4 class="mt-4 mb-3">Отзывы</h4>

    <ReviewForm
      v-if="auth.isAuthenticated && !editingReview"
      :place-id="place.id"
      @submitted="handleCreateReview"
    />

    <ReviewForm
      v-if="editingReview"
      :place-id="place.id"
      :review="editingReview"
      @submitted="handleUpdateReview"
      @cancel="editingReview = null"
    />

    <div v-if="reviewsStore.loading" class="text-center py-3">
      <div class="spinner-border spinner-border-sm"></div>
    </div>
    <div v-else-if="reviewsStore.reviews.length === 0" class="text-muted text-center py-3">
      Отзывов пока нет. Будьте первым!
    </div>
    <div v-else>
      <ReviewCard
        v-for="rv in reviewsStore.reviews"
        :key="rv.id"
        :review="rv"
        :can-edit="canEditReview(rv)"
        @edit="editingReview = rv"
        @delete="handleDeleteReview"
      />
    </div>

    <router-link to="/places" class="btn btn-outline-secondary mt-3">← Назад к списку</router-link>
  </div>

  <div v-else class="text-center py-5 text-muted">
    Заведение не найдено
  </div>
</template>

<script setup>
import { onMounted, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePlacesStore } from '../stores/places'
import { useReviewsStore } from '../stores/reviews'
import { useAuthStore } from '../stores/auth'
import ReviewCard from '../components/ReviewCard.vue'
import ReviewForm from '../components/ReviewForm.vue'
import MapView from '../components/MapView.vue'

const route = useRoute()
const router = useRouter()
const placesStore = usePlacesStore()
const reviewsStore = useReviewsStore()
const auth = useAuthStore()

const editingReview = ref(null)

const place = computed(() => placesStore.currentPlace)
const isOwner = computed(() => auth.user && place.value && place.value.created_by === auth.user.id)

function avgRound(val) {
  return val ? Number(val).toFixed(1) : '–'
}

function canEditReview(rv) {
  if (!auth.user) return false
  return rv.authors?.some(a => a.id === auth.user.id)
}

async function handleCreateReview(data) {
  await reviewsStore.createReview(place.value.id, data)
  await placesStore.fetchPlace(route.params.id)
}

async function handleUpdateReview(data) {
  await reviewsStore.updateReview(place.value.id, editingReview.value.id, data)
  editingReview.value = null
  await placesStore.fetchPlace(route.params.id)
}

async function handleDeleteReview(id) {
  if (!confirm('Удалить отзыв?')) return
  await reviewsStore.deleteReview(place.value.id, id)
  await placesStore.fetchPlace(route.params.id)
}

async function handleDelete() {
  if (!confirm('Удалить заведение?')) return
  await placesStore.deletePlace(place.value.id)
  router.push('/places')
}

onMounted(async () => {
  await placesStore.fetchPlace(route.params.id)
  await reviewsStore.fetchByPlace(route.params.id)
})
</script>
