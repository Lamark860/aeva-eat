<template>
  <div class="sb-paper sb-grain sb-screen place-form-wrap">
    <div class="form-shell">
      <div class="form-header">
        <h1 class="title">{{ headerTitle }}</h1>
        <p class="sub">{{ headerSub }}</p>
      </div>

      <div v-if="error" class="alert alert-danger">{{ error }}</div>

      <form @submit.prevent="handleSubmit">
        <!-- Карта (поиск заведения) — ПЕРВЫМ -->
        <div class="form-section">
          <label class="form-label fw-medium">Найти на Яндекс Картах</label>
          <LocationPicker
            :modelValue="{ lat: form.lat, lng: form.lng }"
            :city="form.city"
            @update:model-value="onLocationChange"
            @address-found="onAddressFound"
            @place-found="onPlaceFound"
          />
        </div>

        <!-- Автозаполненные данные — показываем если нашли -->
        <div v-if="form.name" class="form-section">
          <div class="found-place-card">
            <div class="d-flex align-items-start gap-2">
              <span class="found-icon">📍</span>
              <div class="flex-grow-1">
                <div class="fw-semibold">{{ form.name }}</div>
                <div class="text-muted small">
                  <span v-if="form.city">{{ form.city }}</span>
                  <span v-if="form.city && form.address">, </span>
                  <span v-if="form.address">{{ form.address }}</span>
                </div>
                <div v-if="foundCategories.length" class="mt-1">
                  <span v-for="cat in foundCategories" :key="cat" class="badge bg-light text-dark border me-1">{{ cat }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Ручной ввод — сворачиваемый -->
        <div class="form-section">
          <button
            type="button"
            class="btn btn-link text-decoration-none p-0 mb-2 d-flex align-items-center gap-1"
            @click="showManual = !showManual"
          >
            <span class="manual-arrow" :class="{ open: showManual }">▸</span>
            <span>{{ showManual ? 'Скрыть ручной ввод' : 'Ввести вручную' }}</span>
          </button>
          <div v-if="showManual" class="manual-fields">
            <div class="mb-3">
              <label class="form-label fw-medium">Название *</label>
              <input v-model="form.name" type="text" class="form-control" placeholder="Как называется заведение" required />
            </div>
            <div class="row g-3 mb-3">
              <div class="col-md-5">
                <label class="form-label fw-medium">Город</label>
                <input v-model="form.city" type="text" class="form-control" placeholder="Казань, Москва..." />
              </div>
              <div class="col-md-7">
                <label class="form-label fw-medium">Адрес</label>
                <input v-model="form.address" type="text" class="form-control" placeholder="ул. Баумана, 1" />
              </div>
            </div>
          </div>
        </div>

        <!-- Тип кухни + Категории -->
        <div class="form-section">
          <div class="row g-3">
            <div class="col-md-6">
              <label class="form-label fw-medium">Тип кухни</label>
              <MultiSelect
                :modelValue="form.cuisine_type_ids"
                @update:model-value="form.cuisine_type_ids = $event"
                :options="catalogs.cuisineTypes"
                placeholder="Выберите тип кухни..."
              />
            </div>
            <div class="col-md-6">
              <label class="form-label fw-medium">Категории</label>
              <MultiSelect
                :modelValue="form.category_ids"
                @update:model-value="form.category_ids = $event"
                :options="catalogs.categories"
                placeholder="Выберите категории..."
              />
            </div>
          </div>
        </div>

        <!-- Фото -->
        <div class="form-section">
          <label class="form-label fw-medium">Фото заведения</label>
          <div v-if="currentImageUrl" class="mb-2 position-relative d-inline-block">
            <img :src="currentImageUrl" alt="Фото" class="rounded" style="max-height: 200px; object-fit: cover" />
          </div>
          <div class="photo-upload-area" @click="$refs.fileInput.click()" @dragover.prevent @drop.prevent="onFileDrop">
            <input ref="fileInput" type="file" class="d-none" accept="image/jpeg,image/png,image/webp" @change="handleImageSelect" />
            <div v-if="selectedImage" class="text-center">
              <img :src="previewUrl" class="rounded mb-2" style="max-height: 150px; object-fit: cover" />
              <p class="small text-muted mb-0">{{ selectedImage.name }}</p>
            </div>
            <div v-else class="text-center text-muted">
              <div style="font-size: 2rem">📷</div>
              <p class="small mb-0">Нажмите или перетащите фото сюда</p>
              <p class="small text-muted mb-0">JPEG, PNG, WebP до 5 МБ</p>
            </div>
          </div>
          <button
            v-if="selectedImage && isEdit"
            type="button"
            class="btn btn-outline-success btn-sm mt-2"
            :disabled="uploadingImage"
            @click="handleImageUpload"
          >
            <span v-if="uploadingImage" class="spinner-border spinner-border-sm me-1"></span>
            Загрузить фото
          </button>
          <p v-if="selectedImage && !isEdit" class="small text-muted mt-1">
            Фото будет загружено после создания заведения
          </p>
        </div>

        <!-- Хочу сходить -->
        <div v-if="!isEdit" class="form-section">
          <div class="form-check">
            <input id="addToWishlist" v-model="addToWishlist" type="checkbox" class="form-check-input" />
            <label for="addToWishlist" class="form-check-label">🤍 Хочу сходить</label>
          </div>
        </div>

        <!-- Кнопки — стек вертикально на мобиле, чтобы primary CTA был полной ширины -->
        <div class="form-cta">
          <button type="submit" class="btn btn-apply" :disabled="loading || !form.name">
            <span v-if="loading" class="spinner-border spinner-border-sm me-1"></span>
            {{ isEdit ? 'сохранить' : 'создать' }}
          </button>
          <router-link to="/places" class="cancel-link">отмена</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePlacesStore } from '../stores/places'
