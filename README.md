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
- `POST /api/auth/register` — регистрация `{username, email, password}`
- `POST /api/auth/login` — логин `{email, password}`
- `GET /api/auth/me` — профиль (JWT)

### Places
- `GET /api/places` — список (фильтры: city, cuisine_type_id, category_id, min_rating, is_gem, search, sort)
- `GET /api/places/cities` — список городов (уникальные)
- `GET /api/places/:id` — деталь
- `POST /api/places` — создать (auth)
- `PUT /api/places/:id` — обновить (auth, owner)
- `DELETE /api/places/:id` — удалить (auth, owner)
- `POST /api/places/:id/image` — загрузить фото (auth, owner, multipart/form-data, поле `image`)

### Reviews
- `GET /api/places/:id/reviews` — отзывы заведения
- `POST /api/places/:id/reviews` — создать отзыв (auth)
- `PUT /api/places/:id/reviews/:rid` — обновить (auth, author)
- `DELETE /api/places/:id/reviews/:rid` — удалить (auth, author)
- `POST /api/places/:id/reviews/:rid/image` — загрузить фото к отзыву (auth, author)
- `GET /api/users/:userId/reviews` — отзывы пользователя

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
│   ├── migrations/                  — SQL миграции (001–006)
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

```bash
# 1. Создать .env из шаблона
cp .env.example .env
# отредактировать .env — задать надёжные DB_PASSWORD и JWT_SECRET

# 2. Запустить
docker compose -f docker-compose.prod.yml up --build -d
```

Приложение будет доступно на порту `APP_PORT` (по умолчанию 80).
