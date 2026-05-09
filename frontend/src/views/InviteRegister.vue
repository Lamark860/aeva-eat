<template>
  <div class="row justify-content-center">
    <div class="col-12 col-sm-10 col-md-5">
      <div class="card shadow-sm">
        <div class="card-body">
          <h3 class="card-title text-center mb-4">Регистрация по приглашению</h3>

          <div v-if="checking" class="text-center py-4">
            <div class="spinner-border text-primary" role="status"></div>
            <p class="mt-2 text-muted">Проверяем приглашение…</p>
          </div>

          <div v-else-if="invalidCode" class="text-center py-4">
            <div class="fs-1 mb-2">😔</div>
            <p class="text-danger fw-bold">{{ invalidMessage }}</p>
            <router-link to="/login" class="btn btn-outline-primary mt-2">Войти</router-link>
          </div>

          <template v-else>
            <div class="alert alert-success d-flex align-items-center mb-3">
              <span class="me-2">🎉</span>
              <span>Вас пригласил(а) <strong>{{ creatorName }}</strong></span>
            </div>

            <div v-if="error" class="alert alert-danger">{{ error }}</div>

            <form @submit.prevent="handleRegister">
              <div class="mb-3">
                <label for="username" class="form-label">Логин *</label>
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
                <label for="displayName" class="form-label">Имя <span class="text-muted fw-normal">(необязательно)</span></label>
                <input
                  id="displayName"
                  v-model="displayName"
                  type="text"
                  class="form-control"
                  placeholder="Как вас называть"
                />
              </div>

              <div class="mb-3">
                <label for="password" class="form-label">Пароль *</label>
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
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import http from '../api/http'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const inviteCode = route.params.code

const checking = ref(true)
const invalidCode = ref(false)
const invalidMessage = ref('')
const creatorName = ref('')

const username = ref('')
const displayName = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

onMounted(async () => {
  try {
    const { data } = await http.get(`/invites/validate/${inviteCode}`)
    if (data.valid) {
      creatorName.value = data.creator_name
    } else {
      invalidCode.value = true
      invalidMessage.value = 'Недействительный инвайт-код'
    }
  } catch (e) {
    invalidCode.value = true
    invalidMessage.value = e.response?.data?.error || 'Недействительный инвайт-код'
  } finally {
    checking.value = false
  }
})

async function handleRegister() {
  error.value = ''
  loading.value = true
  try {
    await auth.register(username.value, displayName.value, password.value, inviteCode)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}
</script>
