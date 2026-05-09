-- 015_wishlist_struck.up.sql
--
-- backend.md §Wishlist: общий wishlist круга + флаг "зачёркнуто" для
-- записок-планов на доске. Когда запись wishlist реализована визитом —
-- struck=true (страйк-логика делается в сервисе при создании review,
-- а не БД-триггером, чтобы видеть событие в логах).

ALTER TABLE wishlists
    ADD COLUMN IF NOT EXISTS is_struck   BOOLEAN     NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS struck_at   TIMESTAMPTZ;

-- Backfill: для существующих wishlist-записей, по которым автор уже
-- оставил review, помечаем как зачёркнутые с моментом review.created_at.
UPDATE wishlists w
   SET is_struck = TRUE,
       struck_at = sub.struck_at
  FROM (
        SELECT ra.user_id, rv.place_id, MIN(rv.created_at) AS struck_at
          FROM reviews rv
          JOIN review_authors ra ON ra.review_id = rv.id
         GROUP BY ra.user_id, rv.place_id
       ) sub
 WHERE w.user_id  = sub.user_id
   AND w.place_id = sub.place_id
   AND w.is_struck = FALSE;

CREATE INDEX IF NOT EXISTS idx_wishlists_place ON wishlists(place_id);
