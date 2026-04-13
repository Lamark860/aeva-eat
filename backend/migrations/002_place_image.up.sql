-- 002_place_image.up.sql
ALTER TABLE places ADD COLUMN IF NOT EXISTS image_url TEXT;
