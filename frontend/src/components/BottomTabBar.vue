<template>
  <!-- Bottom tab bar — scrapbook chrome, mobile-first. Hidden on md+ desktops. -->
  <nav v-if="auth.isAuthenticated" class="sb-tabbar bottom-fixed d-md-none">
    <router-link
      v-for="t in tabs"
      :key="t.to"
      :to="t.to"
      class="tab"
      :class="{ active: isActive(t) }"
    >
      <span class="glyph" aria-hidden="true">{{ t.glyph }}</span>
      <span class="lbl">{{ t.label }}</span>
      <!-- C2 — точка-индикатор «есть новое в ленте» на табе Доска. -->
      <span v-if="t.to === '/' && feed.unread > 0" class="dot" aria-hidden="true"></span>
    </router-link>
  </nav>
</template>

<script setup>
import { computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useFeedStore } from '../stores/feed'

const auth = useAuthStore()
const feed = useFeedStore()
const route = useRoute()

// Glyphs lifted from spec/example/v2 — Caveat-rendered to match scrapbook hand.
const tabs = [
  { to: '/',         label: 'Доска',     glyph: '✦', match: (p) => p === '/' },
  { to: '/places',   label: 'Найти',     glyph: '⌕', match: (p) => p === '/places' || p.startsWith('/places/') },
  { to: '/map',      label: 'Карта',     glyph: '◎', match: (p) => p === '/map' },
  { to: '/profile',  label: 'Я',         glyph: '✺', match: (p) => p === '/profile' || p.startsWith('/invites') },
]

const currentPath = computed(() => route.path)
function isActive(t) { return t.match(currentPath.value) }

// Запускаем polling индикатора, когда юзер аутентифицирован.
watch(() => auth.isAuthenticated, (yes) => {
  if (yes) feed.start()
  else feed.stop()
}, { immediate: true })

// При открытии Доски сбрасываем точку (markSeen на сервере).
watch(currentPath, (p) => {
  if (p === '/') feed.markSeen()
}, { immediate: true })

onMounted(() => { if (auth.isAuthenticated) feed.start() })
onUnmounted(() => feed.stop())
</script>

<style scoped lang="scss">
.bottom-fixed {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1030;
}

// Override the .lbl line height for label text inside the scrapbook tabbar.
.sb-tabbar .tab .lbl {
  font-family: var(--sb-serif);
  font-size: 10.5px;
  letter-spacing: 0.04em;
  line-height: 1.1;
  text-transform: lowercase;
}

/* C2 — небольшая точка-индикатор «новости» рядом с глифом таба. */
.sb-tabbar .tab {
  position: relative;
}
.sb-tabbar .tab .dot {
  position: absolute;
  top: 6px;
  right: calc(50% - 16px);
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--sb-terracotta);
  box-shadow: 0 0 0 2px var(--sb-paper-card);
}
</style>
