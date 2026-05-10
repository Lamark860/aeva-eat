# AEVA Eat — заметки сессии (для продолжения после /clear)

Дата: 2026-05-10. Документ временный — удалить когда работа подхвачена в новой сессии.

---

## Где мы остановились

Дизайнер ответил на все 10 вопросов из старой версии `OPEN-QUESTIONS.md`.
Ответы **зафиксированы verbatim** в `DESIGN-DECISIONS.md` → раздел «Раунд 2 решений (после backend-итерации)».

`OPEN-QUESTIONS.md` обновлён — там теперь компактная сводка «что закрыто» со ссылками.

**Имплементация ответов ещё не начата.** Я только начал — успел сделать пункт #6 (удалить «wishlist» из AddArtifactSheet) задержался на отладке ngrok-тестирования (см. ниже). Продолжать с пункта #6 (или сразу с #9, если #6 уже сделан — проверь diff).

---

## Что дальше — 9 пунктов имплементации

Порядок от дешёвых к дорогим. Все полные формулировки и нюансы — в `DESIGN-DECISIONS.md` § «Раунд 2 решений».

| # | Что | Где описано | Прикидка |
|---|---|---|---|
| **#6** | Удалить пункт «wishlist» из `AddArtifactSheet.vue` | Q6 | 5 минут |
| **#9** | Dark-theme audit — все цвета через CSS-переменные, явно игнорить `prefers-color-scheme: dark` | Q9 | 30 минут |
| **#10** | Анимации в ленте: fade+slide-up артефакта 280мс, height+stagger разворота недели 320мс. Hover-приподнятие НЕ делать. | Q10 | 1 час |
| **#8** | «любит грузинскую — 11 раз» — SQL-агрегат favorite_cuisine на /api/users/:id + рендер в Profile/PersonPage. Строго рукописно с глаголом. | Q8 | 30 минут |
| **#2** | Расширенный /api/places/:id: `gem_status.marked_by/first_marked_at` + `attendance.visit_count` + `ratings_per_user`. UI: рукописная подпись «отметила Аня · 12 марта (+ Серёжа, Миша)» под штампом жемчужины, ряд `аватарка · ×3` (рукописно!) под общей оценкой. ratings_per_user — **только для бэка** (sort+future), отдельную таблицу не рисуем. | Q2 | 2 часа |
| **#1** | Wishlist-артефакт на Доске — отдельный класс. Активный = `--paper` + терра-штамп `план` + канцелярская кнопка. Зачёркнутый = `--paper-deep` + **SVG-волнистый штрих** (не CSS line-through!) + moss-штамп `сходили ✓` + мини-полароид внахлёст. Бэк уже готов (`/api/wishlist/all` + триггер struck). | Q1 | 2-3 часа |
| **#5** | UI расширенных фильтров (бэк готов): drawer-секции «кто был» (аватарки multi-select, поиск если >10), «когда» (date-input «с/по» + preset-чипы `этот год` `прошлый год` `последние 30 дней`), sort-pill `по оценке N` с дропдауном друга. Бэк параметры: `attended_by`, `visit_from`, `visit_to`, `sort=rating_user:N`. | Q5 | 2-3 часа |
| **#4** | useFeed → `/api/feed` (event-driven), группировка `(place_id, date(occurred_at), attendees)`. Это закрывает L1 edge case «мы туда ходим раз в месяц». Сейчас группируется только по place_id через `/api/places?sort=new`. | Q4 | 2-3 часа фронт, бэк есть |
| **#3** | `/p/:id` публичный шаринг — Go-handler без auth с OG-метатегами. Cover full-bleed → бумажная плашка → имя серифой, город каракулями, штамп `жемчужина` если is_gem, БЕЗ рейтингов и БЕЗ имён авторов. CTA `войти, чтобы увидеть наши впечатления`. og:title=`Place · City`, og:description=`камерный дневник еды`, og:image=cover. | Q3 | 3-4 часа |

