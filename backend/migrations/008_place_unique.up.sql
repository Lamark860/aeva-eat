-- Unique constraint: no duplicate place with same name+city
CREATE UNIQUE INDEX IF NOT EXISTS idx_places_name_city
    ON places (LOWER(name), LOWER(COALESCE(city, '')));
