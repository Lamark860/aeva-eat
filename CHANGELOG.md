# Changelog

## [0.21.1] — 2026-05-10 — Полировка кружков, FOUC, группировка фида

### Исправлено
- **Видео-кружок рендерился квадратом** через transformed parent — добавлен `border-radius: 50%` на сам `<video>`, Chrome перестал клипать
- **Пустой белый кружок без видео** — `hasKruzhok` гейтится на реальный `video_url`, не на флаг `has_video`
- **Ticket-only артефакт** растягивался на всю строку при наличии видео — `max-width: 360px`
- **Per-event рейтинги и дата** — каждая карточка дня отражает рейтинги именно этого визита, а не среднее по месту
- **▶-индикатор** на видео-кружке — тёмный disc с backdrop-blur, белый CSS-треугольник, всегда центрирован
- **FOUC Bootstrap-навбара** на первом рендере — `await router.isReady()` перед `app.mount()`, paper-bg на body
- **Дублирование архивных секций** — потерянный `v-else` в expanded-секции восстановлен
- **iOS/Chrome poster** для видео — `@loadedmetadata="forcePoster"` сикает на 0.1s

### Добавлено
- **Inline-play видео-кружков** — тап играет/паузит в том же DOM-узле, без перехода на /places/:id
- **Стопка кружков справа от полароида** — мульти-видео: до 3 видимых, остальные → рукописное `+N`, offset `top: i*38px`
- **Группировка по `(place_id, day)`** — два отзыва одного места в один день мерджатся в одну карточку с union-attendees

## [0.21.0] — 2026-05-10 — Раунд 5: ответы дизайнера, wishlist на Доске, шаринг /p/:id

### Добавлено
- **WishlistArtifact** на Доске: активный (paper + терра-штамп «план» + канцелярская кнопка), зачёркнутый (paper-deep + moss-штамп «сходили ✓» + SVG-волнистый штрих + мини-полароид визита)
- **`/p/:id`** — публичная preview-страница для шаринга в мессенджерах (Go-handler без auth, `html/template`, OG/Twitter meta). Без рейтингов и без имён авторов
- **Расширенный `/api/places/:id`**: `gem_status` (отметил X · дата + соавторы), `attendance` (visit_count), `ratings_per_user` (бэк-фундамент). UI: рукописная подпись «отметил(а)…» под штампом жемчужины, ряд аватарок «×N» под общим тикетом
- **Расширенные фильтры в Найти**: drawer «кто был» (avatar-chip multi-select из друзей, поиск при >10), «когда» (date-range + presets «этот год / прошлый год / последние 30 дней»), sort «по оценке друга»
- **Любимая кухня** в профиле: SQL-агрегат + `composables/useCuisine.js` (cuisineAccusative, razSuffix). Рендер «любит грузинскую — 11 раз» Caveat-шрифтом. Прячется при count<2
- **`useFeed` event-driven** — параллельная загрузка places/notes/wishlist/users, группировка review-событий, аватарки через userMap
- **AddArtifactSheet почищен** — убран disabled-пункт «wishlist · скоро». Теперь визит + записка
- **Анимации MVP**: артефакт fade+slide-up 280мс, разворот недели height+stagger 320мс, gem-stamp при включении

### Изменено
- **`color-scheme: light`** зафиксирован в `:root` — Safari/Chrome перестали дарк-модить form controls на OS dark mode (подготовка к G2)
- 30+ литеральных hex-значений в компонентах заменены на CSS-переменные

## [0.20.0] — 2026-05-09 — Раунд 4: бэк-дельта + новые экраны

