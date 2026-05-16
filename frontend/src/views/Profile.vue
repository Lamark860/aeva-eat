<template>
  <div class="sb-paper sb-grain sb-screen me">
    <header class="me-header">
      <div class="avatar-wrap" @click="$refs.avatarInput.click()">
        <img v-if="auth.user?.avatar_url" :src="auth.user.avatar_url" class="avatar-img" alt="" />
        <div v-else class="avatar-placeholder">{{ initial }}</div>
        <div class="avatar-overlay">фото</div>
        <input ref="avatarInput" type="file" accept="image/*" hidden @change="onAvatarChange" />
      </div>

      <div class="me-name">{{ auth.user?.username || '—' }}</div>
      <div v-if="auth.user?.created_at" class="me-since">
        в кругу с&nbsp;{{ formatJoinDate(auth.user.created_at) }}
      </div>

      <!-- Stats ticket -->
      <div class="me-stats sb-t-l1">
        <Ticket
          :food="stats.places"
          :service="stats.gems"
          :vibe="stats.cities"
          :labels="['мест', 'жемчужин', 'городов']"
        />
      </div>

      <div v-if="cuisineLine" class="cuisine-line">{{ cuisineLine }}</div>
    </header>

    <!-- Tabs -->
    <nav class="me-tabs">
      <button
        v-for="t in tabs"
        :key="t.key"
        class="me-tab"
        :class="{ active: activeTab === t.key }"
        @click="activeTab = t.key"
      >
{{ t.label }}
</button>
    </nav>

    <!-- Visits -->
    <section v-if="activeTab === 'visits'" class="me-section">
      <div v-if="reviewsStore.loading" class="sb-empty">…</div>
      <div v-else-if="reviewsStore.reviews.length === 0" class="sb-empty">
        <div>пусто</div>
        <div style="margin-top: 8px">прикнопь первый визит —</div>
        <div style="margin-top: 12px">
          <router-link to="/places/new" class="add-link">+ добавить место</router-link>
        </div>
      </div>
      <div v-else>
        <ReviewForm
          v-if="editingReview"
          :place-id="editingReview.place_id"
          :review="editingReview"
          @submitted="handleUpdateReview"
          @cancel="editingReview = null"
        />

        <ReviewCard
          v-for="rv in displayedReviews"
          :key="rv.id"
          :review="rv"
          can-edit
          show-place
          @edit="editingReview = rv"
          @delete="(id) => handleDeleteReview(rv)"
        />

        <div v-if="reviewsStore.reviews.length > 5" class="me-more">
          <button class="more-link" @click="reviewsExpanded = !reviewsExpanded">
            {{ reviewsExpanded ? '↑ свернуть' : `↓ ещё ${reviewsStore.reviews.length - 5}` }}
          </button>
        </div>
      </div>
    </section>

    <!-- Wishlist -->
    <section v-if="activeTab === 'wishlist'" class="me-section">
      <div class="sb-section-head" style="padding: 0 0 8px">
        <h2>Из круга</h2>
        <span class="sub">места из системы</span>
      </div>

      <div v-if="wishlist.loading" class="sb-empty">…</div>
      <div v-else-if="wishlist.wishlistPlaces.length === 0" class="sb-empty">
        <div>ничего не запланировано</div>
        <div style="margin-top: 6px">тапни ♥ на любой карточке —</div>
      </div>
      <div v-else class="results-list">
        <ResultCard v-for="p in wishlist.wishlistPlaces" :key="p.id" :place="p" />
      </div>

      <div class="sb-section-head" style="padding: 18px 0 8px">
        <h2>Своё</h2>
        <span class="sub">записки от руки</span>
      </div>

      <form @submit.prevent="onAddCustom" class="custom-add">
        <input
          v-model="customName"
          type="text"
          class="form-control"
          placeholder="место (название)"
          required
        />
        <input
          v-model="customNote"
          type="text"
          class="form-control"
          placeholder="заметка (опционально)"
        />
        <button type="submit" class="add-link" :disabled="!customName.trim()">+ записать</button>
      </form>

      <div v-if="wishlist.customItems.length === 0" class="sb-empty" style="padding: 24px 12px">
        здесь можно записать места, которых ещё нет в системе
      </div>

      <div v-else class="custom-notes">
        <div
          v-for="(item, i) in wishlist.customItems"
          :key="item.id"
          class="custom-note"
          :class="customTilt(i)"
        >
          <Note lined>
            <div class="note-name">{{ item.name }}</div>
            <div v-if="item.note" class="note-text">{{ item.note }}</div>
          </Note>
          <Tape :variant="customTape(i)" :style="customTapeStyle(i)" />
          <button class="note-remove" type="button" @click="onDeleteCustom(item.id)" aria-label="Удалить">×</button>
        </div>
      </div>
    </section>

    <!-- Notes -->
    <section v-if="activeTab === 'notes'" class="me-section">
      <div class="sb-section-head" style="padding: 0 0 8px">
        <h2>Записки</h2>
        <span class="sub">мои бумажки на доске</span>
      </div>

      <button v-if="!noteFormOpenInline" class="add-link" @click="noteFormOpenInline = true">
        + новая записка
      </button>

      <form v-else class="custom-add" @submit.prevent="onCreateNote">
        <textarea
          v-model="newNoteText"
          class="form-control"
          rows="3"
          placeholder="что прикнопить?"
          required
        ></textarea>
        <div class="settings-row" style="margin-top: 8px">
          <button type="submit" class="btn-apply" :disabled="!newNoteText.trim()">прикнопить</button>
          <button type="button" class="reset-btn" @click="noteFormOpenInline = false">отмена</button>
        </div>
      </form>

      <div v-if="notesStore.myNotes.length === 0" class="sb-empty" style="padding: 24px 12px">
        пока ни одной записки
      </div>

      <div v-else class="custom-notes" style="margin-top: 18px">
        <div
          v-for="(n, i) in notesStore.myNotes"
          :key="n.id"
          class="custom-note"
          :class="customTilt(i)"
        >
          <NoteArtifact :note="n" can-edit @delete="onDeleteMyNote" />
        </div>
      </div>
    </section>

    <!-- Settings -->
    <section v-if="activeTab === 'settings'" class="me-section">
      <div class="settings-card">
        <h3 class="settings-h">Пароль</h3>
        <button v-if="!showPasswordForm" class="add-link" @click="showPasswordForm = true">сменить пароль</button>
        <form v-else class="settings-form" @submit.prevent="onChangePassword">
          <input v-model="pwForm.old_password" type="password" class="form-control" placeholder="текущий пароль" required />
          <input v-model="pwForm.new_password" type="password" class="form-control" placeholder="новый пароль (мин. 6)" required minlength="6" />
          <div class="settings-row">
            <button type="submit" class="btn-apply" :disabled="pwLoading">
              {{ pwLoading ? '…' : 'сохранить' }}
            </button>
            <button type="button" class="reset-btn" @click="showPasswordForm = false">отмена</button>
          </div>
        </form>
      </div>

      <div class="settings-card">
        <h3 class="settings-h">Инвайты</h3>
        <router-link to="/invites" class="add-link">пригласить друга →</router-link>
      </div>

      <div class="settings-card">
        <button class="danger-link" @click="logout">выйти</button>
      </div>
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useReviewsStore } from '../stores/reviews'
import { useWishlistStore } from '../stores/wishlist'
import { useToast } from '../composables/useToast'
import ReviewForm from '../components/ReviewForm.vue'
import ReviewCard from '../components/ReviewCard.vue'
import Note from '../components/scrapbook/Note.vue'
import Tape from '../components/scrapbook/Tape.vue'
import Ticket from '../components/scrapbook/Ticket.vue'
import ResultCard from '../components/scrapbook/ResultCard.vue'
import NoteArtifact from '../components/scrapbook/NoteArtifact.vue'
import { useNotesStore } from '../stores/notes'
import { favoriteCuisinePhrase } from '../composables/useCuisine'
import http from '../api/http'

