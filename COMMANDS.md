# Полезные команды AEVA Eat

## Docker Compose

```bash
# Поднять всё
docker compose up -d

# Поднять всё с пересборкой
docker compose up --build -d

# Пересобрать только бэкенд
docker compose up -d --build backend

# Пересобрать только фронтенд
docker compose up -d --build frontend

# Остановить всё
docker compose down

# Остановить + удалить volumes (БД тоже!)
docker compose down -v

# Перезапустить конкретный сервис
docker compose restart backend
docker compose restart frontend
docker compose restart nginx
```

## Логи

```bash
# Все логи (follow)
docker compose logs -f --tail=50

# Только бэкенд
docker compose logs -f --tail=50 backend

# Только фронтенд
docker compose logs -f --tail=50 frontend

# Только nginx
docker compose logs -f --tail=50 nginx

# Только postgres
docker compose logs -f --tail=50 postgres
```

## БД (PostgreSQL)

```bash
# Зайти в psql
docker exec -it aeva-postgres psql -U aeva -d aeva_eat

# Быстрый SQL-запрос
docker exec -it aeva-postgres psql -U aeva -d aeva_eat -c "SELECT count(*) FROM places;"

# Дамп базы
docker exec aeva-postgres pg_dump -U aeva aeva_eat > dump.sql

# Восстановить из дампа
cat dump.sql | docker exec -i aeva-postgres psql -U aeva -d aeva_eat
```

Параметры подключения к БД:
| Параметр | Значение |
|----------|----------|
| Host     | localhost |
| Port     | 5434 |
| Database | aeva_eat |
| User     | aeva |
| Password | aeva_secret |

## Тесты и линтеры

```bash
# Go-тесты
cd backend && go test ./... -v

# Go vet
cd backend && go vet ./...

# Vue ESLint (через контейнер)
docker exec aeva-frontend npx eslint src --ext .vue,.js

# Всё через Makefile
make check    # lint + test
make test     # только тесты
make lint     # только линтеры
```

## API — быстрые проверки (curl)

```bash
# Health
curl http://localhost:8091/api/health

# Список заведений
curl "http://localhost:8091/api/places"

# Фильтрация: жемчужины в Казани, сортировка по рейтингу
curl "http://localhost:8091/api/places?city=Казань&is_gem=true&sort=rating"

# Регистрация
curl -X POST http://localhost:8091/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# Логин (получить токен)
curl -X POST http://localhost:8091/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alina","password":"the_best"}'

# Профиль (с токеном)
curl http://localhost:8091/api/auth/me \
  -H "Authorization: Bearer <TOKEN>"

# Отзывы пользователя
curl "http://localhost:8091/api/users/1/reviews"

# Геосаджест
curl "http://localhost:8091/api/suggest?text=кофейня&ll=55.79,49.12"
```

## Доступ

| Откуда | URL |
|--------|-----|
| Браузер (localhost) | http://localhost:8091 |
| Телефон / другое устройство (LAN) | http://192.168.1.192:8091 |
| Vite dev-сервер напрямую | http://localhost:5173 |
| API напрямую (без nginx) | http://localhost:8086/api/... |

## Продакшн

```bash
# Создать .env
cp .env.example .env
# Отредактировать .env — задать DB_PASSWORD и JWT_SECRET

# Запустить продакшн-конфигурацию
docker compose -f docker-compose.prod.yml up --build -d
```

## Пользователи (сид-данные)

| Username | Password | ID |
|----------|----------|----|
| alina    | the_best | 1  |
| lamark   | lamark   | 2  |
