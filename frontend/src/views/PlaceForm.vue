<template>
  <div class="row justify-content-center">
    <div class="col-md-7">
      <h3 class="mb-4">{{ isEdit ? 'Редактировать заведение' : 'Новое заведение' }}</h3>

      <div v-if="error" class="alert alert-danger">{{ error }}</div>

      <form @submit.prevent="handleSubmit">
        <div class="mb-3">
          <label class="form-label">Название *</label>
          <input v-model="form.name" type="text" class="form-control" required />
        </div>

        <div class="row g-3 mb-3">
          <div class="col-md-6">
            <label class="form-label">Город</label>
            <input v-model="form.city" type="text" class="form-control" placeholder="Москва" />
          </div>
          <div class="col-md-6">
            <label class="form-label">Адрес</label>
            <input v-model="form.address" type="text" class="form-control" placeholder="ул. Примерная, 1" />
          </div>
        </div>

        <div class="row g-3 mb-3">
          <div class="col-md-6">
            <label class="form-label">Тип кухни</label>
            <select v-model="form.cuisine_type_id" class="form-select">
              <option :value="null">Не указан</option>
              <option v-for="ct in catalogs.cuisineTypes" :key="ct.id" :value="ct.id">{{ ct.name }}</option>
            </select>
          </div>
          <div class="col-md-6">
            <label class="form-label">Сайт</label>
            <input v-model="form.website" type="url" class="form-control" placeholder="https://..." />
          </div>
        </div>

        <div class="mb-3">
          <label class="form-label">Категории</label>
          <div class="d-flex flex-wrap gap-2">
            <div v-for="cat in catalogs.categories" :key="cat.id" class="form-check">
              <input
                v-model="form.category_ids"
                :value="cat.id"
                class="form-check-input"
                type="checkbox"
                :id="'cat-' + cat.id"
              />
              <label class="form-check-label" :for="'cat-' + cat.id">{{ cat.name }}</label>
            </div>
          </div>
        </div>

        <div class="row g-3 mb-3">
          <div class="col-md-6">
            <label class="form-label">Широта</label>
            <input v-model.number="form.lat" type="number" step="any" class="form-control" placeholder="55.7558" />
          </div>
          <div class="col-md-6">
            <label class="form-label">Долгота</label>
            <input v-model.number="form.lng" type="number" step="any" class="form-control" placeholder="37.6173" />
          </div>
        </div>

        <div class="d-flex gap-2">
          <button type="submit" class="btn btn-primary" :disabled="loading">
            <span v-if="loading" class="spinner-border spinner-border-sm me-1"></span>
            {{ isEdit ? 'Сохранить' : 'Создать' }}
          </button>
          <router-link to="/places" class="btn btn-outline-secondary">Отмена</router-link>
        </div>

        <div v-if="isEdit" class="mt-4">
          <label class="form-label">Фото заведения</label>
          <div v-if="currentImageUrl" class="mb-2">
            <img :src="currentImageUrl" alt="Фото" class="img-thumbnail" style="max-height: 200px" />
          </div>
          <input type="file" class="form-control" accept="image/jpeg,image/png,image/webp" @change="handleImageSelect" />
          <button
            v-if="selectedImage"
            type="button"
            class="btn btn-outline-success btn-sm mt-2"
            :disabled="uploadingImage"
            @click="handleImageUpload"
          >
            <span v-if="uploadingImage" class="spinner-border spinner-border-sm me-1"></span>
            Загрузить фото
          </button>
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

const route = useRoute()
const router = useRouter()
const placesStore = usePlacesStore()
const catalogs = useCatalogsStore()

const isEdit = computed(() => !!route.params.id)
const error = ref('')
const loading = ref(false)
const selectedImage = ref(null)
const uploadingImage = ref(false)
const currentImageUrl = ref('')

const form = ref({
  name: '',
  city: '',
  address: '',
  cuisine_type_id: null,
  website: '',
  category_ids: [],
  lat: null,
  lng: null
})

onMounted(async () => {
  await catalogs.fetchAll()
  if (isEdit.value) {
    const place = await placesStore.fetchPlace(route.params.id)
    form.value = {
      name: place.name,
      city: place.city || '',
      address: place.address || '',
      cuisine_type_id: place.cuisine_type_id || null,
      website: place.website || '',
      category_ids: [],
      lat: place.lat || null,
      lng: place.lng || null
    }
    currentImageUrl.value = place.image_url || ''
  }
})

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    const data = {
      name: form.value.name,
      address: form.value.address || undefined,
      city: form.value.city || undefined,
      cuisine_type_id: form.value.cuisine_type_id || undefined,
      website: form.value.website || undefined,
      category_ids: form.value.category_ids,
      lat: form.value.lat || undefined,
      lng: form.value.lng || undefined
    }

    if (isEdit.value) {
      await placesStore.updatePlace(route.params.id, data)
      router.push(`/places/${route.params.id}`)
    } else {
      const created = await placesStore.createPlace(data)
      router.push(`/places/${created.id}`)
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка при сохранении'
  } finally {
    loading.value = false
  }
}

function handleImageSelect(e) {
  selectedImage.value = e.target.files[0] || null
}

async function handleImageUpload() {
  if (!selectedImage.value) return
  uploadingImage.value = true
  try {
    const result = await placesStore.uploadImage(route.params.id, selectedImage.value)
    currentImageUrl.value = result.image_url
    selectedImage.value = null
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка загрузки фото'
  } finally {
    uploadingImage.value = false
  }
}
</script>