const router = useRouter()
const auth = useAuthStore()
const reviewsStore = useReviewsStore()
const wishlist = useWishlistStore()
const notesStore = useNotesStore()
const toast = useToast()

const noteFormOpenInline = ref(false)
const newNoteText = ref('')

async function onCreateNote() {
  try {
    await notesStore.create({ text: newNoteText.value.trim() })
    newNoteText.value = ''
    noteFormOpenInline.value = false
    toast.success('Записка прикноплена')
  } catch (e) {
    toast.error(e.response?.data?.error || 'Не удалось сохранить')
  }
}

async function onDeleteMyNote(id) {
  if (!confirm('Удалить записку?')) return
  try {
    await notesStore.remove(id)
    toast.info('Записка убрана')
  } catch (e) {
    toast.error(e.response?.data?.error || 'Не удалось удалить')
  }
}

const customName = ref('')
const customNote = ref('')
const reviewsExpanded = ref(false)
const editingReview = ref(null)

const activeTab = ref('visits')
const tabs = [
  { key: 'visits',   label: 'Визиты' },
  { key: 'wishlist', label: 'Wishlist' },
  { key: 'notes',    label: 'Записки' },
  { key: 'settings', label: 'Настройки' },
]

const initial = computed(() => (auth.user?.username || '?').slice(0, 1).toUpperCase())

