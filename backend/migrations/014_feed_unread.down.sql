-- 014_feed_unread.down.sql
ALTER TABLE users DROP COLUMN IF EXISTS last_seen_feed_at;
