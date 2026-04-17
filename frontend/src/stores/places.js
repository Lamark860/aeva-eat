import { defineStore } from 'pinia'
import { ref } from 'vue'
import http from '../api/http'

export const usePlacesStore = defineStore('places', () => {
  const places = ref([])
  const total = ref(0)
  const page = ref(1)
  const limit = ref(20)
  const currentPlace = ref(null)
  const loading = ref(false)
  const filters = ref({
    city: '',
    cuisine_type_ids: [],
    category_ids: [],
    min_rating: '',
    is_gem: false,
    search: '',
    sort: ''
  })

  function _buildParams() {
    const params = {}
    Object.entries(filters.value).forEach(([key, val]) => {
      if (Array.isArray(val)) {
        if (val.length > 0) params[key === 'cuisine_type_ids' ? 'cuisine_type_id' : key === 'category_ids' ? 'category_id' : key] = val.join(',')
      } else if (val !== '' && val !== false && val !== null) {
        params[key] = val
      }
    })
    return params
  }

  async function fetchPlaces() {
    loading.value = true
    try {
      const params = _buildParams()
      params.limit = limit.value
      params.page = page.value
      const { data } = await http.get('/places', { params })
      places.value = data.places || []
      total.value = data.total || 0
    } finally {
      loading.value = false
    }
  }

  async function fetchAllPlaces() {
    loading.value = true
    try {
      const params = _buildParams()
      params.limit = 0
      const { data } = await http.get('/places', { params })
      places.value = data.places || []
      total.value = data.total || 0
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

  return { places, total, page, limit, currentPlace, loading, filters, fetchPlaces, fetchAllPlaces, fetchPlace, createPlace, updatePlace, deletePlace, uploadImage }
})