// Q8 — публичный профиль с favorite_cuisine. Грузим один раз на маунт.
const userProfile = ref(null)
const cuisineLine = computed(() => favoriteCuisinePhrase(userProfile.value))
async function loadUserProfile() {
  if (!auth.user?.id) return
  try {
    const { data } = await http.get(`/users/${auth.user.id}`)
    userProfile.value = data
  } catch { /* профиль необязателен — фразу просто не рендерим */ }
}

// B7 — ячейка «городов» в билетике должна быть числом, а не прочерком.
// Бэкенд считает city_count точно (через визиты), фронт берёт его если
// профиль успел подгрузиться. Fallback на локальный расчёт через
// wishlistPlaces оставлен на случай провала загрузки /users/:id.
const stats = computed(() => {
  const placeIds = new Set()
  let gems = 0
  for (const rv of reviewsStore.reviews) {
    placeIds.add(rv.place_id)
    if (rv.is_gem) gems++
  }
  let cities = userProfile.value?.city_count
  if (cities == null) {
    const localCities = new Set()
    for (const p of wishlist.wishlistPlaces) {
      if (placeIds.has(p.id) && p.city) localCities.add(p.city)
    }
    cities = localCities.size
  }
  return { places: placeIds.size, gems, cities: cities > 0 ? cities : '–' }
})

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

async function onAvatarChange(e) {
  const file = e.target.files?.[0]
  if (!file) return
  const formData = new FormData()
  formData.append('avatar', file)
  try {
    const { data } = await http.post('/auth/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    auth.user.avatar_url = data.avatar_url
    toast.success('Аватар обновлён')
  } catch (err) {
    toast.error(err.response?.data?.error || 'Ошибка загрузки аватара')
  }
}

const displayedReviews = computed(() => {
  const list = reviewsStore.reviews.filter((r) =>
    editingReview.value ? r.id !== editingReview.value.id : true,
  )
  if (reviewsExpanded.value) return list
  return list.slice(0, 5)
})

async function handleUpdateReview(data) {
  const newPhotoFiles   = data._newPhotoFiles || []
  const removedPhotoIds = data._removedPhotoIds || []
  const videoFile       = data._videoFile
  delete data._newPhotoFiles
  delete data._removedPhotoIds
  delete data._videoFile
  const placeId = editingReview.value.place_id
  const reviewId = editingReview.value.id
  try {
    await reviewsStore.updateReview(placeId, reviewId, data)
    for (const pid of removedPhotoIds) {
      await reviewsStore.deleteReviewPhoto(placeId, reviewId, pid)
    }
    if (newPhotoFiles.length > 0) {
      await reviewsStore.uploadReviewPhotos(placeId, reviewId, newPhotoFiles)
    }
    if (videoFile) await reviewsStore.uploadReviewVideo(placeId, reviewId, videoFile)
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

async function onAddCustom() {
  if (!customName.value.trim()) return
  await wishlist.addCustom(customName.value.trim(), customNote.value.trim() || null)
  customName.value = ''
  customNote.value = ''
}

async function onDeleteCustom(id) {
  await wishlist.deleteCustom(id)
}

const customTilts = ['sb-t-l3', 'sb-t-r2', 'sb-t-l2', 'sb-t-r3']
function customTilt(i) { return customTilts[i % customTilts.length] }
const customTapes = ['', 'rose', 'mint', 'blue']
function customTape(i) { return customTapes[i % customTapes.length] }
function customTapeStyle(i) {
  const variants = [
    { top: '-8px', left: '34px', transform: 'rotate(-12deg)' },
    { top: '-8px', right: '40px', transform: 'rotate(8deg)' },
    { top: '-9px', left: '60px', transform: 'rotate(-6deg)' },
  ]
  return variants[i % variants.length]
}

function formatJoinDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' })
}

function logout() {
  auth.logout()
  router.push('/login')
}

function loadAllForUser() {
  if (!auth.user?.id) return
  reviewsStore.fetchByUser(auth.user.id)
  wishlist.fetchPlaces()
  wishlist.fetchCustom()
  notesStore.fetchByAuthor(auth.user.id)
  loadUserProfile()
}

onMounted(loadAllForUser)
// Hard-reloads finish auth/init() async — user becomes available *after* mount.
// Watch fires once the JWT resolves to /auth/me and we can then hydrate data.
watch(() => auth.user?.id, (id) => { if (id) loadAllForUser() })
</script>

<style scoped lang="scss">
.me {
  padding-top: calc(18px + var(--aeva-safe-top, 0px));
  padding-inline: 16px;
}

.me-header {
  text-align: center;
  margin-bottom: 14px;
}

.avatar-wrap {
  position: relative;
  width: 96px;
  height: 96px;
  border-radius: 50%;
  overflow: hidden;
  margin: 0 auto 10px;
  cursor: pointer;
  box-shadow:
    0 0 0 4px var(--sb-paper-card),
    0 2px 4px rgba(40, 30, 20, 0.18),
    0 8px 18px rgba(40, 30, 20, 0.12);
}
.avatar-img,
.avatar-placeholder {
  width: 100%;
  height: 100%;
}
.avatar-img { object-fit: cover; }
.avatar-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: oklch(0.78 0.10 30);
  color: oklch(0.22 0.05 30);
  font-family: var(--sb-serif);
  font-weight: 600;
  font-size: 32px;
}
.avatar-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(40, 30, 20, 0.4);
  color: var(--sb-on-accent);
  font-family: var(--sb-hand);
  font-size: 16px;
  opacity: 0;
  transition: opacity 0.2s;
}
.avatar-wrap:hover .avatar-overlay { opacity: 1; }

