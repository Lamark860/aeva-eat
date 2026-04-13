import { defineStore } from 'pinia'
import { ref } from 'vue'
import http from '../api/http'

export const useWishlistStore = defineStore('wishlist', () => {
  const wishlistIds = ref(new Set())
  const wishlistPlaces = ref([])
  const customItems = ref([])
  const loading = ref(false)

  async function fetchIds() {
    try {
      const { data } = await http.get('/wishlist/ids')
      wishlistIds.value = new Set(data || [])
    } catch { /* ignore if not authed */ }
  }

  async function fetchPlaces() {
    loading.value = true
    try {
      const { data } = await http.get('/wishlist')
      wishlistPlaces.value = data || []
    } finally {
      loading.value = false
    }
  }

  async function fetchCustom() {
    try {
      const { data } = await http.get('/wishlist/custom')
      customItems.value = data || []
    } catch { /* ignore */ }
  }

  function isWishlisted(placeId) {
    return wishlistIds.value.has(placeId)
  }

  async function toggle(placeId) {
    if (isWishlisted(placeId)) {
      await http.delete(`/wishlist/${placeId}`)
      wishlistIds.value.delete(placeId)
      wishlistPlaces.value = wishlistPlaces.value.filter(p => p.id !== placeId)
    } else {
      await http.post(`/wishlist/${placeId}`)
      wishlistIds.value.add(placeId)
    }
  }

  async function addCustom(name, note) {
    const { data } = await http.post('/wishlist/custom', { name, note: note || undefined })
    customItems.value.unshift(data)
    return data
  }

  async function deleteCustom(id) {
    await http.delete(`/wishlist/custom/${id}`)
    customItems.value = customItems.value.filter(i => i.id !== id)
  }

  return { wishlistIds, wishlistPlaces, customItems, loading, fetchIds, fetchPlaces, fetchCustom, isWishlisted, toggle, addCustom, deleteCustom }
})
