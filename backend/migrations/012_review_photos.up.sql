-- 012_review_photos.up.sql
--
-- Несколько фото на отзыв. Дизайнерская стопка полароидов в ленте,
-- горизонтальный ряд на странице места (DESIGN-DECISIONS §L3).
--
-- reviews.image_url НЕ удаляем — остаётся как backwards-compat первое
-- фото и используется в COALESCE-fallback'е cover'а места.

CREATE TABLE IF NOT EXISTS review_photos (
    id         BIGSERIAL PRIMARY KEY,
    review_id  BIGINT NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    url        VARCHAR(500) NOT NULL,
    position   INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_review_photos_review
    ON review_photos(review_id, position);

-- Backfill: существующие отзывы с image_url получают первое фото в стопке.
INSERT INTO review_photos (review_id, url, position, created_at)
SELECT id, image_url, 0, created_at
  FROM reviews
 WHERE image_url IS NOT NULL AND image_url <> ''
   AND NOT EXISTS (
       SELECT 1 FROM review_photos rp WHERE rp.review_id = reviews.id
   );
