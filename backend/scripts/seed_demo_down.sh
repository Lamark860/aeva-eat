#!/usr/bin/env bash
# seed_demo_down.sh — снимает всё что залил seed_demo.sh.
#
# Идемпотентно: безопасно запускать многократно.
#
# Удаляет:
#   1) places.created_by IN seed_users — каскадом подметает reviews,
#      review_authors, review_photos, place_categories, wishlists.
#   2) users WHERE username LIKE 'seed_%' — каскадом подметает notes,
#      wishlist_custom, wishlists, invites.
#   3) /app/uploads/seed_*.jpg.
#
# Реальные данные (lamark/alina/charlie и их объекты) НЕ затрагиваются.

set -euo pipefail

PG_CONTAINER=${PG_CONTAINER:-aeva-postgres}
BE_CONTAINER=${BE_CONTAINER:-aeva-backend}
DB=${DB:-aeva_eat}
DB_USER=${DB_USER:-aeva}

echo "→ removing seed objects from DB"
docker exec -i "$PG_CONTAINER" psql -U "$DB_USER" -d "$DB" -v ON_ERROR_STOP=1 <<'SQL'
BEGIN;

DELETE FROM places
 WHERE created_by IN (SELECT id FROM users WHERE username LIKE 'seed_%');

DELETE FROM users WHERE username LIKE 'seed_%';

COMMIT;

\echo
\echo '=== after-unseed counts ==='
SELECT
  (SELECT COUNT(*) FROM users WHERE username LIKE 'seed_%')   AS seed_users_left,
  (SELECT COUNT(*) FROM users)                                 AS users_total,
  (SELECT COUNT(*) FROM places)                                AS places_total,
  (SELECT COUNT(*) FROM reviews)                               AS reviews_total;
SQL

echo "→ removing /app/uploads/seed_*.jpg"
docker exec "$BE_CONTAINER" sh -c 'rm -f /app/uploads/seed_*.jpg && echo "  removed: $(ls /app/uploads/seed_* 2>/dev/null | wc -l) files left"'

echo "✓ unseed done"