Итого по плану ~16-20 часов фокус-работы.

---

## ngrok-туннель

**Статус:** запущен в фоне на ноуте, URL: `https://enunciate-deity-grill.ngrok-free.dev`

**Управление:**
```bash
# Узнать текущий URL (если процесс жив):
curl -s http://localhost:4040/api/tunnels | python3 -c 'import sys,json; print(json.load(sys.stdin)["tunnels"][0]["public_url"])'

# Остановить:
pkill -f 'ngrok http 8091'

# Перезапустить (authtoken уже сохранён в ngrok config):
ngrok http 8091 --log=stdout --log-format=json &
```

**Запоминай:** при перезапуске URL **изменится** (бесплатный план). Чтобы постоянный домен — Cloudflare Tunnel или платный ngrok.

**ngrok inspector:** http://localhost:4040 — видно все запросы в реальном времени.

---

## Незакоммиченные правки (7 файлов)

Это **общие баг-фиксы + ngrok-специфика**, обнаруженные при тестировании. Должны попасть в коммит, но порционно.

### Группа A — общие баг-фиксы (КОММИТИТЬ ОБЯЗАТЕЛЬНО)

Полезны независимо от ngrok. Чинят реальные баги:

- **`frontend/src/api/http.js`** — http-interceptor больше не делает hard-reload на 401 для silent-poll эндпоинтов (`/feed/seen`, `/feed/unread-count`) и не зацикливает редирект на /login если уже там. Это разрывает «401 → reload → 401 → reload» цикл при истёкшем токене.
- **`frontend/src/stores/feed.js`** — `useFeedStore` проверяет наличие токена перед вызовами `/feed/*`. Не дёргает API на /login и без auth.
- **`frontend/public/sw.js`** — kill-switch SW. Старый PWA-SW кэшировал app shell, что давало «мерцание старого дизайна» при любых reload. Новый SW тихо разрегистрирует себя и чистит кэши. Полезно для всех будущих юзеров с накопленным SW-кэшем.

### Группа B — ngrok-специфика (КОММИТИТЬ ИЛИ ОТКАТИТЬ — на твой выбор)

Полезны пока туннель открыт; для локального dev можно вернуть как было:

- **`frontend/vite.config.js`** — добавлены `.ngrok-free.dev`/`.ngrok.app`/`.trycloudflare.com` в allowedHosts + **`hmr: false`**. Без `hmr: false` Vite периодически делает full-reload через ngrok туннель (WS-флапает). Когда туннель закроем — вернуть `hmr: { /* default */ }` или удалить строку.
- **`nginx/nginx.conf`** — `map $http_upgrade $connection_upgrade` + WS upgrade на корне (`location /`). Раньше upgrade был только для `/ws`. Это **общая** правка — пригодится для любого WS через nginx, не только ngrok. Можно оставить навсегда.
- **`spec/DESIGN-DECISIONS.md`** — добавлен раздел «Раунд 2 решений» с verbatim-ответами дизайнера. ВАЖНО НЕ ПОТЕРЯТЬ.
- **`spec/OPEN-QUESTIONS.md`** — переписан как сводка «всё закрыто, см. DESIGN-DECISIONS». ВАЖНО НЕ ПОТЕРЯТЬ.

### Рекомендация по коммитам

