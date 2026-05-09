-- 015_wishlist_struck.down.sql
DROP INDEX IF EXISTS idx_wishlists_place;
ALTER TABLE wishlists
    DROP COLUMN IF EXISTS struck_at,
    DROP COLUMN IF EXISTS is_struck;
