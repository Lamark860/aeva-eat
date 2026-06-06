# AEVA Eat

Личный дневник впечатлений от заведений для друзей. Оценки кухни, сервиса, вайба, интерактивная карта, фильтры.

## Стек

| Компонент | Технология |
|-----------|-----------|
| Backend | Go 1.25, Chi router |
| Database | PostgreSQL 16 |
| Frontend | Vue 3, Vite, Bootstrap 5 |
| Карта | Яндекс Карты JS API 2.1 |
| Proxy | Nginx |
| Инфра | Docker Compose |

## Быстрый старт

```bash
docker compose up --build -d
```

Приложение: http://localhost:8091

## Порты

| Сервис | Порт (хост) | Порт (контейнер) |
|--------|-------------|-------------------|
| Nginx | 8091 | 80 |
| Backend API | 8086 | 8085 |
| Frontend (Vite dev) | 5173 | 5173 |
| PostgreSQL | 5434 | 5432 |

## API

### Health
- `GET /api/health` — healthcheck

### Auth
- `POST /api/auth/register` — регистрация `{username, password, invite_code}`
- `POST /api/auth/login` — логин `{username, password}`
- `GET /api/auth/me` — профиль (JWT)
- `PUT /api/auth/password` — смена пароля `{old_password, new_password}` (auth)
- `POST /api/auth/avatar` — загрузка аватарки (auth, multipart/form-data, поле `avatar`)

### Invites
- `GET /api/invites/validate/:code` — проверка инвайт-кода (публичный)
- `GET /api/invites` — мои инвайты (auth)
- `POST /api/invites` — создать инвайт (auth)
- `DELETE /api/invites/:id` — удалить инвайт (auth, только создатель или superuser)
- `GET /api/invites/all` — все инвайты (superuser)
- `GET /api/admin/users` — все пользователи (superuser)

### Places
- `GET /api/places` — список `{places: [...], total: N}`
  - фильтры: `city`, `cuisine_type_id`, `category_id`, `min_rating`, `is_gem`, `search` (по названию), `attended_by`, `visit_from`, `visit_to`
  - sort: `rating` | `rating_asc` | `name` | `rating_user:<userId>` (по умолчанию — новизна)
  - пагинация: `limit` (default 20, max 100, `0`=все), `page`
- `GET /api/places/cities` — список городов (уникальные)
- `GET /api/places/:id` — деталь (+ `gem_status`, `attendance`, `ratings_per_user`)
- `POST /api/places` — создать (auth). **Идентичность места = `name + address + city`** (уникальный индекс `idx_places_identity`, миграция 016). Разные филиалы сети (разный адрес) создаются нормально. На точный дубль → **`409 {error:"duplicate", existing:{…}}`** с телом уже существующего места, чтобы UI предложил перейти и оставить отзыв
- `PUT /api/places/:id` — обновить (auth, любой участник круга; конфликт идентичности → `409`)
- `DELETE /api/places/:id` — мягкое удаление (архив; auth, создатель/superuser). Отзывы/фото/видео круга сохраняются
- `POST /api/places/:id/restore` — вернуть место из архива (superuser)
- `POST /api/places/:id/image` — загрузить фото (auth, owner, multipart/form-data, поле `image`)
- `GET /api/random` — случайное место под фильтры (`?exclude_visited_by=me`)

### Reviews
- `GET /api/places/:id/reviews` — отзывы заведения
- `POST /api/places/:id/reviews` — создать отзыв (auth)
- `PUT /api/places/:id/reviews/:rid` — обновить оценки/коммент/дату (auth, author; фото и видео НЕ затрагиваются — у них свои эндпоинты)
- `DELETE /api/places/:id/reviews/:rid` — удалить (auth, author)
- `POST /api/places/:id/reviews/:rid/image` — добавить одно фото (legacy, auth, author)
- `POST /api/places/:id/reviews/:rid/photos` — добавить до 5 фото одним запросом (auth, author)
- `DELETE /api/places/:id/reviews/:rid/photos/:pid` — удалить фото (auth, author)
- `POST /api/places/:id/reviews/:rid/video` — загрузить видео-кружок (auth, author, mp4/webm, до 20MB)
- `GET /api/users/:userId/reviews` — отзывы пользователя

