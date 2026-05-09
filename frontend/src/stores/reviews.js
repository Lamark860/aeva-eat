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

  // Multi-фото за один multipart-запрос. Бэк добавляет в стопку review.photos
  // в порядке прихода. Лимит — 5 фото на отзыв; если превышен, бэк отвергает
  // целиком (HTTP 400).
  async function uploadReviewPhotos(placeId, reviewId, files) {
    if (!files || files.length === 0) return null
    const formData = new FormData()
    for (const f of files) formData.append('photos', f)
    const { data } = await http.post(`/places/${placeId}/reviews/${reviewId}/photos`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    const idx = reviews.value.findIndex(r => r.id === reviewId)
    if (idx !== -1) {
      const existing = reviews.value[idx].photos || []
      reviews.value[idx].photos = [...existing, ...(data.photos || [])]
    }
    return data
  }

  async function deleteReviewPhoto(placeId, reviewId, photoId) {
    await http.delete(`/places/${placeId}/reviews/${reviewId}/photos/${photoId}`)
    const idx = reviews.value.findIndex(r => r.id === reviewId)
    if (idx !== -1) {
      reviews.value[idx].photos = (reviews.value[idx].photos || []).filter(p => p.id !== photoId)
    }
  }

  async function uploadReviewVideo(placeId, reviewId, blob) {
    const formData = new FormData()
    // Pick filename extension based on the actual blob type so the backend's
    // Content-Type validation (video/webm | video/mp4) accepts it.
    const ext = (blob.type || '').includes('mp4') ? 'mp4' : 'webm'
    formData.append('video', blob, `video.${ext}`)
    const { data } = await http.post(`/places/${placeId}/reviews/${reviewId}/video`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    const idx = reviews.value.findIndex(r => r.id === reviewId)
    if (idx !== -1) reviews.value[idx].video_url = data.video_url
    return data
  }

  return {
    reviews,
    loading,
    fetchByPlace,
    fetchByUser,
    createReview,
    updateReview,
    deleteReview,
    uploadReviewImage,
    uploadReviewPhotos,
    deleteReviewPhoto,
    uploadReviewVideo,
  }
})