### Добавлено
- **Миграция 012** `review_photos` (несколько фото в review): таблица с `position`, FK CASCADE, backfill 9 `image_url` → `position=0`. API: `POST .../photos` (multipart, до 5 за раз), `DELETE .../photos/:pid`. `image_url` сохранён для backwards-compat и cover-fallback
- **Миграция 013** `notes` + VIEW `feed_events` (review_added + note_added, фундамент для будущих типов). API: `GET/POST/PUT/DELETE /api/notes`, `PUT /api/notes/:id/strike`, `GET /api/feed`, `GET /api/feed/weeks`
- **Миграция 014** `users.last_seen_feed_at` + `GET /api/feed/unread-count` + `POST /api/feed/seen`. Точка-индикатор «новости» на табе Доска (poll каждые 60s + при focus)
- **Миграция 015** `wishlists.is_struck` + `struck_at`, триггер при создании review автоматически помечает wishlist как зачёркнутый. API `GET /api/wishlist/all`
- **Cities API**: `GET /api/cities`, `/cities/:name`, `/cities/:name/places`, `/cities/:name/gems`
- **Users API**: `GET /api/users`, `/users/:id`, `/users/:id/places`, `/users/:id/gems`, `/users/:id/cities`
- **Gems Hub API**: `GET /api/gems` (places + by_city + by_user)
- **`GET /api/random`** с фильтрами `city`, `cuisine_type_id`, `is_gem`, `exclude_visited_by=me|<id>`
- **`/api/places` расширен**: `attended_by=1,3,7`, `visit_from/to`, `sort=rating_user:<id>` (NULLS LAST)
- **Новые экраны**: `views/CityPage.vue` (/cities/:name), `views/PersonPage.vue` (/people/:id), `views/GemsHub.vue` (/gems) — серифа-имя, билетик-стата, полки `ResultCard`
- **PolaroidStack.vue** — 1 фото одиночка, 2 внахлёст −2°/+3°, 3+ стопка из 3-х с каракулевым «+ ещё N»
- **Onboarding-режим** (C1): юзер с places_count===0 видит read-only ленту + soft-CTA «+ новое место»
- **Доска: featured-артефакт** (A1) — grid с dense flow, один полароид во всю ширину раз в бакет
- **Полка «По друзьям»** в Найти — горизонтальная карусель аватарок 60px

### Изменено
- Доска переведена с flex на `display: grid` + `grid-auto-flow: dense`
- `repository/aggregate.go` — все SQL-агрегации в одном файле

## [0.19.0] — 2026-05-08 — Phase 7: хендофф дизайнера

### Добавлено
- **`spec/DESIGN-DECISIONS.md`** — 19 ответов дизайнера verbatim, M1/M2/F1/R2/L4/R3/G1 микро-улучшения реализованы
- **`spec/NEXT.md`** — продуктовые заметки дизайнера после хендоффа (A1–A3, B1–B5, C1–C6)
- **M1** маркер карты охрой для wishlist (без рейтинга — терра = «плохо» в светофоре)
- **M2** стек авторов в balloon
- **F1** список городов вертикальный
- **L4** paper-deep фон для развёрнутого архива
- **R3** stamp-press анимация для gem-toggle
- **G1** фон вокруг колонки на десктопе

## [0.18.0] — 2026-05-07 — Phase 6: миграция на скрапбук-дизайн

### Добавлено
- **`scrapbook.scss`** — палитра OKLCH в CSS-переменных, бумажный фон с грэном, sb-* примитивы, Lora + Caveat шрифты
- **Vue-примитивы**: `Polaroid`, `Tape`, `Stamp`, `Ticket` (food/service/vibe), `Note`, `GemBadge` (анимация блика), `AuthorTag` (4 цвета terra/ochre/moss/plum), `VideoKruzhok`, `PinButton`, `CollapsedStrip`, `AddArtifactSheet`, `ArtifactCard`, `ResultCard`
- **Доска** (`views/Home.vue`): wordmark, 2-колоночная масонри-сетка, шапка недели + PinButton, «↓ ещё N», свернутые полоски прошлых недель/месяцев
- **Перевёрстаны на скрапбук**: PlaceDetail, Places (Найти), Profile, PlaceForm, Login, InviteRegister, Invites, MapPage, MapView (маркеры-канцелярки и balloon-плашки), MultiSelect, VideoRecorder, RatingInput, ReviewForm, ReviewCard
- **`/login`, `/invite/:code`, `/invites`, `/map`** с `meta.scrapbook=true` → навбар Bootstrap скрыт, рендер edge-to-edge
- **Светофор маркеров карты** — pushpin head меняет цвет по рейтингу (≥8 moss / ≥5 ochre / <5 terra), рейтинг внутри головки