```bash
# Коммит 1: общие баг-фиксы — отдельный коммит для чёткой истории
git add frontend/src/api/http.js frontend/src/stores/feed.js frontend/public/sw.js nginx/nginx.conf
git commit -m "fix: 401 redirect loop + SW kill-switch + nginx WS upgrade

- Stop hard-reload on 401 for background polls (/feed/seen,
  /feed/unread-count). Loop occurred when stale token in
  localStorage triggered redirect-to-login while already on /login.
- Guard feed-store polls behind hasToken() check.
- New SW = kill-switch: deletes caches, claims clients, navigates,
  unregisters. Wipes stale prod-SW from prior phase 5 build.
- nginx: WS upgrade on / (not just /ws) so proxied dev servers
  with HMR keep working through tunnels and reverse proxies."

# Коммит 2: спек-обновления (важнее — там ответы дизайнера)
git add spec/
git commit -m "spec: round 2 designer decisions — close all 10 OPEN-QUESTIONS"

# Коммит 3: ngrok-специфика (опционально — можно перед закрытием туннеля откатить)
git add frontend/vite.config.js
git commit -m "frontend: ngrok/cloudflare allowedHosts + hmr off for tunnel demos

Vite HMR via public tunnel flaps and triggers full-reload churn.
Disable for now; flip back when developing locally."
```

---

## Состояние task-листа на момент /clear

Активный батч (#30-#39) — почти все pending, кроме #29 (записать решения) и #30 (#6 — удалить wishlist пункт; начат, но я не уверен закоммитил ли изменение в `AddArtifactSheet.vue` — ПРОВЕРЬ).

Старые таски (#1-#28) — все завершены либо deleted. Можно их **очистить через TaskUpdate→deleted** в начале новой сессии, чтобы список не захламлялся.

---

## Как зайти в новую сессию

1. **Прочитай этот документ + `DESIGN-DECISIONS.md` § «Раунд 2 решений»** — это весь контекст.
2. **Чек незакоммиченного:** `git diff` и `git status` чтобы увидеть что висит.
3. **Закоммить группу A** (общие баг-фиксы) сразу — они не должны висеть.
4. **Туннель:** проверь `pgrep -fl ngrok` — жив ли. Если нет, `ngrok http 8091 &` — токен сохранён.
5. **Стартовый шаг:** проверить `frontend/src/components/scrapbook/AddArtifactSheet.vue` — удалён ли блок `wishlist` (пункт #6). Если нет — удалить, закоммитить, идти на #9.
6. **Дальше по таблице** «9 пунктов имплементации» (см. выше).

---

## Контекст по проекту

- Документация: `spec/README.md` → точка входа. `STATUS.md` (журнал что сделано), `DESIGN-DECISIONS.md` (решения), `NEXT.md` (изначальные продуктовые предложения), `backend.md`/`mvp-scope.md`/`product.md`/`design.md` (исходное ТЗ).
- Стек: Go-бэкенд (chi router) + Vue 3 / Pinia / vite-dev фронт. Postgres. Nginx прокси на 8091, фронт dev на 5173, бэк на 8086 (через nginx — `/api/*`).
- Запуск всего: `docker compose up -d` из корня репо.
- Билд бэка после изменений: `docker compose build backend && docker compose up -d backend`.
- Сборка фронта: `cd frontend && npx vite build` (если нужно), но в dev — Vite ловит изменения автоматически (без HMR — нужен ручной refresh).
- Тесты бэка: `cd backend && go test ./...` (всё зелёное на момент паузы).
- Lint фронта: `cd frontend && npm run lint` — есть **3 pre-existing warnings в `LocationPicker.vue`** (зафиксированы как low-priority, не трогать).

---

## Запомнить про user-предпочтения

- Хочет коммиты осмысленные, в стиле прошлых: короткий заголовок + body
- Любит когда сначала закоммичены чистые баг-фиксы, потом фичи
- При работе с дизайн-системой — строго следовать `DESIGN-DECISIONS.md`. Если что-то не описано — задавать вопрос явно (через `OPEN-QUESTIONS.md`), а не угадывать
- В UI — никаких эмодзи кроме 🎲 (явно одобренного дизайнером)
- Строго `Caveat` для рукописных подписей, `Lora` для серифа. Цифры в стиле «×3» (рукописные), не «3» (плоско)
- Углы наклона полароидов: 1–5°, не больше
- Готов к большим батчам работы — «делай всё что понятно, потом обсудим»

---

После того как этот документ перенесён в новую сессию и работа подхвачена — **удалить файл**.
