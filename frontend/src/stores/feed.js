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

  async function refresh() {
    try {
      const { data } = await http.get('/feed/unread-count')
      unread.value = data?.count ?? 0
    } catch { /* ignore — индикатор не должен ломать UX */ }
  }

  async function markSeen() {
    unread.value = 0
    try { await http.post('/feed/seen') } catch { /* ignore */ }
  }

  function onFocus() { refresh() }

  function start() {
    if (started || typeof window === 'undefined') return
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
