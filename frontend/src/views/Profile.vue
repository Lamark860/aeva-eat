<template>
  <div>
    <h2 class="mb-4">Мой профиль</h2>

    <div v-if="auth.user" class="card mb-4">
      <div class="card-body">
        <div class="d-flex align-items-center gap-3 mb-3">
          <div class="avatar-wrap" @click="$refs.avatarInput.click()">
            <img v-if="auth.user.avatar_url" :src="auth.user.avatar_url" class="avatar-img" />
            <div v-else class="avatar-placeholder">{{ auth.user.username?.charAt(0).toUpperCase() }}</div>
            <div class="avatar-overlay">📷</div>
            <input ref="avatarInput" type="file" accept="image/*" hidden @change="onAvatarChange" />
          </div>
          <div>
            <h5 class="mb-0">{{ auth.user.username }}</h5>
            <small class="text-muted">Зарегистрирован: {{ new Date(auth.user.created_at).toLocaleDateString('ru') }}</small>
          </div>
        </div>

        <!-- Смена пароля -->
        <button class="btn btn-sm btn-outline-secondary" @click="showPasswordForm = !showPasswordForm">
          {{ showPasswordForm ? 'Отмена' : '🔒 Сменить пароль' }}
        </button>
        <form v-if="showPasswordForm" class="mt-3" @submit.prevent="onChangePassword" style="max-width: 360px;">
          <div class="mb-2">
            <input v-model="pwForm.old_password" type="password" class="form-control form-control-sm" placeholder="Текущий пароль" required />
          </div>
          <div class="mb-2">
            <input v-model="pwForm.new_password" type="password" class="form-control form-control-sm" placeholder="Новый пароль (мин. 6 символов)" required minlength="6" />
          </div>
          <button type="submit" class="btn btn-primary btn-sm" :disabled="pwLoading">
            {{ pwLoading ? 'Сохранение...' : 'Сохранить' }}
          </button>
        </form>
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

      <!-- Inline edit form -->
      <ReviewForm
        v-if="editingReview"
        :place-id="editingReview.place_id"
        :review="editingReview"
        @submitted="handleUpdateReview"
        @cancel="editingReview = null"
      />

      <div v-for="rv in displayedReviews" :key="rv.id" class="card mb-3">
        <div v-if="rv.image_url" class="review-img-wrap-sm">
          <img :src="rv.image_url" class="review-img-sm" />
        </div>
        <div class="card-body">
          <div class="d-flex justify-content-between align-items-start">
            <div>
              <router-link :to="`/places/${rv.place_id}`" class="fw-bold text-decoration-none">
                {{ rv.place_name || `Заведение #${rv.place_id}` }}
              </router-link>
              <span v-if="rv.visited_at" class="text-muted ms-2">· {{ formatDate(rv.visited_at) }}</span>
            </div>
            <div class="d-flex gap-2 align-items-center">
              <span class="badge bg-warning text-dark">🍴 {{ rv.food_rating }}</span>
              <span class="badge bg-info text-dark">🤝 {{ rv.service_rating }}</span>
              <span class="badge bg-success">✨ {{ rv.vibe_rating }}</span>
              <span v-if="rv.is_gem" class="badge bg-primary">💎</span>
              <button class="btn btn-outline-primary btn-sm ms-1" @click="editingReview = rv" title="Редактировать">✏️</button>
              <button class="btn btn-outline-danger btn-sm" @click="handleDeleteReview(rv)" title="Удалить">🗑</button>
            </div>
          </div>
          <p v-if="rv.comment" class="mb-0 mt-2 text-muted">{{ rv.comment }}</p>
          <div v-if="rv.video_url" class="mt-2">
            <div class="video-circle-profile">
              <video :src="rv.video_url" class="video-preview-profile" playsinline loop muted preload="metadata" @click="toggleVideo($event)"></video>
            </div>
          </div>
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
import { useToast } from '../composables/useToast'
import PlaceCard from '../components/PlaceCard.vue'
import ReviewForm from '../components/ReviewForm.vue'
import http from '../api/http'

const auth = useAuthStore()
const reviewsStore = useReviewsStore()
const wishlist = useWishlistStore()
const toast = useToast()

const customName = ref('')
const customNote = ref('')
const reviewsExpanded = ref(false)
const editingReview = ref(null)

// Password change
const showPasswordForm = ref(false)
const pwLoading = ref(false)
const pwForm = ref({ old_password: '', new_password: '' })

async function onChangePassword() {
  pwLoading.value = true
  try {
    await http.put('/auth/password', pwForm.value)
    toast.success('Пароль успешно изменён')
    showPasswordForm.value = false
    pwForm.value = { old_password: '', new_password: '' }
  } catch (e) {
    toast.error(e.response?.data?.error || 'Ошибка при смене пароля')
  } finally {
    pwLoading.value = false
  }
}

// Avatar upload
async function onAvatarChange(e) {
  const file = e.target.files?.[0]
  if (!file) return
  const formData = new FormData()
  formData.append('avatar', file)
  try {
    const { data } = await http.post('/auth/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    auth.user.avatar_url = data.avatar_url
    toast.success('Аватар обновлён')
  } catch (err) {
    toast.error(err.response?.data?.error || 'Ошибка загрузки аватара')
  }
}

const displayedReviews = computed(() => {
  const list = reviewsStore.reviews.filter(r => editingReview.value ? r.id !== editingReview.value.id : true)
  if (reviewsExpanded.value) return list
  return list.slice(0, 5)
})

async function handleUpdateReview(data) {
  const photoFile = data._photoFile
  const videoFile = data._videoFile
  delete data._photoFile
  delete data._videoFile
  const placeId = editingReview.value.place_id
  const reviewId = editingReview.value.id
  try {
    await reviewsStore.updateReview(placeId, reviewId, data)
    if (photoFile) {
      await reviewsStore.uploadReviewImage(placeId, reviewId, photoFile)
    }
    if (videoFile) {
      await reviewsStore.uploadReviewVideo(placeId, reviewId, videoFile)
    }
    editingReview.value = null
    toast.success('Отзыв обновлён')
  } catch (e) {
    toast.error(e.response?.data?.error || 'Ошибка при обновлении отзыва')
  }
  if (auth.user) await reviewsStore.fetchByUser(auth.user.id)
}

async function handleDeleteReview(rv) {
  if (!confirm(`Удалить отзыв к "${rv.place_name || 'заведению'}"?`)) return
  try {
    await reviewsStore.deleteReview(rv.place_id, rv.id)
    toast.info('Отзыв удалён')
  } catch (e) {
    toast.error(e.response?.data?.error || 'Ошибка при удалении отзыва')
  }
  if (auth.user) await reviewsStore.fetchByUser(auth.user.id)
}

function toggleVideo(e) {
  const v = e.target
  if (v.paused) { v.muted = false; v.play() } else { v.pause() }
}

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

<style scoped>
.avatar-wrap {
  position: relative;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
  flex-shrink: 0;
}
.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #6c757d;
  color: #fff;
  font-size: 1.5rem;
  font-weight: 600;
}
.avatar-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0,0,0,0.4);
  opacity: 0;
  transition: opacity 0.2s;
  font-size: 1.2rem;
}
.avatar-wrap:hover .avatar-overlay {
  opacity: 1;
}
.review-img-wrap-sm {
  overflow: hidden;
  border-radius: 0.75rem 0.75rem 0 0;
}
.review-img-sm {
  width: 100%;
  max-height: 150px;
  object-fit: cover;
}
.video-circle-profile {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  background: #000;
  cursor: pointer;
  border: 2px solid #dee2e6;
}
.video-preview-profile {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
