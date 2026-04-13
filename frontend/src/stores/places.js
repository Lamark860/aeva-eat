import { defineStore } from 'pinia'
import { ref } from 'vue'
import http from '../api/http'

export const usePlacesStore = defineStore('places', () => {
  const places = ref([])
  const currentPlace = ref(null)
  const loading = ref(false)
  const filters = ref({
    city: '',
    cuisine_type_id: '',
    category_id: '',
    min_rating: '',
    is_gem: false,
    search: '',
    sort: ''
  })

  async function fetchPlaces() {
    loading.value = true
    try {
      const params = {}
      Object.entries(filters.value).forEach(([key, val]) => {
        if (val !== '' && val !== false && val !== null) {
          params[key] = val
        }
      })
      const { data } = await http.get('/places', { params })
      places.value = data
    } finally {
      loading.value = false
    }
  }

  async function fetchPlace(id) {
    loading.value = true
    try {
      const { data } = await http.get(`/places/${id}`)
      currentPlace.value = data
      return data
    } finally {
      loading.value = false
    }
  }

  async function createPlace(placeData) {
    const { data } = await http.post('/places', placeData)
    places.value.unshift(data)
    return data
  }

  async function updatePlace(id, placeData) {
    const { data } = await http.put(`/places/${id}`, placeData)
    const idx = places.value.findIndex(p => p.id === id)
    if (idx !== -1) places.value[idx] = data
    currentPlace.value = data
    return data
  }

  async function deletePlace(id) {
    await http.delete(`/places/${id}`)
    places.value = places.value.filter(p => p.id !== id)
  }

  async function uploadImage(id, file) {
    const formData = new FormData()
    formData.append('image', file)
    const { data } = await http.post(`/places/${id}/image`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    if (currentPlace.value && currentPlace.value.id === Number(id)) {
      currentPlace.value.image_url = data.image_url
    }
    return data
  }

  return { places, currentPlace, loading, filters, fetchPlaces, fetchPlace, createPlace, updatePlace, deletePlace, uploadImage }
})