import { useCatalogsStore } from '../stores/catalogs'
import { useWishlistStore } from '../stores/wishlist'
import { useToast } from '../composables/useToast'
import LocationPicker from '../components/LocationPicker.vue'
import MultiSelect from '../components/MultiSelect.vue'

const route = useRoute()
const router = useRouter()
const placesStore = usePlacesStore()
const catalogs = useCatalogsStore()
const wishlist = useWishlistStore()
const toast = useToast()

const isEdit = computed(() => !!route.params.id)
const intent = computed(() => route.query.intent)

const headerTitle = computed(() => {
  if (isEdit.value) return 'Редактировать место'
  if (intent.value === 'visit') return 'Новый визит'
  return 'Новое место'
})
const headerSub = computed(() => {
  if (isEdit.value) return 'правки сохраняются мгновенно'
  if (intent.value === 'visit') return 'сначала находим место — потом оценишь'
  return 'найдите на карте или введите вручную'
})
const error = ref('')
const loading = ref(false)
const selectedImage = ref(null)
const previewUrl = ref(null)
const uploadingImage = ref(false)
const currentImageUrl = ref('')
const showManual = ref(false)
const addToWishlist = ref(false)
const foundCategories = ref([])

const form = ref({
  name: '',
  city: '',
  address: '',
  cuisine_type_ids: [],
  category_ids: [],
  lat: null,
  lng: null
})

onMounted(async () => {
  await catalogs.fetchAll()
  if (isEdit.value) {
    showManual.value = true
    const place = await placesStore.fetchPlace(route.params.id)
    form.value = {
      name: place.name,
      city: place.city || '',
      address: place.address || '',
      cuisine_type_ids: place.cuisine_type_id ? [place.cuisine_type_id] : [],
      category_ids: place.category_ids || [],
      lat: place.lat || null,
      lng: place.lng || null
    }
    currentImageUrl.value = place.image_url || ''
  }
})

function onLocationChange(loc) {
  form.value.lat = loc.lat
  form.value.lng = loc.lng
}

function onAddressFound(info) {
  // Fallback for click-on-map (no place-found)
  if (!form.value.address && info.address) {
    const parts = info.address.split(', ')
    if (!form.value.city && parts.length > 1) {
      form.value.city = parts.length >= 3 ? parts[parts.length - 3] : parts[parts.length - 2]
    }
  }
}

