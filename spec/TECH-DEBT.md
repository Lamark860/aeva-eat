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

### ✅ Исправлено — батч 6 (2026-06-05, продуктовые развилки)

| Что | Где | Суть |
|---|---|---|
| Мягкое удаление мест | migration 017 (`deleted_at` + партиал-индексы), `SoftDelete`/`Restore`, List/Random/loadPlaces/share фильтруют | место уходит в архив, отзывы/фото/видео круга сохраняются; superuser возвращает через `POST /api/places/:id/restore` |
| Редактирование любым из круга | `handler/place.go` Update/UploadImage | общие места правит любой участник (по философии «без модерации»); удаление — создатель/superuser |
| Wishlist strike круговой | `handler/review.go`, `MarkStruckByPlace` | визит кем угодно гасит записку общего wishlist у всех (по спеке) |
| Auto-merge на rebase | `.github/workflows/automerge.yml` | `--squash` → `--rebase`: dev и master не расходятся, следующий PR не конфликтует |

### ✅ Исправлено — батч 7 (2026-06-05, safe-полировка)

| Что | Где | Суть |
|---|---|---|
| a11y + валидный HTML подсказок | `LocationPicker.vue` | `role=combobox/listbox/option`, `aria-activedescendant`, «Ничего не найдено» теперь `<li>` внутри `<ul>` (был невалидный HTML) |
| TOCTOU лимита фото | `repository/review.go` `AddPhoto` | `SELECT … FOR UPDATE` + проверка лимита в транзакции → параллельные загрузки не пробьют 5 фото (`ErrPhotoLimit`) |

---

## 🔸 Открытый тех-долг (по убыванию приоритета)

### Критичное / высокое

- **🔸 [ops] Бэкап на проде не на расписании.** Скрипт есть, но cron + off-site
  копию надо завести на VPS (`DEPLOY.md`). Пока не сделано — риск потери данных
  при отказе диска остаётся.
- **✅ [product] Права на места — СДЕЛАНО** (батч 2: superuser-bypass; батч 6:
  редактирование любым участником круга, удаление — создатель/superuser). Остаётся
  мелочь: переназначать `created_by` при NULL (сейчас не критично — мест без
  владельца в проде нет).
- **✅ [data-model] Удаление места — мягкое (батч 6).** `deleted_at` + Restore;
  чужой контент больше нельзя снести безвозвратно (hard-purge оставлен superuser'у
  на будущее). Мелкая неточность: счётчики-агрегаты («Москва · 87») пока включают
  архивированные места (редко; листинги их уже не показывают).
- **🔸 [ops] Авто-деплой после auto-merge молча не срабатывает** (нет `GH_PAT` →
  merge от `github-actions[bot]` не триггерит `deploy.yml`). Деплой дотягиваем
  `gh workflow run deploy.yml` руками. → завести `GH_PAT` (repo→Settings→Secrets).
  (✅ батч 6: squash→rebase убрал расхождение dev↔master, но это отдельная от
  GH_PAT проблема.)
- **✅→🔸 [ops] Тест миграций + апрув прода — СДЕЛАНО (батч 10).** CI-job
  `migrations` поднимает Postgres-service и реально бутит backend (прогон всех
  миграций + проверка леджера ≥16) — сломанный/не-идемпотентный SQL ловится до
  прода. `deploy.yml` получил `environment: production` — апрув включается заданием
  required reviewers в Settings→Environments (репо-настройка). Осталось 🔸: задать
  reviewers (твоё действие) и forward-only откат (по тегу/SHA).
### Среднее

- **✅ [data-model] N+1 в листингах — СДЕЛАНО (батч 8 + 11).** `List` («Найти»)
  батчит reviewers+feedPhotos (батч 8). `loadPlaces` (профиль/жемчужины) переведён
  на `GetManyByIDs` — облегчённая батч-загрузка (base + reviewers + feedPhotos +
  gem_status одним набором запросов вместо N×8); detail-only поля (attendance,
  ratings) карточкам не отдаются. Порядок входа сохранён, soft-deleted исключены.
  Проверено интеграционным тестом против реального Postgres (`place_db_test.go`,
  gated `AEVA_TEST_DSN`).
- **✅ [ops] Миграции с леджером — СДЕЛАНО (батч 9).** Таблица `schema_migrations`:
  каждый файл применяется максимум один раз, в транзакции (миграция + запись
  версии коммитятся вместе — нет частично-применённого состояния). На каждом
  старте больше не гоняется весь SQL. Проверено на реальном Postgres: чистый
  старт, идемпотентный рестарт, prod-сценарий (схема есть, леджера нет → разовый
  идемпотентный перенакат + запись, партиал-индекс сохранён). Остаётся 🔸 (опц.):
  применение `*.down.sql` для отката — пока вручную.
- **✅→🔸 [backend] Suggest-прокси — таймаут СДЕЛАН (батч 2)** (`http.Client{Timeout:5s}`
  + контекст; батч 5: in-memory кэш TTL 60s). Дебаунс — на фронте уже есть (400ms).
