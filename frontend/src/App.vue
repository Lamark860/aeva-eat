<template>
  <div>
    <nav class="navbar navbar-expand-md navbar-dark bg-dark">
      <div class="container">
        <router-link class="navbar-brand" to="/">🍽 AEVA Eat</router-link>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navMain">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navMain">
          <div class="navbar-nav me-auto">
            <router-link class="nav-link" to="/places">Заведения</router-link>
            <router-link class="nav-link" to="/map">Карта</router-link>
          </div>
          <div class="navbar-nav ms-auto">
            <template v-if="auth.isAuthenticated">
              <router-link class="nav-link" to="/profile">{{ auth.user?.username }}</router-link>
              <button class="btn btn-outline-light btn-sm ms-2" @click="logout">Выйти</button>
            </template>
            <template v-else>
              <router-link class="nav-link" to="/login">Войти</router-link>
              <router-link class="nav-link" to="/register">Регистрация</router-link>
            </template>
          </div>
        </div>
      </div>
    </nav>

    <main class="container-fluid container-lg py-4">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from './stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

auth.init()

function logout() {
  auth.logout()
  router.push('/login')
}
</script>
