<template>
  <div>
    <div class="text-center" v-if="auth.isAuthenticated">
      <h1 class="display-5 mb-3">🍽 Привет, {{ auth.user?.username }}!</h1>
      <p class="lead text-muted">
        Добро пожаловать в AEVA Eat — твой личный дневник впечатлений от заведений.
      </p>
      <p class="text-muted">Скоро здесь появятся заведения, отзывы и карта.</p>
    </div>
    <div class="text-center" v-else>
      <h1 class="display-5 mb-3">🍽 AEVA Eat</h1>
      <p class="lead text-muted">
        Храни впечатления о заведениях вместе с друзьями.
      </p>
      <div class="d-flex gap-2 justify-content-center mt-4">
        <router-link to="/login" class="btn btn-primary">Войти</router-link>
        <router-link to="/register" class="btn btn-outline-primary">Регистрация</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
</script>