### Доска / Лента
- `GET /api/feed` — единая хронология (review + note); `GET /api/feed/weeks` — агрегаты по неделям
- `GET /api/feed/unread-count`, `POST /api/feed/seen` — индикатор непрочитанного
- `GET /api/notes`, `POST /api/notes`, `PUT /api/notes/:id`, `DELETE /api/notes/:id`, `PUT /api/notes/:id/strike` — записки от руки (auth, author)

### Разрезы (города / друзья / жемчужины)
- `GET /api/cities`, `GET /api/cities/:name`, `GET /api/cities/:name/places`, `GET /api/cities/:name/gems`
- `GET /api/users`, `GET /api/users/:id`, `GET /api/users/:id/places`, `GET /api/users/:id/gems`, `GET /api/users/:id/cities`
- `GET /api/gems` — hub жемчужин (+ агрегаты по городам/друзьям)

### Шаринг
- `GET /p/:token` — публичная preview-страница места по неугадываемому UUID-токену (без auth, OG-теги; без авторов/оценок). Шаринг по токену, не по id — нет энумерации

### Геосаджест
- `GET /api/suggest?text=...&ll=lat,lng` — подсказки Яндекс Геосаджест (проксируется через бэкенд, под auth)

### Справочники
- `GET /api/cuisine-types`
- `GET /api/categories`

### Wishlist
- `GET /api/wishlist` — мои запланированные заведения (auth)
- `GET /api/wishlist/ids` — ID запланированных (auth)
- `POST /api/wishlist/:id` — добавить в план (auth)
- `DELETE /api/wishlist/:id` — убрать из плана (auth)
- `GET /api/wishlist/custom` — кастомный список (auth)
- `POST /api/wishlist/custom` — добавить запись `{name, note?}` (auth)
- `DELETE /api/wishlist/custom/:id` — удалить запись (auth)

## Структура проекта

```
aeva-eat/
├── docker-compose.yml               — dev-окружение
├── docker-compose.prod.yml          — продакшн
├── .env.example                     — шаблон переменных
├── Makefile                         — lint, test, check, up/down
├── backend/
│   ├── cmd/api/main.go              — точка входа
│   ├── internal/
│   │   ├── config/                  — конфигурация из env
│   │   ├── handler/                 — HTTP хендлеры
│   │   ├── middleware/              — JWT auth
│   │   ├── model/                   — структуры данных
│   │   ├── repository/              — SQL запросы
│   │   └── service/                 — бизнес-логика
│   ├── migrations/                  — SQL миграции (001–016, накатываются на старте)
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/                     — axios instance
│   │   ├── components/              — переиспользуемые компоненты
│   │   ├── views/                   — страницы
│   │   ├── stores/                  — Pinia stores
│   │   ├── router/                  — Vue Router
│   │   └── composables/             — хуки
│   ├── Dockerfile                   — dev
│   └── Dockerfile.prod              — multi-stage production build
└── nginx/
    ├── nginx.conf                   — dev
    └── nginx.prod.conf              — production (gzip, кеш)
```

## Тесты

```bash
cd backend && go test ./...
```

## Продакшн-деплой

Полная пошаговая инструкция (перенос реальной БД, restore, Traefik/HTTPS) — в **[`DEPLOY.md`](./DEPLOY.md)**.

Кратко:

```bash
# 1. Создать .env из шаблона
cp .env.example .env
# задать: DB_PASSWORD, JWT_SECRET, APP_PORT, GEOSUGGEST_KEY, VITE_YANDEX_MAPS_KEY

# 2. Запустить
docker compose -f docker-compose.prod.yml up --build -d
```

Переменные окружения (`.env`):

| Переменная | Назначение |
|---|---|
| `APP_ENV` | `production` на проде — включает fail-fast при слабом `JWT_SECRET`. Локально `development` (дефолт) |
| `DB_PASSWORD` | пароль PostgreSQL |
| `JWT_SECRET` | секрет для подписи JWT (`openssl rand -base64 48`). В `production` пустой/дефолтный = отказ старта |
| `APP_PORT` | порт наружу (за Traefik — свой, напр. `8091`) |
| `GEOSUGGEST_KEY` | ключ Яндекс Геосаджеста (бэкенд, рантайм) |
| `VITE_YANDEX_MAPS_KEY` | ключ Яндекс JS API Карт (печётся во фронт на сборке) |

> Известные ограничения и тех-долг (бэкапы, права на редактирование общих мест,
> наблюдаемость, rate-limit и др.) задокументированы в **[`spec/TECH-DEBT.md`](./spec/TECH-DEBT.md)**.

Приложение будет доступно на порту `APP_PORT`.
