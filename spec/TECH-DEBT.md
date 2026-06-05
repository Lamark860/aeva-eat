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

---

## 🔸 Открытый тех-долг (по убыванию приоритета)

### Критичное / высокое

- **🔸 [ops] Бэкап на проде не на расписании.** Скрипт есть, но cron + off-site
  копию надо завести на VPS (`DEPLOY.md`). Пока не сделано — риск потери данных
  при отказе диска остаётся.
- **🔸 [product] Общее место правит/удаляет только создатель, нет admin-bypass.**
  `handler/place.go` Update/Delete/UploadImage проверяют только `ownerID==userID`.
  Чужое место с опечаткой никто (даже superuser) не исправит; при `created_by=NULL`
  место блокируется навсегда. Конфликт с философией «общая на круг» (`product.md`).
  → дать superuser-оверрайд или разрешить редактирование любому из круга.
- **🔸 [data-model] Удаление места каскадит чужие отзывы/фото/видео без
  предупреждения** (`001_init.sql` CASCADE, `place.go:Delete` голый DELETE).
  Файлы в `/uploads` остаются сиротами. → confirm с числом отзывов + мягкое
  удаление; чистить файлы.
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
- **🔸 [frontend] `ymaps` без проверки загрузки.** `index.html` без `onerror`,
  `LocationPicker`/`MapView` дёргают глобальный `ymaps`. Блокировщик/оффлайн →
  пустой div без сообщения, поиск молча мёртв. → guard + fallback на ручной ввод.

### Среднее

- **🔸 [data-model] N+1 в листингах.** `aggregate.go:loadPlaces` зовёт
  `GetByID` (7 запросов) в цикле; списки `UserPlaceIDs`/`GemPlaceIDs` без LIMIT.
  → батч-загрузка / облегчённый `GetMany` без detail-полей.
- **🔸 [ops] Самописные миграции без леджера.** Хардкод-список в `main.go`,
  прогон на каждом старте, `*.down.sql` не применяются, идемпотентность вручную.
  → `golang-migrate`/`goose` или таблица `schema_migrations`.
- **🔸 [backend] Suggest-прокси без таймаута/кэша** (`handler/suggest.go` —
  `http.Get`, без `Context`/`Timeout`). Зависший Яндекс копит горутины; каждое
  нажатие = платный вызов. → `http.Client{Timeout}`, кэш TTL, дебаунс.
- **🔸 [backend] HTTP-сервер без таймаутов и без лимита тела** (`main.go`
  `ListenAndServe`, нет `MaxBytesReader`). → `http.Server{ReadTimeout,…}`.
- **🔸 [security] Нет rate-limit на `/login` `/register` `/suggest`** — брутфорс
  и слив квоты Яндекса. → `chi/httprate` или nginx `limit_req`.
- **🔸 [backend] Загрузка видео доверяет `Content-Type`** без sniff magic-bytes
  (`review.go:UploadVideo`, в отличие от картинок). → проверять сигнатуру.
- **🔸 [product] Город — свободный текст без нормализации** → полка «По городам»
  дробится («Москва»/«москва»/«Moscow»). → canonical-case/справочник.
- **🔸 [product/bug] Wishlist «исполнение желания» не срабатывает, если место
  посетил не автор плана** (`MarkStruck` бьёт только по `user_id` визита), хотя
  `product.md` описывает фичу как круговую. → struck по факту визита кем угодно.
- **🔸 [product] Поиск ищет только по названию** (`place.go:List` — `LOWER(name)
  LIKE`), хотя плейсхолдер обещает «место, кухня, город». → расширить на
  city+cuisine (или сузить плейсхолдер).
- **🔸 [frontend] Данные классификатора Яндекса теряются.** `uri`/`categories`/
  `rating`/`url` доходят до фронта, но не сохраняются; `website` не пишется
  никогда. → прокинуть `website`, маппить категории; `uri` — задел под фазу 3
  идентичности (внешний org-id как ключ).
- **🔸 [frontend] Непоследовательная обработка ошибок** — `handleDelete`/
  `handleDeleteReview`/`wishlist.toggle` без `catch` → тихий провал. → единый
  враппер с toast.
- **🔸 [frontend] a11y подсказок `LocationPicker`** — невалидный HTML (`<li>` вне
  `<ul>`), нет `role`/`aria`. Образец рядом — `MultiSelect.vue`.

### Низкое

- 🔸 [data-model] Громоздкий place-SELECT продублирован 4× (List/GetByID/wishlist) и уже разошёлся по полям.
- 🔸 [backend] Правка заметки без лимита длины (в Create лимит есть); гонка TOCTOU на лимите 5 фото; review Update/Delete не сверяет `place_id` из URL.
- 🔸 [security] CORS `*` + `AllowCredentials:true` (инертно — токен в заголовке); публичная перечислимая `/p/:id` (перебор id выгружает каталог) → шарить по UUID.
- 🔸 [ops] Нет лимитов ресурсов и ротации docker-логов на общем VPS; у backend нет healthcheck в compose; деплой `git reset --hard` на проде.
- 🔸 [frontend] Разнобой копирайта «место/заведение/location»; мёртвый `VideoKruzhok.vue`; Bootstrap + кастомный scrapbook-CSS дублируются.
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
