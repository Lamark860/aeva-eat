<template>
  <!-- Mobile-only bottom tab bar. Hidden on md+ where the top navbar serves the same role. -->
  <nav v-if="auth.isAuthenticated" class="bottom-tab-bar d-md-none">
    <router-link
      v-for="t in tabs"
      :key="t.to"
      :to="t.to"
      class="bottom-tab"
      :class="{ active: isActive(t) }"
    >
      <span class="bottom-tab-icon" aria-hidden="true">{{ t.icon }}</span>
      <span class="bottom-tab-label">{{ t.label }}</span>
    </router-link>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const route = useRoute()

const tabs = [
  { to: '/',         label: 'Главная',   icon: '🏠', match: (p) => p === '/' },
  { to: '/places',   label: 'Заведения', icon: '🍽', match: (p) => p === '/places' || p.startsWith('/places/') },
  { to: '/map',      label: 'Карта',     icon: '🗺', match: (p) => p === '/map' },
  { to: '/profile',  label: 'Профиль',   icon: '👤', match: (p) => p === '/profile' || p.startsWith('/invites') },
]

// router-link's auto active class only matches exact paths; we want
// /places/142 to highlight the Заведения tab, etc.
const currentPath = computed(() => route.path)
function isActive(t) { return t.match(currentPath.value) }
</script>

<style scoped lang="scss">
.bottom-tab-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1030;
  display: flex;
  justify-content: space-around;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(12px);
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.04);
  padding-bottom: var(--aeva-safe-bottom, 0px);
}

.bottom-tab {
  flex: 1 1 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 8px 4px 6px;
  text-decoration: none;
  color: #6c757d;
  font-size: 11px;
  font-weight: 500;
  min-height: 56px;
  transition: color 0.15s ease;

  &.active,
  &.router-link-exact-active {
    color: var(--bs-primary);
  }
}

.bottom-tab-icon {
  font-size: 22px;
  line-height: 1;
  margin-bottom: 2px;
}

.bottom-tab-label {
  letter-spacing: -0.01em;
}
</style>
