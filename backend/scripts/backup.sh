#!/usr/bin/env bash
# Бэкап прод-данных AEVA Eat: дамп БД + архив загруженных фото/видео/аватаров.
# Единственная невосстановимая ценность приложения — пользовательский контент
# (отзывы, заметки, фото, видео-кружки). Запускать по cron (см. DEPLOY.md).
#
# Использование:
#   ./backend/scripts/backup.sh [DEST_DIR]
# DEST_DIR — куда складывать (по умолчанию ./backups рядом с .env прод-стека).
#
# Cron (ежедневно в 04:30, хранить 14 копий — ротация ниже):
#   30 4 * * * cd /opt/projects/aeva-eat && ./backend/scripts/backup.sh >> /var/log/aeva-backup.log 2>&1
set -euo pipefail

DEST="${1:-./backups}"
KEEP="${BACKUP_KEEP:-14}"            # сколько последних копий хранить
PG_CONTAINER="${PG_CONTAINER:-aeva-postgres-prod}"
UPLOADS_VOLUME="${UPLOADS_VOLUME:-aeva-eat_uploads_data}"

# Имя БД/пользователя берём из .env прод-стека (рядом со скриптом запуска).
DB_USER="${DB_USER:-$(grep -E '^DB_USER=' .env 2>/dev/null | cut -d= -f2- | tr -d '[:space:]')}"
DB_NAME="${DB_NAME:-$(grep -E '^DB_NAME=' .env 2>/dev/null | cut -d= -f2- | tr -d '[:space:]')}"
DB_USER="${DB_USER:-aeva}"
DB_NAME="${DB_NAME:-aeva_eat}"

STAMP="$(date +%Y%m%d_%H%M%S)"
mkdir -p "$DEST"

echo "==> [$(date)] Дамп БД ${DB_NAME} из ${PG_CONTAINER}"
docker exec "$PG_CONTAINER" pg_dump -U "$DB_USER" -d "$DB_NAME" --no-owner --no-privileges \
  | gzip > "${DEST}/aeva_eat_${STAMP}.sql.gz"

echo "==> [$(date)] Архив uploads (${UPLOADS_VOLUME})"
docker run --rm -v "${UPLOADS_VOLUME}":/u -v "$(cd "$DEST" && pwd)":/out alpine \
  sh -c "cd /u && tar czf /out/uploads_${STAMP}.tar.gz ."

echo "==> Ротация: оставляю последние ${KEEP} копий каждого вида"
ls -1t "${DEST}"/aeva_eat_*.sql.gz 2>/dev/null | tail -n +$((KEEP + 1)) | xargs -r rm -f
ls -1t "${DEST}"/uploads_*.tar.gz  2>/dev/null | tail -n +$((KEEP + 1)) | xargs -r rm -f

echo "==> Готово: ${DEST}/aeva_eat_${STAMP}.sql.gz + uploads_${STAMP}.tar.gz"
echo "    ВАЖНО: периодически копируйте бэкапы off-site (S3/другой хост) —"
echo "    локальная копия на том же диске не спасёт при отказе диска VPS."
