# Тех-долг и продуктовые пробелы — AEVA Eat

> Реестр заведён 2026-06-05 по итогам сквозного аудита «MVP → продукт».
> Severity откалибрована состязательной перепроверкой по коду. Статусы:
> ✅ исправлено в этом цикле · 🔸 открыто · 💤 отложено осознанно.
>
> Связано: [`product.md`](./product.md) · [`backend.md`](./backend.md) ·
> [`../ROADMAP.md`](../ROADMAP.md) · [`../CHANGELOG.md`](../CHANGELOG.md)

---

## ✅ Исправлено (2026-06-05)

| Что | Где | Суть |
|---|---|---|
| Идентичность места = `name+address+city` | `migrations/016_place_identity`, `repository/place.go`, `handler/place.go` | Ключ `(name, city)` блокировал второй филиал сети («Волконский»). Теперь в ключ входит адрес → филиалы сосуществуют |
| 409 больше не тупик | `handler/place.go` Create, `stores/places.js`, `views/PlaceForm.vue` | На дубль возвращается тело существующего места; форма даёт развилку «перейти и оставить отзыв» / «уточнить адрес» |
| Типизированный детект конфликта | `repository/place.go` | `pq.Error` code `23505` вместо `strings.Contains(err, "idx_…")` — переименование индекса не ломает |
| Потеря видео/обложки при правке отзыва | `repository/review.go` Update | Update не трогает `image_url`/`video_url` (управляются upload-эндпоинтами) |
| IDOR на удалении инвайта | `handler/invite.go`, `repository/invite.go` | Удалить может только создатель/superuser (проверка в SQL, 403 иначе) |
| Fail-fast по `JWT_SECRET` | `config.go`, `cmd/api/main.go`, оба compose | В `APP_ENV=production` пустой/дефолтный секрет → отказ старта |
| Скрипт бэкапа | `backend/scripts/backup.sh`, `DEPLOY.md` | Дамп БД + архив `uploads_data` + ротация; осталось повесить cron на проде |

### ✅ Исправлено — батч 2 (2026-06-05)

| Что | Где | Суть |
|---|---|---|
| Superuser-bypass на места | `handler/place.go` (`canMutatePlace`) | superuser может править/удалять/менять фото любого места (общие данные круга) |
| Чистка файлов при удалении места | `handler/place.go` Delete, `repository/place.go` `CollectUploadPaths` | удаляются обложка места + фото/видео всех отзывов из `/uploads` (не сироты) |
| Предупреждение при удалении места | `views/PlaceDetail.vue` | confirm показывает число отзывов круга, которые удалятся |
| Guard загрузки `ymaps` + fallback | `index.html`, `LocationPicker.vue`, `MapView.vue`, `PlaceForm.vue` | при недоступном скрипте — сообщение + авто-раскрытие ручного ввода, без молчаливого пустого div |
| Таймауты HTTP-сервера | `cmd/api/main.go` | `ReadHeaderTimeout`/`IdleTimeout` (Slowloris), не ограничивая заливку видео |
| Таймаут+контекст suggest-прокси | `handler/suggest.go` | `http.Client{Timeout:5s}` + `NewRequestWithContext` |
| Видео по сигнатуре, не по заголовку | `handler/review.go` UploadVideo | `http.DetectContentType` (magic bytes), расширение из реального типа |
| Лимит длины при правке заметки | `handler/note.go` Update | 2000 симв., как в Create |
| Обработка ошибок delete/wishlist | `views/PlaceDetail.vue` | try/catch + toast.error на удалении места/отзыва и toggle wishlist |

### ✅ Исправлено — батч 3 (2026-06-05)

| Что | Где | Суть |
|---|---|---|
| Rate-limit | `middleware/ratelimit.go`, `main.go` | in-memory token-bucket: login/register 20/min по IP, suggest 60/min по userID; без внешних зависимостей |
| Лимит JSON-тела | `middleware.BodyLimit`, `main.go` | `MaxBytesReader` 1 MB на `application/json` (multipart-загрузки не трогает) |
| Поиск по name+city+cuisine | `repository/place.go` List | строка поиска ищет по названию, городу и кухне (как обещает плейсхолдер) |
| Сохранение `website` | `views/PlaceForm.vue` | сайт из Яндекс-подсказки и ручного поля теперь пишется (раньше терялся) |

