import { defineStore } from 'pinia'
import { ref } from 'vue'
import http from '../api/http'

// Feed-индикатор «новостей» (NEXT.md §C2). Хранит unread-count, дёргает
// /api/feed/unread-count периодически и на focus вкладки. Сбрасывается
// после открытия Доски через markSeen → POST /api/feed/seen.
export const useFeedStore = defineStore('feed', () => {
  const unread = ref(0)
  let timer = null
  let started = false

  function hasToken() {
    return typeof localStorage !== 'undefined' && !!localStorage.getItem('token')
  }

  async function refresh() {
    // Не дёргаем фоновый poll, если нет токена — иначе генерим лишние
    // 401 и (хоть теперь без redirect) зря шумим в логах.
    if (!hasToken()) return
    try {
      const { data } = await http.get('/feed/unread-count')
      unread.value = data?.count ?? 0
    } catch (e) {
      // 401 → токен протух, но redirect делать не нужно (это фон).
      // Просто остановим polling: при следующем валидном входе start()
      // снова запустится через watch isAuthenticated в BottomTabBar.
      if (e.response?.status === 401) stop()
    }
  }

  async function markSeen() {
    if (!hasToken()) return
    unread.value = 0
    try { await http.post('/feed/seen') } catch { /* ignore */ }
  }

  function onFocus() { refresh() }

  function start() {
    if (started || typeof window === 'undefined') return
    if (!hasToken()) return
    started = true
    refresh()
    // Каждые 60 секунд плюс при возврате фокуса.
    timer = window.setInterval(refresh, 60_000)
    window.addEventListener('focus', onFocus)
  }

  function stop() {
    if (!started) return
    started = false
    if (timer) { window.clearInterval(timer); timer = null }
    if (typeof window !== 'undefined') {
      window.removeEventListener('focus', onFocus)
    }
  }

  return { unread, refresh, markSeen, start, stop }
})