### Удалено
- `components/PlaceCard.vue` — legacy Bootstrap, заменён ResultCard/ArtifactCard

### Исправлено
- EXIF orientation — `imaging.Decode(AutoOrientation(true))`, iPhone-фото больше не «лёжа»
- Cover-fallback на первое review-фото (`COALESCE(p.image_url, review.image_url)`)
- VideoRecorder на iOS Safari ≤16 — выбор доступного mimeType из кандидатов
- 401 redirect loop, SW kill-switch, nginx WS upgrade

## [0.17.0] — 2026-05-05 — Phase 5: PWA

### Добавлено
- **Манифест** `manifest.webmanifest` + иконки 192/512/maskable + apple-touch-icon
- **Service Worker** `sw.js` с кэшем shell и стратегией network-first для API
- Регистрация SW в `main.js`

## [0.16.0] — 2026-05-04 — Phases 1-4: мобильная адаптация

### Добавлено
- **Phase 1**: viewport meta, scss-хелперы для мобильных, safe-area-insets
- **Phase 2**: `BottomTabBar` (Доска / Найти / Карта / Я / Добавить)
- **Phase 3**: filter drawer (Bootstrap offcanvas со скрапбук-обёрткой) для Places и MapPage
- **Phase 4**: touch-friendly inputs, stacked CTAs на мобильных
- Fix: sort по рейтингу падал с 500 (DISTINCT + ORDER BY конфликт)

## [0.15.0] — 2026-04-17 — Профиль, пагинация, аватар

### Добавлено
- **Редактирование / удаление отзывов из профиля**: inline-редактирование через ReviewForm, кнопки ✏️ и 🗑 на каждом отзыве, подтверждение удаления
- **Смена пароля**: `PUT /api/auth/password` — валидация старого пароля, минимум 6 символов; UI-форма в профиле
- **Загрузка аватарки**: `POST /api/auth/avatar` — загрузка с автоматическим сжатием через imageutil; круглый аватар в профиле с hover-оверлеем, плейсхолдер с первой буквой имени
- **Пагинация заведений**: `GET /api/places` возвращает `{places, total}`, поддержка `limit` / `page` параметров
  - Фронтенд: навигация по страницам (до 5 видимых), сброс на страницу 1 при смене фильтров
  - Карта: отдельный запрос `fetchAllPlaces()` с `limit=0` для загрузки всех точек

### Изменено
- Place handler теперь возвращает объект `{places: [...], total: N}` вместо голого массива
- Профиль: блок юзера дополнен аватаром и кнопкой смены пароля

## [0.14.0] — 2026-04-16 — Приватный доступ, сжатие фото, видеосообщения

### Добавлено
- **Invite-only система**: самостоятельная регистрация отключена, доступ только по инвайт-коду
  - Миграция `009_invites_roles` — таблица `invites` + колонка `role` (superuser/user) в `users`
  - API: `POST /api/invites` — создание инвайта, `GET /api/invites` — мои инвайты, `GET /api/invites/validate/:code` — проверка кода (публичный)
  - Регистрация `POST /api/auth/register` теперь требует `invite_code`
  - Admin API: `GET /api/admin/users`, `GET /api/invites/all` (только superuser)
  - Страница `/invite/:code` — регистрация по приглашению с валидацией
  - Страница `/invites` — управление инвайтами (генерация ссылок, копирование)