### ✅ Исправлено — батч 4 (2026-06-05)

| Что | Где | Суть |
|---|---|---|
| Канонизация города | `handler/place.go` `normalizePlaceReq`, `repository/place.go` `CanonicalCity` | trim + приведение к существующему написанию → полка «По городам» не дробится на «Москва»/«москва» |
| Trim полей места | `handler/place.go` | name/address/city/website тримятся на сохранении (меньше дублей с хвостовыми пробелами) |
| Backend healthcheck в compose | `docker-compose.prod.yml` | `wget /api/health` + nginx ждёт `service_healthy` (раньше стартовал раньше готовности бэка → 502) |
| Удалён мёртвый `VideoKruzhok.vue` | `components/scrapbook/` | не использовался нигде |

### ✅ Исправлено — батч 5 (2026-06-05)

| Что | Где | Суть |
|---|---|---|
| Кэш suggest (TTL 60s) | `handler/suggest.go` | повторные запросы автокомплита не бьют по платной квоте Яндекса |
| review принадлежит place из URL | `handler/review.go` `reviewBelongsToPlace`, `repository/review.go` `PlaceIDOf` | Update/Delete проверяют, что `{rid}` принадлежит `{id}` (404 иначе) |
| Ротация docker-логов | `docker-compose.prod.yml` | `json-file max-size 10m × 3` всем сервисам — логи не забьют диск VPS |

---

## 🔸 Открытый тех-долг (по убыванию приоритета)

### Критичное / высокое

- **🔸 [ops] Бэкап на проде не на расписании.** Скрипт есть, но cron + off-site
  копию надо завести на VPS (`DEPLOY.md`). Пока не сделано — риск потери данных
  при отказе диска остаётся.
- **✅→🔸 [product] Admin-bypass на места — СДЕЛАНО (батч 2).** superuser теперь
  правит/удаляет/меняет фото любого места (`canMutatePlace`). Осталось 🔸: дать
  обычному участнику круга править общие места (сейчас только создатель/superuser)
  и переназначать `created_by` при NULL.
- **✅→🔸 [data-model] Каскадное удаление места — частично (батч 2).** Файлы
  `/uploads` теперь чистятся (`CollectUploadPaths`), confirm показывает число
  отзывов. Осталось 🔸: мягкое удаление/архив вместо физического DELETE, чтобы
  чужой контент нельзя было снести безвозвратно одним тапом.
- **🔸 [ops] Авто-деплой после auto-merge молча не срабатывает.** `automerge.yml`
  мержит дефолтным `GITHUB_TOKEN` (секрет `GH_PAT` не задан), а GitHub не
  триггерит workflow на push от `github-actions[bot]` (защита от рекурсии).
  Итог: PR смержился в master, но `deploy.yml` не запустился — прод остался на
  старом коде, пока не дёрнешь `workflow_dispatch` руками (так и было 2026-06-05).
  → завести `GH_PAT` (repo→Settings→Secrets) или сделать deploy явно зависимым от
  события merge, либо оставить деплой только ручным/по тегу.
- **🔸 [ops] Авто-конвейер dev→master→прод без апрува и без теста миграций.**
  `automerge.yml`/`deploy.yml`; CI гоняет только vet/build/unit. Сломанная
  миграция уезжает на живую БД, откат forward-only. → job с поднятием postgres и
  прогоном `runMigrations`; ручной апрув (environment protection) на прод.
### Среднее

- **🔸 [data-model] N+1 в листингах.** `aggregate.go:loadPlaces` зовёт
  `GetByID` (7 запросов) в цикле; списки `UserPlaceIDs`/`GemPlaceIDs` без LIMIT.
  → батч-загрузка / облегчённый `GetMany` без detail-полей.
- **🔸 [ops] Самописные миграции без леджера.** Хардкод-список в `main.go`,
  прогон на каждом старте, `*.down.sql` не применяются, идемпотентность вручную.
  → `golang-migrate`/`goose` или таблица `schema_migrations`.
