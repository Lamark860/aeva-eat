<template>
  <router-link :to="`/places/${place.id}`" class="text-decoration-none">
    <div class="card card-hover h-100">
      <div class="card-img-wrap">
        <img
          v-if="place.image_url"
          :src="place.image_url"
          :alt="place.name"
          class="card-img-top"
        />
        <div v-else class="card-img-placeholder">
          <span>🍽</span>
        </div>
        <span v-if="place.is_gem_place" class="gem-badge">💎</span>
        <button
          v-if="auth.isAuthenticated"
          class="wishlist-btn"
          :class="{ active: wishlist.isWishlisted(place.id) }"
          @click.prevent="onToggleWishlist"
          :title="wishlist.isWishlisted(place.id) ? 'Убрать из планов' : 'Хочу сходить'"
        >
          {{ wishlist.isWishlisted(place.id) ? '❤️' : '🤍' }}
        </button>
      </div>
      <div class="card-body d-flex flex-column">
        <h5 class="card-title mb-1">{{ place.name }}</h5>
        <p class="text-muted small mb-1" v-if="place.cuisine_type">
          {{ place.cuisine_type }}
        </p>
        <p class="text-muted small mb-2" v-if="place.city">
          📍 {{ place.city }}<span v-if="place.address">, {{ place.address }}</span>
        </p>

        <div class="mt-auto">
          <div v-if="place.review_count > 0">
            <div class="d-flex align-items-center gap-2 mb-1">
              <span class="overall-rating">{{ overallRating }}</span>
              <span class="text-muted small ms-auto">{{ place.review_count }} отз.</span>
            </div>
            <div class="d-flex gap-2 align-items-center">
              <div class="rating-bar">
                <span class="rating-item food">🍴 {{ avgRound(place.avg_food) }}</span>
                <span class="rating-item service">🤝 {{ avgRound(place.avg_service) }}</span>
                <span class="rating-item vibe">✨ {{ avgRound(place.avg_vibe) }}</span>
              </div>
            </div>
            <div v-if="place.reviewers && place.reviewers.length" class="reviewers-row mt-2">
              <div class="reviewers-stack" :title="reviewerNames">
                <span
                  v-for="(rev, i) in place.reviewers.slice(0, 4)"
                  :key="rev.id"
                  class="reviewer-avatar"
                  :style="{ zIndex: 10 - i, marginLeft: i > 0 ? '-8px' : '0' }"
                >{{ rev.username.charAt(0).toUpperCase() }}</span>
                <span v-if="place.reviewers.length > 4" class="reviewer-avatar reviewer-more" :style="{ marginLeft: '-8px' }">+{{ place.reviewers.length - 4 }}</span>
              </div>
            </div>
          </div>
          <div v-else class="text-muted small">
            Ещё нет отзывов
          </div>
        </div>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useWishlistStore } from '../stores/wishlist'

const props = defineProps({
  place: { type: Object, required: true }
})

const auth = useAuthStore()
const wishlist = useWishlistStore()

const overallRating = computed(() => {
  const p = props.place
  if (!p.avg_food || !p.avg_service || !p.avg_vibe) return '–'
  return ((Number(p.avg_food) + Number(p.avg_service) + Number(p.avg_vibe)) / 3).toFixed(1)
})

const reviewerNames = computed(() => {
  if (!props.place.reviewers) return ''
  return props.place.reviewers.map(r => r.username).join(', ')
})

function avgRound(val) {
  return val ? Number(val).toFixed(1) : '–'
}

function onToggleWishlist() {
  wishlist.toggle(props.place.id)
}
</script>

<style scoped>
.card-img-wrap {
  position: relative;
  overflow: hidden;
  border-radius: 0.75rem 0.75rem 0 0;
}

.card-img-top {
  height: 180px;
  object-fit: cover;
  width: 100%;
  transition: transform 0.3s ease;
}

.card-hover:hover .card-img-top {
  transform: scale(1.05);
}

.card-img-placeholder {
  height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #fef3ef, #fde8df);
  font-size: 3rem;
}

.gem-badge {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
}

.wishlist-btn {
  position: absolute;
  top: 0.5rem;
  left: 0.5rem;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  cursor: pointer;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  transition: transform 0.2s;
  z-index: 2;
}
.wishlist-btn:hover {
  transform: scale(1.15);
}

.rating-bar {
  display: flex;
  gap: 0.35rem;
}

.rating-item {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.2em 0.5em;
  border-radius: 0.4rem;
  font-family: 'Inter', sans-serif;
  letter-spacing: -0.01em;
}
.rating-item.food    { background: #fff3cd; color: #856404; }
.rating-item.service { background: #d1ecf1; color: #0c5460; }
.rating-item.vibe    { background: #d4edda; color: #155724; }

.card-title {
  color: var(--bs-body-color);
}

.overall-rating {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--bs-primary);
  line-height: 1;
  font-family: 'Inter', sans-serif;
  letter-spacing: -0.02em;
}

.reviewers-row {
  display: flex;
  align-items: center;
}
.reviewers-stack {
  display: flex;
  align-items: center;
  cursor: default;
}
.reviewer-avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--bs-primary);
  color: #fff;
  font-size: 0.65rem;
  font-weight: 700;
  border: 2px solid #fff;
  font-family: 'Inter', sans-serif;
}
.reviewer-more {
  background: #6c757d;
  font-size: 0.6rem;
}
</style>
