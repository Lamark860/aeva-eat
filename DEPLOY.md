# Деплой AEVA Eat на VPS

Прод-стек: `postgres + backend (Go) + frontend (Vite build) + nginx`, всё через
`docker-compose.prod.yml`. Backend применяет миграции схемы автоматически при старте
(идемпотентно — `IF NOT EXISTS` / `ON CONFLICT`). Демо-данные (`005_seed_data`)
**больше не сидятся автоматически** — только реальные данные из дампа (см. §0).

Текущий режим доступа — по IP и порту: `http://<IP>:8091`. Перевод на домен +
HTTPS (Traefik) — см. в конце.

> ⚠️ `APP_PORT=8091` обязателен в `.env`: на сервере 80/443 уже заняты Traefik'ом,
> и без своего порта nginx приложения упадёт с `port is already allocated`.

---

## 0. Перенос реальных данных (важно)

Реальные данные (юзеры alina/lamark, места, отзывы) живут **в БД, а не в репозитории**.
Чтобы на проде были настоящие данные, их нужно перенести дампом.

> Демо-данные (`005_seed_data` — alice/bob/charlie, и `scripts/seed_demo.sh` — seed_*)
> в прод НЕ нужны. `005` убран из авто-миграций; если случайно зальёшь демо —
> чистка: `DELETE FROM places WHERE created_by IN (SELECT id FROM users WHERE username LIKE 'seed\_%'); DELETE FROM users WHERE username LIKE 'seed\_%';`

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
| `APP_ENV` | `production` (включает fail-fast: пустой/дефолтный `JWT_SECRET` → backend не стартует) |
| `DB_PASSWORD` | надёжный пароль БД |
| `JWT_SECRET` | `openssl rand -base64 48` (в `production` обязан быть не-дефолтным) |
| `APP_PORT` | `8091` |
| `GEOSUGGEST_KEY` | ключ Яндекс Геосаджеста |
| `VITE_YANDEX_MAPS_KEY` | ключ Яндекс JS API Карт |

## 3. Поднять только postgres и восстановить БД ДО старта backend

Порядок критичен: восстановление в свежую БД должно произойти раньше, чем
backend накатит миграции (иначе `CREATE TABLE` из дампа упрётся в уже созданную схему).

```bash
docker compose -f docker-compose.prod.yml up -d postgres
sleep 5
# на Linux — zcat (gzcat это macOS). Альтернатива: gunzip -c
zcat aeva_eat_<STAMP>.sql.gz | docker exec -i aeva-postgres-prod psql -U aeva -d aeva_eat
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
- **Сменить слабые пароли.** У `lamark` пароль из разряда демо (`demo12345`).
  На публичном проде смени пароль суперюзера через профиль сразу после первого входа.
- **DB_PASSWORD со спецсимволами** теперь безопасен (DSN URL-экранируется в коде).
  Но если разворачиваешь старый образ — избегай `( ) @ / ?` в пароле, иначе
  backend упадёт с `invalid port after host` (502). Сменить пароль роли без потери данных:
  ```bash
  docker exec aeva-postgres-prod psql -U aeva -d aeva_eat -c "ALTER USER aeva WITH PASSWORD '<новый>';"
  ```
  (затем тот же пароль в `.env` → `up -d backend`)
- **Яндекс-ключи** привязываются к домену/хосту в кабинете. По голому IP карта может
  не отрисоваться (referer); геосаджест (серверный) работает в любом случае.
- **Бэкап БД + фото (по расписанию).** Пользовательский контент (отзывы, заметки,
  фото, видео) — единственная невосстановимая ценность; держать его без бэкапа на
  одном диске VPS = риск безвозвратной потери. В репозитории есть готовый скрипт
  `backend/scripts/backup.sh` (дамп БД + архив volume `uploads_data` + ротация):
  ```bash
  # разовый запуск из каталога прод-стека
  cd /opt/projects/aeva-eat && ./backend/scripts/backup.sh

  # cron: ежедневно в 04:30, хранить 14 копий
  crontab -e
  30 4 * * * cd /opt/projects/aeva-eat && ./backend/scripts/backup.sh >> /var/log/aeva-backup.log 2>&1
  ```
  Раз в период копируйте дампы **off-site** (S3/другой хост) и хотя бы раз
  проверьте восстановление из бэкапа.

  Разовый дамп только БД (без скрипта):
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