function onPlaceFound(info) {
  if (info.name) form.value.name = info.name
  if (info.city) form.value.city = info.city
  if (info.address) form.value.address = info.address
  if (info.lat) form.value.lat = info.lat
  if (info.lng) form.value.lng = info.lng
  foundCategories.value = info.categories || []
}

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    const data = {
      name: form.value.name,
      address: form.value.address || undefined,
      city: form.value.city || undefined,
      cuisine_type_id: form.value.cuisine_type_ids.length ? form.value.cuisine_type_ids[0] : undefined,
      category_ids: form.value.category_ids,
      lat: form.value.lat || undefined,
      lng: form.value.lng || undefined
    }

    if (isEdit.value) {
      await placesStore.updatePlace(route.params.id, data)
      toast.success('Заведение обновлено')
      router.push(`/places/${route.params.id}`)
    } else {
      const created = await placesStore.createPlace(data)
      if (selectedImage.value) {
        try {
          await placesStore.uploadImage(created.id, selectedImage.value)
        } catch { /* photo upload is best-effort */ }
      }
      if (addToWishlist.value) {
        try {
          await wishlist.toggle(created.id)
        } catch { /* wishlist is best-effort */ }
      }
      toast.success('Заведение создано!')
      // If user came in via "+ добавить → новый визит", forward them to the
      // detail screen with `review=open` so the rating form auto-expands.
      // Otherwise (wishlist-style add or direct visit) just land on the place.
      const intent = route.query.intent
      if (intent === 'visit') {
        router.push({ path: `/places/${created.id}`, query: { review: 'open' } })
      } else {
        router.push(`/places/${created.id}`)
      }
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка при сохранении'
  } finally {
    loading.value = false
  }
}

function handleImageSelect(e) {
  const file = e.target.files[0] || null
  setSelectedFile(file)
}

function onFileDrop(e) {
  const file = e.dataTransfer.files[0]
  if (file && file.type.startsWith('image/')) {
    setSelectedFile(file)
  }
}

function setSelectedFile(file) {
  selectedImage.value = file
  if (file) {
    previewUrl.value = URL.createObjectURL(file)
  } else {
    previewUrl.value = null
  }
}

async function handleImageUpload() {
  if (!selectedImage.value) return
  uploadingImage.value = true
  try {
    const result = await placesStore.uploadImage(route.params.id, selectedImage.value)
    currentImageUrl.value = result.image_url
    selectedImage.value = null
    previewUrl.value = null
    toast.success('Фото загружено')
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка загрузки фото'
  } finally {
    uploadingImage.value = false
  }
}
</script>

<style scoped lang="scss">
.place-form-wrap {
  padding-top: calc(18px + var(--aeva-safe-top, 0px));
  padding-inline: 16px;
}

.form-shell {
  max-width: 640px;
  margin: 0 auto;
}

.form-header {
  margin-bottom: 18px;
  text-align: center;

  .title {
    font-family: var(--sb-serif);
    font-style: italic;
    font-weight: 500;
    font-size: 26px;
    color: var(--sb-ink);
    margin: 0 0 4px;
  }
  .sub {
    font-family: var(--sb-hand);
    font-size: 16px;
    color: var(--sb-ink-mute);
    margin: 0;
  }
}

.form-section {
  background: #fdfcf7;
  padding: 14px;
  border-radius: 1px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 14px rgba(40, 30, 20, 0.07);
  margin-bottom: 14px;
}

/* Bootstrap form-label override on this screen — handwritten subdued */
.form-section :deep(.form-label) {
  font-family: var(--sb-hand);
  font-size: 16px;
  font-weight: 400;
  color: var(--sb-ink-mute);
  margin-bottom: 4px;
}

.found-place-card {
  background: oklch(0.94 0.05 85 / 0.6);
  border: 1.4px dashed rgba(40, 30, 20, 0.25);
  border-radius: 4px;
  padding: 10px 12px;
}
.found-icon { font-size: 1.3rem; line-height: 1; }
.manual-arrow {
  display: inline-block;
  transition: transform 0.2s;
  font-size: 0.85rem;
}
.manual-arrow.open { transform: rotate(90deg); }
.manual-fields {
  border-left: 2px dashed rgba(40, 30, 20, 0.2);
  padding-left: 1rem;
}
.photo-upload-area {
  border: 2px dashed rgba(40, 30, 20, 0.25);
  border-radius: 4px;
  padding: 1.5rem;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
  background: oklch(0.97 0.018 82);
}
.photo-upload-area:hover {
  border-color: var(--sb-terracotta);
  background: oklch(0.95 0.04 85 / 0.5);
}

.form-cta {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 8px 0 24px;
}
.btn-apply {
  background: var(--sb-terracotta);
  color: #fff;
  border: none;
  border-radius: 999px;
  padding: 12px 22px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 15px;
  cursor: pointer;
  flex: 1;
  &:disabled { opacity: 0.5; cursor: not-allowed; }
  &:hover:not(:disabled) { background: oklch(0.55 0.14 30); color: #fff; }
}
.cancel-link {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 15px;
  color: var(--sb-ink-mute);
  text-decoration: none;
  padding: 12px 8px;
  &:hover { color: var(--sb-ink); }
}
</style>
