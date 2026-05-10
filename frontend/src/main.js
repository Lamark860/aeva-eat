import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

import './assets/scss/main.scss'
import 'bootstrap/dist/js/bootstrap.bundle.min.js'

// Дожидаемся резолва текущего маршрута до mount'а — иначе на первый кадр
// `route.meta.scrapbook` ещё undefined, App.vue считает страницу
// «не-скрапбук» и мигает Bootstrap-навбар прежде чем переключиться.
const app = createApp(App)
app.use(createPinia())
app.use(router)

router.isReady().then(() => {
  app.mount('#app')
})

// Register the service worker only in production builds. The dev server doesn't
// need it (and SW caching makes Vite HMR confusing).
if ('serviceWorker' in navigator && import.meta.env.PROD) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js').catch((err) => {
      // eslint-disable-next-line no-console
      console.warn('SW registration failed:', err)
    })
  })
}
