-- 010 down: revert video messages support
ALTER TABLE reviews DROP COLUMN IF EXISTS video_url;
