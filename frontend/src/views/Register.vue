<template>
  <div class="row justify-content-center">
    <div class="col-md-5">
      <div class="card shadow-sm">
        <div class="card-body">
          <h3 class="card-title text-center mb-4">Регистрация</h3>

          <div v-if="error" class="alert alert-danger">{{ error }}</div>

          <form @submit.prevent="handleRegister">
            <div class="mb-3">
              <label for="username" class="form-label">Имя пользователя</label>
              <input
                id="username"
                v-model="username"
                type="text"
                class="form-control"
                placeholder="username"
                required
              />
            </div>

            <div class="mb-3">
              <label for="email" class="form-label">Email</label>
              <input
                id="email"
                v-model="email"
                type="email"
                class="form-control"
                placeholder="your@email.com"
                required
              />
            </div>

            <div class="mb-3">
              <label for="password" class="form-label">Пароль</label>
              <input
                id="password"
                v-model="password"
                type="password"
                class="form-control"
                placeholder="Минимум 6 символов"
                minlength="6"
                required
              />
            </div>

            <button type="submit" class="btn btn-primary w-100" :disabled="loading">
              <span v-if="loading" class="spinner-border spinner-border-sm me-1"></span>
              Создать аккаунт
            </button>
          </form>

          <p class="text-center mt-3 mb-0">
            Уже есть аккаунт?
            <router-link to="/login">Войти</router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleRegister() {
  error.value = ''
  loading.value = true
  try {
    await auth.register(username.value, email.value, password.value)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}
</script>
