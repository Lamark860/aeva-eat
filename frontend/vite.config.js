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
    ],
    watch: {
      usePolling: true
    }
  }
})