- **Router guards**: все маршруты кроме `/login` и `/invite/:code` требуют авторизацию
- **Все API-роуты закрыты JWT**: places, reviews, catalogs, wishlist, geosuggest — без токена 401
- **Сжатие фото**: загружаемые изображения автоматически:
  - Ресайзятся до max 1920×1920px (с сохранением пропорций)
  - Конвертируются в JPEG с quality 80%
  - EXIF-метаданные удаляются (decode → re-encode)
  - Библиотека `disintegration/imaging`
- **Видеосообщения к отзывам** (как в Telegram):
  - Миграция `010_review_video` — колонка `video_url` в `reviews`
  - `POST /api/places/:id/reviews/:rid/video` — загрузка видео (mp4/webm, до 20МБ)
  - `VideoRecorder` компонент: запись с камеры (circular preview), лимит 60 сек, MediaRecorder API
  - Круглое превью видео в `ReviewCard` с play/pause по клику
  - Интеграция в `ReviewForm` — кнопка «Записать видео» + предпросмотр

### Изменено
- Navbar: показывает ссылки (Заведения, Карта, Пригласить) только для авторизованных
- Убрана ссылка на открытую регистрацию из Login.vue
- Nginx: `client_max_body_size` увеличен до 25MB (для видео)
- `YANDEX_MAPS_RULES.md` обновлён — заметка про invite-only и ToS

## [0.13.0] — 2026-04-13 — UI-фиксы и полировка

### Исправлено
- **Автозаполнение названия**: для гео-результатов (`type !== 'business'`) берётся оригинальное имя из Geosuggest, а не адрес улицы из ymaps.search
- **Счётчик жемчужин** на главной: добавлено поле `is_gem_place` (вычисляемое через `EXISTS`) в модель Place и SQL-запросы List/GetByID
- **Профиль — имя заведения**: `ListByUser` теперь JOIN-ит таблицу `places` и возвращает `place_name` вместо «Заведение #ID»
- **Сортировка по рейтингу**: бэкенд теперь поддерживает `sort=rating` (убывание) и `sort=rating_asc` (возрастание)

### Добавлено
- **Превью общего рейтинга** в форме отзыва: динамический показатель `(food + service + vibe) / 3` с цветовой индикацией (зелёный/жёлтый/красный)
- **Сворачивание отзывов** в профиле: по умолчанию показываются 5, кнопка «Показать все / Свернуть»
- **Tooltip на маркерах карты**: `hintContent` — при наведении показывается название заведения
- **Сортировка «По рейтингу ↑»** в фильтрах списка заведений

## [0.12.0] — 2026-04-13 — Геосаджест, данные, уникальность

### Добавлено
- **Geosuggest-прокси**: `GET /api/suggest` проксирует запросы к Яндекс Геосаджест (ключ на бэкенде, без CORS-проблем)
- **Autocomplete в PlaceForm**: dropdown из Geosuggest → выбор → ymaps.search для координат и метаданных
- **Сид-данные**: 174 заведения (Казань + Нижний Новгород), 141 совместный отзыв, 2 пользователя (alina, lamark)
- **Уникальность заведений**: миграция 008 — `UNIQUE INDEX` по `(LOWER(name), LOWER(city))`, бэкенд возвращает понятную ошибку при дубле
- **Wishlist чекбокс** в форме создания заведения («🤍 Хочу сходить»)

### Изменено
- **Авторизация по username** (email полностью удалён из модели, репозиториев и UI)
- **PlaceForm переработан**: карта-поиск первым экраном, авто-заполнение карточки, ручные поля под спойлером
- Удалён лейбл «были» с PlaceCard — остались только аватары рецензентов

## [0.11.0] — 2026-04-13 — Миграция на Яндекс Карты

