-- 003_rating_float.down.sql

ALTER TABLE reviews ALTER COLUMN food_rating TYPE INT USING food_rating::INT;
ALTER TABLE reviews ALTER COLUMN service_rating TYPE INT USING service_rating::INT;
ALTER TABLE reviews ALTER COLUMN vibe_rating TYPE INT USING vibe_rating::INT;

ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_food_rating_check;
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_service_rating_check;
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_vibe_rating_check;

ALTER TABLE reviews ADD CONSTRAINT reviews_food_rating_check CHECK (food_rating BETWEEN 1 AND 10);
ALTER TABLE reviews ADD CONSTRAINT reviews_service_rating_check CHECK (service_rating BETWEEN 1 AND 10);
ALTER TABLE reviews ADD CONSTRAINT reviews_vibe_rating_check CHECK (vibe_rating BETWEEN 1 AND 10);
