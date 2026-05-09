// AEVA Eat — minimal service worker.
// Strategy:
//   - Precache the app shell on install (index.html + manifest + icons).
//   - Runtime cache same-origin GET requests for static assets (cache-first).
//   - Pass through everything else to the network unchanged. API requests are
//     never cached because the data is mutable and we don't want to serve a
//     stale list of places offline.
//
// Bumping CACHE_VERSION evicts old caches on the next activate.

const CACHE_VERSION = 'aeva-eat-v1';
const SHELL_ASSETS = [
  '/',
  '/manifest.webmanifest',
  '/icon-192.png',
  '/icon-512.png',
  '/apple-touch-icon.png',
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_VERSION).then((cache) => cache.addAll(SHELL_ASSETS))
  );
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(keys.filter((k) => k !== CACHE_VERSION).map((k) => caches.delete(k)))
    )
  );
  self.clients.claim();
});

self.addEventListener('fetch', (event) => {
  const req = event.request;
  if (req.method !== 'GET') return;

  const url = new URL(req.url);
  if (url.origin !== self.location.origin) return;

  // Never cache the API — places/reviews/etc must be live.
  if (url.pathname.startsWith('/api/')) return;

  // SPA navigations: try network first so users get fresh HTML; fall back
  // to the cached shell when offline.
  if (req.mode === 'navigate') {
    event.respondWith(
      fetch(req).catch(() => caches.match('/'))
    );
    return;
  }

  // Static assets (Vite hashes asset filenames, so cache-first is safe — a
  // new build produces new URLs).
  event.respondWith(
    caches.match(req).then((cached) => {
      if (cached) return cached;
      return fetch(req).then((res) => {
        if (res.ok && res.type === 'basic') {
          const copy = res.clone();
          caches.open(CACHE_VERSION).then((cache) => cache.put(req, copy));
        }
        return res;
      });
    })
  );
});