### Изменено
- **Карта**: Leaflet + CartoDB Voyager → Яндекс Карты JS API 2.1
  - MapView.vue: полностью переписан на `ymaps.Map` + `ymaps.Placemark` с inline SVG-маркерами
  - LocationPicker.vue: переписан на `ymaps.Map` + `ymaps.geocode` (прямое и обратное геокодирование)
  - Удалена npm-зависимость `leaflet`
  - Скрипт Яндекс Карт подключён в index.html
- **Футер**: добавлена ссылка на условия использования Яндекс Карт (обязательная по ToS)
- Добавлен документ `YANDEX_MAPS_RULES.md` — выжимка правил использования API

## [0.10.0] — 2026-04-13 — Phase 10: Фото отзывов, кастомный wishlist, аватары

### Добавлено
- **Фото к отзывам**:
  - Миграция `006` — колонка `image_url` в таблице `reviews`
  - `POST /api/places/:id/reviews/:rid/image` — загрузка фото к отзыву (auth, author)
  - ReviewCard: фото отображается над текстом отзыва
  - ReviewForm: drag & drop / выбор фото при создании/редактировании отзыва
- **Кастомный wishlist** (свободный текст):
  - Таблица `wishlist_custom` (name, note, user_id)
  - API: `GET/POST /api/wishlist/custom`, `DELETE /api/wishlist/custom/:id`
  - Профиль: секция «Ещё хочется попробовать» — форма для записи заведений не из системы
- **Аватары рецензентов на PlaceCard**: стопка маленьких кружков с инициалами (при наведении — полный список имён), бейдж «были»
- **Совместные отзывы** отмечены бейджем «совместный» в ReviewCard, авторы с аватарами
- **Reviewers в API**: эндпоинт `/api/places` теперь возвращает `reviewers` — список уникальных авторов отзывов (id + username)

### Исправлено
- **Маркеры карты** — заменены CSS-пины на inline SVG: цветная капля с рейтингом всегда видна, не зависит от порядка загрузки CSS
- Мок-данные: 6 отзывов с фото (Unsplash), 4 кастомных записи wishlist

### Изменено
- Мок-данные расширены: фото еды в отзывах, кастомные wishlist записи

## [0.9.0] — 2026-04-13 — Phase 9: Оценки float, Wishlist, Мок-данные

### Добавлено
- **Рейтинги 0–10 с шагом 0.1** (float):
  - Миграция `003_rating_float` — `NUMERIC(3,1)`, CHECK 0–10
  - Бэкенд: `float64` для всех рейтингов, валидация 0–10
  - Фронтенд: `RatingInput` — слайдер min=0, max=10, step=0.1
- **Wishlist (Хочу посетить)**:
  - Миграция `004_wishlist` — таблица `wishlists` (user_id, place_id)
  - Бэкенд: `WishlistHandler` + `WishlistRepo` — CRUD, `/api/wishlist`, `/api/wishlist/ids`
  - Фронтенд: Pinia store `wishlist.js`, кнопка ❤️ на PlaceCard и PlaceDetail, секция «Хочу посетить» в профиле
- **Мок-данные** — миграция `005_seed_data`:
  - 3 пользователя (alice, bob, charlie)
  - 8 заведений (Москва + СПб) с фотографиями (Unsplash)
  - 10 отзывов: сольные и совместные (alice+bob, bob+charlie, alice+charlie)
  - 4 жемчужины 💎
  - 5 записей wishlist
- **Кастомные маркеры карты**: SVG-пин с цветом по рейтингу (зелёный ≥8, жёлтый ≥5, красный <5), бейдж gem, стилизованные popup с фото
- **Тесты**: `wishlist_test.go` (unauthorized endpoints), расширен `review_test.go` (10 boundary cases)

### Изменено
- Убран суффикс «/10» у оценки на карточке и в слайдере рейтинга — компактнее
- Маркеры карты крупнее (40×50px), белая обводка, тень — лучше видно на карте
- Размер общей оценки на PlaceCard уменьшен до 1.1rem — гармоничнее с текстом

## [0.8.0] — 2026-04-13 — Phase 8: UI/UX Overhaul

