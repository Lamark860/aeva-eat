-- Откат soft-delete. ВНИМАНИЕ: вернёт прежнее поведение (hard-delete) и
-- не-партиальную уникальность.
DROP INDEX IF EXISTS idx_places_active;
DROP INDEX IF EXISTS idx_places_identity;
CREATE UNIQUE INDEX IF NOT EXISTS idx_places_identity
    ON places (LOWER(name), LOWER(COALESCE(address, '')), LOWER(COALESCE(city, '')));
ALTER TABLE places DROP COLUMN IF EXISTS deleted_at;
