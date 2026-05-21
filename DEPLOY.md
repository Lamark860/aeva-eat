# Деплой AEVA Eat на VPS

Прод-стек: `postgres + backend (Go) + frontend (Vite build) + nginx`, всё через
`docker-compose.prod.yml`. Backend применяет миграции автоматически при старте
(идемпотентно — `IF NOT EXISTS` / `ON CONFLICT`).

Текущий режим доступа — по IP и порту: `http://<IP>:8091`. Перевод на домен +
HTTPS (Traefik) — см. в конце.

---

## 0. Перенос реальных данных (важно)

Реальные данные (юзеры alina/lamark, места, отзывы) живут **в БД, а не в репозитории**.
Демо-seed `005_seed_data.up.sql` (alice/bob/charlie) — фейковый, накатывается миграцией.
Чтобы на проде были настоящие данные, их нужно перенести дампом.

### Снять дамп с исходной машины (где крутится рабочая БД)

```bash
mkdir -p ~/aeva-eat-deploy
STAMP=$(date +%Y%m%d)

# БД (схема + данные)
docker exec aeva-postgres pg_dump -U aeva -d aeva_eat --no-owner --no-privileges \
  | gzip > ~/aeva-eat-deploy/aeva_eat_${STAMP}.sql.gz

# Загруженные фото/аватары (volume uploads_data)
docker run --rm -v aeva-eat_uploads_data:/u -v ~/aeva-eat-deploy:/out alpine \
  sh -c "cd /u && tar czf /out/uploads_${STAMP}.tar.gz ."
```

> Дампы содержат персональные данные и хэши паролей — **в git их не коммитим**.

---

## 1. Залить артефакты на сервер

```bash
scp ~/aeva-eat-deploy/aeva_eat_<STAMP>.sql.gz \
    ~/aeva-eat-deploy/uploads_<STAMP>.tar.gz \
    maxim@<IP>:/opt/projects/aeva-eat/
```

## 2. Код + переменные окружения

```bash
cd /opt/projects/aeva-eat
git clone git@github.com:Lamark860/aeva-eat.git .
cp .env.example .env
nano .env
```

Заполнить `.env`:

| Переменная | Значение |
|---|---|
| `DB_PASSWORD` | надёжный пароль БД |
| `JWT_SECRET` | `openssl rand -base64 48` |
| `APP_PORT` | `8091` |
| `GEOSUGGEST_KEY` | ключ Яндекс Геосаджеста |
| `VITE_YANDEX_MAPS_KEY` | ключ Яндекс JS API Карт |

## 3. Поднять только postgres и восстановить БД ДО старта backend

Порядок критичен: восстановление в свежую БД должно произойти раньше, чем
backend накатит миграции (иначе `CREATE TABLE` из дампа упрётся в уже созданную схему).

```bash
docker compose -f docker-compose.prod.yml up -d postgres
sleep 5
gzcat aeva_eat_<STAMP>.sql.gz | docker exec -i aeva-postgres-prod psql -U aeva -d aeva_eat
```

## 4. Восстановить фото в volume

```bash
docker run --rm -v aeva-eat_uploads_data:/u -v $(pwd):/in alpine \
  sh -c "cd /u && tar xzf /in/uploads_<STAMP>.tar.gz"
```

## 5. Поднять остальной стек

Миграции backend поверх восстановленных данных просто no-op.

```bash
docker compose -f docker-compose.prod.yml up --build -d
docker compose -f docker-compose.prod.yml ps
curl -s http://localhost:8091/api/health
```

Приложение: `http://<IP>:8091`.

---

## Заметки

- **Суперюзеры** alina/lamark приезжают в дампе с ролью `superuser` — ручной SQL не нужен.
  Если стартуешь с пустой БД (без дампа), назначить вручную:
  ```bash
  docker exec -it aeva-postgres-prod psql -U aeva -d aeva_eat \
    -c "UPDATE users SET role='superuser' WHERE username='<твой_логин>';"
  ```
- **Сменить слабые пароли.** У демо-аккаунтов (charlie, seed_*) и у `lamark` пароль
  из разряда демо (`demo12345`/`password123`). На публичном проде смени пароль
  суперюзера через профиль сразу после первого входа.
- **Яндекс-ключи** привязываются к домену/хосту в кабинете. По голому IP карта может
  не отрисоваться (referer); геосаджест (серверный) работает в любом случае.
- **Бэкап БД:**
  ```bash
  docker exec aeva-postgres-prod pg_dump -U aeva -d aeva_eat --no-owner \
    | gzip > backup_$(date +%F).sql.gz
  ```

---

## Перевод на домен + HTTPS (Traefik)

Когда DNS домена будет указывать на сервер, переключение с порта на Traefik —
правка `docker-compose.prod.yml` сервиса `nginx`:

1. Убрать проброс порта `ports: ["${APP_PORT}:80"]`.
2. Подключить к внешней сети `proxy` (создаётся Traefik'ом).
3. Добавить лейблы:
   ```yaml
   labels:
     - "traefik.enable=true"
     - "traefik.http.routers.aeva.rule=Host(`<домен>`)"
     - "traefik.http.routers.aeva.entrypoints=websecure"
     - "traefik.http.routers.aeva.tls.certresolver=letsencrypt"
     - "traefik.http.services.aeva.loadbalancer.server.port=80"
   ```

postgres/backend/frontend остаются во внутренней сети и наружу не публикуются.
