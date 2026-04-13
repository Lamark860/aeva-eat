<template>
  <div>
    <nav class="navbar navbar-expand-md navbar-light">
      <div class="container">
        <router-link class="navbar-brand d-flex align-items-center gap-2" to="/">
          <span class="brand-icon">🍽</span> AEVA Eat
        </router-link>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navMain">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navMain">
          <div class="navbar-nav me-auto">
            <router-link class="nav-link" to="/places">Заведения</router-link>
            <router-link class="nav-link" to="/map">Карта</router-link>
          </div>
          <div class="navbar-nav ms-auto align-items-center">
            <template v-if="auth.isAuthenticated">
              <router-link class="nav-link d-flex align-items-center gap-1" to="/profile">
                <span class="avatar-sm">{{ auth.user?.username?.charAt(0)?.toUpperCase() }}</span>
                {{ auth.user?.username }}
              </router-link>
              <button class="btn btn-outline-secondary btn-sm ms-2" @click="logout">Выйти</button>
            </template>
            <template v-else>
              <router-link class="nav-link" to="/login">Войти</router-link>
              <router-link class="btn btn-primary btn-sm ms-2" to="/register">Регистрация</router-link>
            </template>
          </div>
        </div>
      </div>
    </nav>

    <main class="container-fluid container-lg py-4">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <ToastContainer />

    <footer class="app-footer text-center text-muted py-3 mt-4">
      <small>
        Картографические данные предоставлены
        <a href="https://yandex.ru/legal/maps_termsofuse/" target="_blank" rel="noopener">Яндекс Картами</a>
      </small>
    </footer>
  </div>
</template>

<script setup>
import { useAuthStore } from './stores/auth'
import { useWishlistStore } from './stores/wishlist'
import { useRouter } from 'vue-router'
import ToastContainer from './components/ToastContainer.vue'

const auth = useAuthStore()
const wishlist = useWishlistStore()
const router = useRouter()

auth.init()
if (auth.isAuthenticated) {
  wishlist.fetchIds()
}

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.avatar-sm {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--bs-primary);
  color: #fff;
  font-weight: 700;
  font-size: 0.8rem;
}

.brand-icon {
  font-size: 1.4rem;
}
</style>
