// AEVA Eat — kill-switch service worker (агрессивная версия).
//
// Контекст: предыдущая версия SW (CACHE_VERSION='aeva-eat-v1') кэшировала
// app shell и статику cache-first. На время dev-туннеля (ngrok) и пока
// идёт активная разработка это даёт «мерцание старого дизайна».
//
// Этот файл байт-отличается от старого, поэтому браузер автоматически
// установит его как новую версию SW. На активации:
//   1. skipWaiting + claim — забираем контроль над открытыми вкладками
//      (без claim новый SW не видит существующих клиентов)
//   2. Удаляем все кэши, которые SW успел накопить
//   3. unregister — снимаем регистрацию самого себя
//   4. Hard navigate всех вкладок — следующий load пойдёт через сеть

self.addEventListener('install', (event) => {
  // skipWaiting в install, чтобы новая версия не висела в waiting state
  // когда у клиента всё ещё открыта вкладка под старым SW.
  event.waitUntil(self.skipWaiting());
});

self.addEventListener('activate', (event) => {
  event.waitUntil((async () => {
    // 1. Сначала забрать контроль — иначе clients.matchAll() пуст
    //    (новый SW не контролирует никого по умолчанию).
    await self.clients.claim();

    // 2. Все кэши под нож.
    const keys = await caches.keys();
    await Promise.all(keys.map((k) => caches.delete(k)));

    // 3. Найти все вкладки в нашем scope, перезагрузить их.
    //    Делаем reload ДО unregister, чтобы клиенты ещё видели нас как
    //    активного SW и согласились на навигацию.
    const clients = await self.clients.matchAll({ type: 'window', includeUncontrolled: true });
    for (const client of clients) {
      try { await client.navigate(client.url); } catch { /* cross-origin / closed */ }
    }

    // 4. Снять регистрацию. После этого следующая навигация пойдёт уже
    //    без SW-прослойки.
    await self.registration.unregister();
  })());
});

// Никаких fetch-перехватов — пропускаем всё в сеть. Минимальное
// поведение, пока активный SW всё ещё в процессе самоуничтожения.
self.addEventListener('fetch', () => {
  // explicit no-op — браузер пойдёт в сеть напрямую
});
