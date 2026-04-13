import { defineStore } from 'pinia'
import { ref } from 'vue'
import http from '../api/http'

export const useReviewsStore = defineStore('reviews', () => {
  const reviews = ref([])
  const loading = ref(false)

  async function fetchByPlace(placeId) {
    loading.value = true
    try {
      const { data } = await http.get(`/places/${placeId}/reviews`)
      reviews.value = data
    } finally {
      loading.value = false
    }
  }

  async function fetchByUser(userId) {
    loading.value = true
    try {
      const { data } = await http.get(`/users/${userId}/reviews`)
      reviews.value = data
    } finally {
      loading.value = false
    }
  }

  async function createReview(placeId, reviewData) {
    const { data } = await http.post(`/places/${placeId}/reviews`, reviewData)
    reviews.value.unshift(data)
    return data
  }

  async function updateReview(placeId, reviewId, reviewData) {
    const { data } = await http.put(`/places/${placeId}/reviews/${reviewId}`, reviewData)
    const idx = reviews.value.findIndex(r => r.id === reviewId)
    if (idx !== -1) reviews.value[idx] = data
    return data
  }

  async function deleteReview(placeId, reviewId) {
    await http.delete(`/places/${placeId}/reviews/${reviewId}`)
    reviews.value = reviews.value.filter(r => r.id !== reviewId)
  }

  async function uploadReviewImage(placeId, reviewId, file) {
    const formData = new FormData()
    formData.append('image', file)
    const { data } = await http.post(`/places/${placeId}/reviews/${reviewId}/image`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    const idx = reviews.value.findIndex(r => r.id === reviewId)
    if (idx !== -1) reviews.value[idx].image_url = data.image_url
    return data
  }

  return { reviews, loading, fetchByPlace, fetchByUser, createReview, updateReview, deleteReview, uploadReviewImage }
})