### Добавлено
- **MultiSelect** компонент — мультивыбор кухонь и категорий в фильтрах:
  - Кастомный dropdown с тегами, поддержка нескольких значений одновременно
  - Используется в фильтрах списка заведений и карты
- Бэкенд: поддержка множественных `cuisine_type_id` и `category_id` через запятую (`IN` вместо `=`)
- Общая (суммарная) оценка заведения — среднее из кухня/сервис/вайб:
  - Отображается крупно в PlaceCard, PlaceDetail и ReviewCard
- Drag & drop зона для загрузки фото с превью
- Фото можно прикрепить при создании заведения (автозагрузка после создания)

### Изменено
- **Карта**: тайлы CartoDB Voyager вместо стандартных OSM (чище, современнее)
- **Цветовая схема**: контрастнее — белые карточки на `#f5f5f7` фоне, чёткие тени
- **Navbar**: светлый (белый с blur) вместо чёрного тёмного, оранжевый бренд
- **RatingInput**: слайдер с цветной заливкой вместо 10 квадратных кнопок, удобен на мобиле
- **Формат дат**: `2026-04-10T00:00:00Z` → «10 апреля 2026 г.» (locale `ru-RU`)
- **ReviewCard**: показывает общую оценку, дата в человеческом формате
- **Profile**: дата посещения в человеческом формате

## [0.7.0] — 2026-04-13 — Phase 7: UI/UX Polish

