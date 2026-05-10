import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    // Allow LAN access by hostname / IP. Vite 5+ blocks unknown Host headers
    // by default as anti-DNS-rebinding; whitelist the names we use locally.
    allowedHosts: [
      'localhost',
      '.local',                  // any *.local (Bonjour) — MacBook-Pro-Maxim.local, aeva.local, etc.
      '192.168.1.10',
      '.lan',                    // some routers expose hostnames as *.lan
      '.ts.net',                 // Tailscale MagicDNS, if used later
      '.ngrok-free.dev',         // ngrok-туннель для шаринга снаружи (free tier)
      '.ngrok.app',              // ngrok платный домен
      '.trycloudflare.com',      // Cloudflare Tunnel — на случай переключения
    ],
    // HMR через ngrok: WS-канал по разным причинам бьётся (ws-upgrade
    // через CDN-туннель, browser-interstitial, free-tier лимиты), и
    // Vite после нескольких неудачных пингов делает full-reload страницы —
    // отсюда «мерцание / постоянные перезагрузки» в туннеле.
    // Отключаем HMR на время демо-сессии. Локальная разработка через
    // localhost тоже без HMR — обновляем вручную Cmd+R, это терпимо.
    // Когда туннель закроем, можно вернуть `hmr: true`.
    hmr: false,
    watch: {
      usePolling: true
    }
  }
})