- **✅ [backend] HTTP-сервер — таймауты + лимит тела СДЕЛАНЫ** (батч 2:
  `ReadHeaderTimeout`/`IdleTimeout`; батч 3: `MaxBytesReader` 1 MB на JSON).
- **✅ [security] Rate-limit СДЕЛАН (батч 3)** — login/register 20/min по IP,
  suggest 60/min по userID (`middleware/ratelimit.go`, in-memory token-bucket).
- **✅ [product] Город — канонизация СДЕЛАНА (батч 4)** (`CanonicalCity`: приведение
  к существующему написанию + trim). Полка «По городам» больше не дробится по
  регистру. Осталось 🔸 (опц.): справочник/автокомплит городов, rename города.
- **✅ [product/bug] Wishlist «исполнение желания» — СДЕЛАНО (батч 6).** Визит кем
  угодно из круга гасит записку общего wishlist у всех (`MarkStruckByPlace`).
- **✅ [product] Поиск по name+city+cuisine СДЕЛАН (батч 3)** (`place.go:List`).
- **✅→🔸 [frontend] Классификатор Яндекса — `website` СДЕЛАН (батч 3).** Сайт
  из подсказки/ручного поля теперь сохраняется. Осталось 🔸: маппить `categories`
  в кухню/категории; `uri`/org-id — задел под фазу 3 идентичности места.
- **✅→🔸 [frontend] Обработка ошибок — частично (батч 2).** `PlaceDetail`:
  delete места/отзыва и toggle wishlist обёрнуты в try/catch + toast. Осталось
  🔸: единый враппер мутаций (паттерн повторяется по views).
- **✅ [frontend] a11y подсказок `LocationPicker` — СДЕЛАНО (батч 7)**: валидный
  HTML + `role`/`aria-activedescendant`.

### Низкое

- ✅→🔸 [data-model] Place-SELECT: List/GetByID/GetManyByIDs сведены к общим `placeSelectCols`+`placeBaseFrom`+`scanPlace` (батч 11) — дрейф List↔GetByID устранён. Остаётся 🔸: 2 более лёгкие копии в `wishlist.go` (без videos/top_comment) — консолидация добавила бы им ненужные колонки, оставлены намеренно.
- 🔸 [backend] (✅ батч 2: лимит длины заметки; ✅ батч 5: review сверяет `place_id`; ✅ батч 7: TOCTOU лимита фото устранён через FOR UPDATE).
- ✅→🔸 [security] Публичная `/p/` — СДЕЛАНО (батч 12): шарится по неугадываемому UUID `share_token` (миграция 018), не по инкрементному id — энумерация каталога закрыта; nginx-регэксп под UUID, кнопка «поделиться» в PlaceDetail. Остаётся 🔸: CORS `*` + `AllowCredentials:true` (инертно — токен в заголовке, не в куках).
- 🔸 [ops] Деплой `git reset --hard` на проде (forward-only). (✅ батч 4: backend healthcheck + nginx ждёт healthy; ✅ батч 5: ротация docker-логов; ✅ батч 10: лимиты mem/cpu — потолки postgres 768m/backend 512m/front+nginx 128m)
- ✅→🔸 [frontend] Копирайт унифицирован на «место» (батч 13: все user-facing «заведение»→«место»). Остаётся 🔸: внутреннее имя компонента `LocationPicker` (не пользовательское, рефактор-рейн пропущен); Bootstrap + scrapbook-CSS сосуществуют (полное вытеснение Bootstrap — крупный риск, отложено). (✅ батч 4: мёртвый `VideoKruzhok.vue` удалён)
- 💤 [tech-debt] Единый серверный `/feed` — **осознанно НЕ делаем** (см. ниже): лента собирается на клиенте (BFF-enrichment) намеренно.

---

## 💤 Отложено осознанно

- **💤 Единый серверный `/feed` вместо клиентской сборки.** Разобрано в батче 13:
  `useFeed.js` держит плотную ПРЕЗЕНТАЦИОННУЮ логику доски — недельные/месячные
  бакеты, группировку по (place_id, неделя), усреднение оценок, сбор видео, резолв
  attendees, выбор Q-layout-цитаты, мердж note+wishlist. Перенос в API запёк бы
  вёрстку доски в бэкенд и означал высокорисковый рефактор сердца приложения при
  сомнительной выгоде на камерном масштабе. Текущий дизайн (хронология из
  `feed_events` + клиентское обогащение из /places,/notes,/wishlist,/users) —
  оправданный BFF-паттерн, оставлен намеренно. Если доска вырастет — выносить
  бакетинг в API отдельной взвешенной работой.

- **💤 Модель «бренд → филиал» (brands над places).** Правильная доменная
  модель, но для дружеского круга преждевременна (двухшаговая форма, авто-склейка
  одноимённых мест в разных городах). Вводить, когда филиалов одной сети станет
  столько, что доска без группировки шумит. Идентичность места уже чинит баг без
  этой сущности. Идеальный end-state — единая таблица `places`=филиал с приоритетом
  ключей: Яндекс `org_id` → координатный guard → текстовый `name+address+city`.
- **💤 Тёмная тема, page-transitions, «Ваш год»** — см. `ROADMAP.md` (next-next).
