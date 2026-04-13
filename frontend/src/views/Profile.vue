<template>
  <div>
    <h2 class="mb-4">Мой профиль</h2>

    <div v-if="auth.user" class="card mb-4">
      <div class="card-body">
        <h5>{{ auth.user.username }}</h5>
        <small class="text-muted">Зарегистрирован: {{ new Date(auth.user.created_at).toLocaleDateString('ru') }}</small>
      </div>
    </div>

    <h4 class="mb-3">Мои отзывы <span v-if="reviewsStore.reviews.length" class="text-muted fs-6">({{ reviewsStore.reviews.length }})</span></h4>

    <div v-if="reviewsStore.loading" class="text-center py-3">
      <div class="spinner-border spinner-border-sm"></div>
    </div>

    <div v-else-if="reviewsStore.reviews.length === 0" class="text-muted text-center py-3">
      У вас пока нет отзывов.
    </div>

    <div v-else>
      <button class="btn btn-sm btn-outline-secondary mb-3" @click="reviewsExpanded = !reviewsExpanded">
        {{ reviewsExpanded ? 'Свернуть' : 'Показать все' }}
      </button>
      <div v-for="(rv, idx) in displayedReviews" :key="rv.id" class="card mb-3">
        <div class="card-body">
          <div class="d-flex justify-content-between">
            <div>
              <router-link :to="`/places/${rv.place_id}`" class="fw-bold text-decoration-none">
                {{ rv.place_name || `Заведение #${rv.place_id}` }}
              </router-link>
              <span v-if="rv.visited_at" class="text-muted ms-2">· {{ formatDate(rv.visited_at) }}</span>
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

    <!-- Wishlist -->
    <h4 class="mt-4 mb-3">Хочу посетить</h4>

    <div v-if="wishlist.loading" class="text-center py-3">
      <div class="spinner-border spinner-border-sm"></div>
    </div>

    <div v-else-if="wishlist.wishlistPlaces.length === 0" class="text-muted text-center py-3">
      Пока ничего не запланировано. Нажмите 🤍 на карточке заведения.
    </div>

    <div v-else class="row g-3 mb-4">
      <div v-for="place in wishlist.wishlistPlaces" :key="place.id" class="col-12 col-md-6 col-lg-4">
        <PlaceCard :place="place" />
      </div>
    </div>

    <!-- Custom wishlist -->
    <h4 class="mt-4 mb-3">Ещё хочется попробовать</h4>

    <form @submit.prevent="onAddCustom" class="d-flex gap-2 mb-3">
      <input
        v-model="customName"
        type="text"
        class="form-control form-control-sm"
        placeholder="Название заведения..."
        required
      />
      <input
        v-model="customNote"
        type="text"
        class="form-control form-control-sm"
        placeholder="Заметка (опционально)"
      />
      <button type="submit" class="btn btn-primary btn-sm text-nowrap" :disabled="!customName.trim()">+ Добавить</button>
    </form>

    <div v-if="wishlist.customItems.length === 0" class="text-muted text-center py-3">
      Здесь можно записать заведения, которые хотите посетить, но их ещё нет в системе.
    </div>

    <div v-else class="list-group mb-4">
      <div v-for="item in wishlist.customItems" :key="item.id" class="list-group-item d-flex justify-content-between align-items-center">
        <div>
          <strong>{{ item.name }}</strong>
          <small v-if="item.note" class="text-muted ms-2">— {{ item.note }}</small>
        </div>
        <button class="btn btn-outline-danger btn-sm" @click="onDeleteCustom(item.id)" title="Удалить">✕</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useReviewsStore } from '../stores/reviews'
import { useWishlistStore } from '../stores/wishlist'
import PlaceCard from '../components/PlaceCard.vue'

const auth = useAuthStore()
const reviewsStore = useReviewsStore()
const wishlist = useWishlistStore()

const customName = ref('')
const customNote = ref('')
const reviewsExpanded = ref(false)

const displayedReviews = computed(() => {
  if (reviewsExpanded.value) return reviewsStore.reviews
  return reviewsStore.reviews.slice(0, 5)
})

async function onAddCustom() {
  if (!customName.value.trim()) return
  await wishlist.addCustom(customName.value.trim(), customNote.value.trim() || null)
  customName.value = ''
  customNote.value = ''
}

async function onDeleteCustom(id) {
  await wishlist.deleteCustom(id)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(() => {
  if (auth.user) {
    reviewsStore.fetchByUser(auth.user.id)
    wishlist.fetchPlaces()
    wishlist.fetchCustom()
  }
})
</script>
