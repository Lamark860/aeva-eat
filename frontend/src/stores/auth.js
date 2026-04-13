import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import http from '../api/http'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(null)

  const isAuthenticated = computed(() => !!token.value)

  function init() {
    const savedToken = localStorage.getItem('token')
    if (savedToken) {
      token.value = savedToken
      fetchMe()
    }
  }

  async function register(username, email, password) {
    const { data } = await http.post('/auth/register', { username, email, password })
    token.value = data.token.access_token
    user.value = data.user
    localStorage.setItem('token', token.value)
  }

  async function login(email, password) {
    const { data } = await http.post('/auth/login', { email, password })
    token.value = data.token.access_token
    user.value = data.user
    localStorage.setItem('token', token.value)
  }

  async function fetchMe() {
    try {
      const { data } = await http.get('/auth/me')
      user.value = data
    } catch {
      logout()
    }
  }

  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
  }

  return { user, token, isAuthenticated, init, register, login, fetchMe, logout }
})