.me-name {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 24px;
  color: var(--sb-ink);
  word-break: break-word;
  overflow-wrap: break-word;
}
.me-since {
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  margin-top: 2px;
}

.me-stats {
  display: inline-block;
  margin: 14px 0 6px;
}

.cuisine-line {
  margin-top: 6px;
  font-family: var(--sb-hand);
  font-size: 18px;
  color: var(--sb-ink-soft);
  line-height: 1.2;
}

.me-tabs {
  display: flex;
  gap: 6px;
  justify-content: center;
  margin: 18px 0 14px;
  flex-wrap: wrap;
}
.me-tab {
  background: transparent;
  border: 1.4px solid rgba(40, 30, 20, 0.18);
  border-radius: 999px;
  padding: 6px 14px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink-soft);
  cursor: pointer;
  white-space: nowrap;
  &.active {
    background: var(--sb-ink);
    color: var(--sb-paper);
    border-color: var(--sb-ink);
  }
}
@media (max-width: 400px) {
  .me-tabs { gap: 4px; }
  .me-tab {
    padding: 5px 10px;
    font-size: 13px;
  }
}

.me-section {
  padding-bottom: 6px;
}

.add-link {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-terracotta);
  text-decoration: none;
  background: transparent;
  border: none;
  cursor: pointer;
  &:disabled { color: var(--sb-ink-mute); cursor: not-allowed; }
  &:hover:not(:disabled) { color: oklch(0.5 0.14 30); }
}

.me-more {
  text-align: center;
  padding: 8px 0 18px;
}
.more-link {
  background: transparent;
  border: none;
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  cursor: pointer;
  padding: 6px 14px;
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.custom-add {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 14px;
}
.custom-add .form-control {
  flex: 1 1 140px;
  min-height: 40px;
  font-family: var(--sb-serif);
}

.custom-notes {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 22px 18px;
  margin: 12px 0;
}
.custom-note {
  position: relative;
  display: block;
}
.note-name {
  font-family: var(--sb-hand);
  font-size: 20px;
  color: var(--sb-ink);
  font-weight: 500;
  line-height: 1.05;
}
.note-text {
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  margin-top: 4px;
}
.note-remove {
  position: absolute;
  top: -6px;
  right: -6px;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: oklch(0.92 0.07 25);
  color: var(--sb-terracotta);
  border: none;
  font-size: 16px;
  line-height: 1;
  cursor: pointer;
  z-index: 4;
  box-shadow: 0 1px 2px rgba(40, 15, 5, 0.3);
  &:hover { background: var(--sb-terracotta); color: var(--sb-on-accent); }
}

.settings-card {
  background: var(--sb-paper-card);
  padding: 14px;
  border-radius: 1px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 14px rgba(40, 30, 20, 0.07);
  margin-bottom: 14px;
}
.settings-h {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 16px;
  color: var(--sb-ink);
  margin: 0 0 8px;
}
.settings-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.settings-row {
  display: flex;
  gap: 14px;
  align-items: center;
}
.btn-apply {
  background: var(--sb-terracotta);
  color: var(--sb-on-accent);
  border: none;
  border-radius: 999px;
  padding: 8px 16px;
  font-family: var(--sb-serif);
  font-style: italic;
  cursor: pointer;
  &:disabled { opacity: 0.5; cursor: not-allowed; }
}
.reset-btn {
  background: transparent;
  border: none;
  font-family: var(--sb-serif);
  font-style: italic;
  color: var(--sb-ink-mute);
  cursor: pointer;
}
.danger-link {
  background: transparent;
  border: none;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-terracotta);
  cursor: pointer;
  padding: 4px 0;
}
</style>
