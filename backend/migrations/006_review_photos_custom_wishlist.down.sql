-- 006_review_photos_custom_wishlist.down.sql

ALTER TABLE reviews DROP COLUMN IF EXISTS image_url;
DROP TABLE IF EXISTS wishlist_custom;
