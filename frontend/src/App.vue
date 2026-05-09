<template>
  <div :class="{ 'is-scrapbook': isScrapbook }">
    <nav v-if="!isScrapbook" class="navbar navbar-expand-md navbar-light">
      <div class="container">
        <router-link class="navbar-brand d-flex align-items-center gap-2" to="/">
          <span class="brand-icon">🍽</span> AEVA Eat
        </router-link>
        <!-- Hamburger toggler is suppressed on mobile; the bottom tab bar provides nav -->
        <button class="navbar-toggler d-none" type="button" data-bs-toggle="collapse" data-bs-target="#navMain">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navMain">
          <div class="navbar-nav me-auto">
            <template v-if="auth.isAuthenticated">
              <router-link class="nav-link" to="/places">Заведения</router-link>
              <router-link class="nav-link" to="/map">Карта</router-link>
              <router-link class="nav-link" to="/invites">Пригласить</router-link>
            </template>
          </div>
          <div class="navbar-nav ms-auto align-items-center">
            <template v-if="auth.isAuthenticated">
              <router-link class="nav-link d-flex align-items-center gap-1" to="/profile">
                <span class="avatar-sm">{{ auth.user?.username?.charAt(0)?.toUpperCase() }}</span>
                {{ auth.user?.username }}
              </router-link>
              <button class="btn btn-outline-secondary btn-sm ms-2" @click="logout">Выйти</button>
            </template>
          </div>
        </div>
        <!-- Mobile-only avatar shortcut to profile (since collapsed nav is hidden) -->
        <router-link
          v-if="auth.isAuthenticated"
          class="d-md-none ms-auto avatar-link"
          to="/profile"
          aria-label="Профиль"
        >
          <span class="avatar-sm">{{ auth.user?.username?.charAt(0)?.toUpperCase() }}</span>
        </router-link>
      </div>
    </nav>

    <main :class="isScrapbook ? 'sb-main' : 'container-fluid container-lg py-4 main-content'">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <ToastContainer />

    <BottomTabBar />

    <footer v-if="!isScrapbook" class="app-footer text-center text-muted py-3 mt-4 d-none d-md-block">
      <small>
        Картографические данные предоставлены
        <a href="https://yandex.ru/legal/maps_termsofuse/" target="_blank" rel="noopener">Яндекс Картами</a>
      </small>
    </footer>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from './stores/auth'
import { useWishlistStore } from './stores/wishlist'
import { useRoute, useRouter } from 'vue-router'
import ToastContainer from './components/ToastContainer.vue'
import BottomTabBar from './components/BottomTabBar.vue'

const auth = useAuthStore()
const wishlist = useWishlistStore()
const router = useRouter()
const route = useRoute()

// Scrapbook routes own their own chrome (wordmark / paper bg / scrapbook tabbar)
// — hide the global Bootstrap navbar/footer on those routes.
const isScrapbook = computed(() => !!route.meta.scrapbook)

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
.sb-main {
  /* Scrapbook screens render edge-to-edge — no Bootstrap container padding. */
  padding: 0;
  margin: 0;
}

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

.avatar-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  text-decoration: none;
}

.avatar-link .avatar-sm {
  width: 36px;
  height: 36px;
  font-size: 1rem;
}
</style>
