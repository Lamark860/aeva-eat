-- 003_rating_float.up.sql — ratings from INT(1-10) to NUMERIC(0-10, step 0.1)

ALTER TABLE reviews ALTER COLUMN food_rating TYPE NUMERIC(3,1) USING food_rating::NUMERIC(3,1);
ALTER TABLE reviews ALTER COLUMN service_rating TYPE NUMERIC(3,1) USING service_rating::NUMERIC(3,1);
ALTER TABLE reviews ALTER COLUMN vibe_rating TYPE NUMERIC(3,1) USING vibe_rating::NUMERIC(3,1);

ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_food_rating_check;
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_service_rating_check;
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_vibe_rating_check;

ALTER TABLE reviews ADD CONSTRAINT reviews_food_rating_check CHECK (food_rating >= 0 AND food_rating <= 10);
ALTER TABLE reviews ADD CONSTRAINT reviews_service_rating_check CHECK (service_rating >= 0 AND service_rating <= 10);
ALTER TABLE reviews ADD CONSTRAINT reviews_vibe_rating_check CHECK (vibe_rating >= 0 AND vibe_rating <= 10);
