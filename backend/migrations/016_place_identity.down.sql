-- Откат к старому ключу уникальности (name + city). ВНИМАНИЕ: вернёт баг с
-- невозможностью завести второй филиал сети в одном городе.
DROP INDEX IF EXISTS idx_places_identity;
CREATE UNIQUE INDEX IF NOT EXISTS idx_places_name_city
    ON places (LOWER(name), LOWER(COALESCE(city, '')));