- **✅→🔸 [backend] Suggest-прокси — таймаут СДЕЛАН (батч 2)** (`http.Client{Timeout:5s}`
  + контекст; батч 5: in-memory кэш TTL 60s). Дебаунс — на фронте уже есть (400ms).
- **✅ [backend] HTTP-сервер — таймауты + лимит тела СДЕЛАНЫ** (батч 2:
  `ReadHeaderTimeout`/`IdleTimeout`; батч 3: `MaxBytesReader` 1 MB на JSON).
- **✅ [security] Rate-limit СДЕЛАН (батч 3)** — login/register 20/min по IP,
  suggest 60/min по userID (`middleware/ratelimit.go`, in-memory token-bucket).
- **✅ [product] Город — канонизация СДЕЛАНА (батч 4)** (`CanonicalCity`: приведение
  к существующему написанию + trim). Полка «По городам» больше не дробится по
  регистру. Осталось 🔸 (опц.): справочник/автокомплит городов, rename города.
- **🔸 [product/bug] Wishlist «исполнение желания» не срабатывает, если место
  посетил не автор плана** (`MarkStruck` бьёт только по `user_id` визита), хотя
  `product.md` описывает фичу как круговую. → struck по факту визита кем угодно.
- **✅ [product] Поиск по name+city+cuisine СДЕЛАН (батч 3)** (`place.go:List`).
- **✅→🔸 [frontend] Классификатор Яндекса — `website` СДЕЛАН (батч 3).** Сайт
  из подсказки/ручного поля теперь сохраняется. Осталось 🔸: маппить `categories`
  в кухню/категории; `uri`/org-id — задел под фазу 3 идентичности места.
- **✅→🔸 [frontend] Обработка ошибок — частично (батч 2).** `PlaceDetail`:
  delete места/отзыва и toggle wishlist обёрнуты в try/catch + toast. Осталось
  🔸: единый враппер мутаций (паттерн повторяется по views).
- **🔸 [frontend] a11y подсказок `LocationPicker`** — невалидный HTML (`<li>` вне
  `<ul>`), нет `role`/`aria`. Образец рядом — `MultiSelect.vue`.

### Низкое

- 🔸 [data-model] Громоздкий place-SELECT продублирован 4× (List/GetByID/wishlist) и уже разошёлся по полям.
- 🔸 [backend] (✅ батч 2: лимит длины заметки; ✅ батч 5: review Update/Delete сверяет `place_id`); осталось: гонка TOCTOU на лимите 5 фото.
- 🔸 [security] CORS `*` + `AllowCredentials:true` (инертно — токен в заголовке); публичная перечислимая `/p/:id` (перебор id выгружает каталог) → шарить по UUID.
- 🔸 [ops] Нет лимитов ресурсов (mem/cpu) на общем VPS — нужна привязка к реальному потреблению. Деплой `git reset --hard` на проде. (✅ батч 4: backend healthcheck + nginx ждёт healthy; ✅ батч 5: ротация docker-логов)
- 🔸 [frontend] Разнобой копирайта «место/заведение/location»; Bootstrap + кастомный scrapbook-CSS дублируются. (✅ батч 4: мёртвый `VideoKruzhok.vue` удалён)
- 🔸 [tech-debt] `feed_events` покрывает только review+note; лента склеивается на клиенте из 3 запросов вместо единого `/feed`.

---

## 💤 Отложено осознанно

- **💤 Модель «бренд → филиал» (brands над places).** Правильная доменная
  модель, но для дружеского круга преждевременна (двухшаговая форма, авто-склейка
  одноимённых мест в разных городах). Вводить, когда филиалов одной сети станет
  столько, что доска без группировки шумит. Идентичность места уже чинит баг без
  этой сущности. Идеальный end-state — единая таблица `places`=филиал с приоритетом
  ключей: Яндекс `org_id` → координатный guard → текстовый `name+address+city`.
- **💤 Тёмная тема, page-transitions, «Ваш год»** — см. `ROADMAP.md` (next-next).
