import { defineStore } from 'pinia'
import { ref } from 'vue'
import http from '../api/http'

export const useCatalogsStore = defineStore('catalogs', () => {
  const cuisineTypes = ref([])
  const categories = ref([])

  async function fetchCuisineTypes() {
    if (cuisineTypes.value.length) return
    const { data } = await http.get('/cuisine-types')
    cuisineTypes.value = data || []
  }

  async function fetchCategories() {
    if (categories.value.length) return
    const { data } = await http.get('/categories')
    categories.value = data || []
  }

  async function fetchAll() {
    await Promise.all([fetchCuisineTypes(), fetchCategories()])
  }

  return { cuisineTypes, categories, fetchCuisineTypes, fetchCategories, fetchAll }
})