### Добавлено
- Кастомная SCSS-тема:
  - Тёплая палитра (primary #e85d30), шрифт Inter, мягкие скругления
  - Фон `#faf7f5`, бордерлесс карточки с тенями, стилизованный скроллбар
  - Заменён дефолтный Bootstrap CSS на кастомную SCSS-сборку
- **LocationPicker** компонент — выбор координат заведения:
  - Клик по карте → маркер (перетаскиваемый)
  - Поиск по адресу через Nominatim (OpenStreetMap geocoder)
  - Обратный геокодинг (клик → адрес)
  - Автозаполнение города из результата геокодинга
  - Заменены числовые поля широта/долгота в форме заведения
- **ToastContainer** + `useToast` composable — система уведомлений:
  - Toast-уведомления при создании/обновлении/удалении заведений и отзывов
  - 4 типа: success, error, info, warning с автоскрытием
- Анимации переходов между страницами (`<transition name="page">`)

### Изменено
- **Home** — hero-секция + карточки статистики (заведения, отзывы, жемчужины) + последние заведения
- **PlaceCard** — hover-эффект (подъём + zoom изображения), плейсхолдер без фото, бейдж жемчужины, вся карточка кликабельна
- **App.vue** — обновлённый navbar с аватаром-инициалом пользователя, кнопка регистрации стилизована как primary
- Navbar: полупрозрачный с backdrop-blur, убран класс `bg-dark`

## [0.6.0] — 2026-04-13 — Phase 6: Фото, адаптивность, продакшн

### Добавлено
- Загрузка фото заведений:
  - `POST /api/places/:id/image` — загрузка изображения (JPEG/PNG/WebP, до 5 МБ)
  - Миграция `002_place_image` — колонка `image_url` в таблице `places`
  - Фото отображается на карточке, странице заведения и в форме редактирования
  - Nginx раздаёт `/uploads/` как статику с кешированием
  - Docker volume `uploads_data` для персистентного хранения
- Адаптивная вёрстка:
  - Бургер-меню в навбаре (mobile)
  - Респонсивная сетка карточек (1/2/3 колонки)
- Продакшн-конфигурация:
  - `frontend/Dockerfile.prod` — multi-stage build (node → nginx)
  - `docker-compose.prod.yml` — env-файл, restart-политики, без volume-маунтов кода
  - `nginx/nginx.prod.conf` — gzip, кеширование, проксирование
  - `.env.example` — шаблон переменных окружения

## [0.5.0] — 2026-04-13 — Phase 4+5: Фильтры + Карта

### Добавлено
- Go API:
  - `GET /api/places/cities` — динамический список городов
- Vue frontend:
  - **MapView** компонент (Leaflet + OpenStreetMap): маркеры, popup с названием/оценкой/ссылкой, автоподгонка bounds
  - **Страница /map** — полноэкранная карта с фильтрами (город, кухня, категория, жемчужина, поиск)
  - **Мини-карта** на странице заведения (если есть координаты)
  - Фильтры синхронизируются с URL-параметрами (можно шарить ссылку с фильтрами)
  - Города подтягиваются с бэкенда динамически
  - Ссылка «Карта» в навбаре
  - Жемчужины отмечены особым маркером 💎 на карте
- Leaflet 1.9 подключён через npm

## [0.3.0] — 2026-04-13 — Phase 3: Reviews + Ratings

### Добавлено
- Go API:
  - `GET /api/places/:id/reviews` — отзывы заведения
  - `POST /api/places/:id/reviews` — создание отзыва (auth, автоматическое добавление себя в авторы)
  - `PUT /api/places/:id/reviews/:rid` — редактирование (auth, author only)
  - `DELETE /api/places/:id/reviews/:rid` — удаление (auth, author only)
  - `GET /api/users/:userId/reviews` — отзывы пользователя
  - Валидация рейтингов 1-10
  - M2M review_authors: несколько авторов на один отзыв
- Go тесты: validateRatings (boundaries, edge cases) — всего 18 тестов
- Vue frontend:
  - Компонент RatingInput (кликабельная шкала 1-10 с цветами)
  - Компонент ReviewCard (отображение отзыва с авторами)
  - Компонент ReviewForm (создание/редактирование отзыва)
  - Отзывы на странице заведения (создание, редактирование, удаление)
  - Средние оценки пересчитываются автоматически
  - Страница профиля с лентой отзывов пользователя
  - Ссылка на профиль в навбаре
  - Pinia store: reviews

## [0.2.0] — 2026-04-13 — Phase 2: Places + Справочники

### Добавлено
- Go API:
  - `GET /api/places` — список с фильтрами (city, cuisine_type_id, category_id, min_rating, is_gem, search, sort)
  - `GET /api/places/:id` — деталь заведения с категориями и средними оценками
  - `POST /api/places` — создание (auth, с привязкой категорий)
  - `PUT /api/places/:id` — редактирование (auth, owner only)
  - `DELETE /api/places/:id` — удаление (auth, owner only)
  - `GET /api/cuisine-types` — справочник типов кухонь
  - `GET /api/categories` — справочник категорий
- Go тесты: auth handler, auth service, JWT middleware (11 тестов)
- Vue frontend:
  - Страница списка заведений с фильтрами и поиском
  - Форма создания/редактирования заведения
  - Страница детали заведения с оценками
  - Pinia stores: places, catalogs
  - Компонент PlaceCard
  - Навигация «Заведения» в navbar
- README.md с документацией API и структуры проекта
- CHANGELOG.md

## [0.1.0] — 2026-04-13 — Phase 1: Скелет + Auth

### Добавлено
- Docker Compose: postgres, backend (Go), frontend (Vue), nginx
- PostgreSQL миграция: users, cuisine_types, categories, places, place_categories, reviews, review_authors
- Seed-данные: 15 типов кухонь, 10 категорий
- Go API (Chi router):
  - `POST /api/auth/register` — регистрация
  - `POST /api/auth/login` — логин (JWT)
  - `GET /api/auth/me` — профиль по токену
  - `GET /api/health` — healthcheck
  - JWT middleware
- Vue 3 frontend (Vite + Bootstrap 5):
  - Страницы: Home, Login, Register
  - Pinia auth store с localStorage persistence
  - Axios instance с JWT interceptor
  - Navbar с состоянием авторизации
- Nginx reverse proxy: `/api` → backend, `/` → frontend
