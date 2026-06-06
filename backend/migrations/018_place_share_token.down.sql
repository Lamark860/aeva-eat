DROP INDEX IF EXISTS idx_places_share_token;
ALTER TABLE places DROP COLUMN IF EXISTS share_token;
