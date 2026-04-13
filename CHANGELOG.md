# Changelog

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
