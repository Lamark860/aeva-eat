-- 010: Video messages support for reviews

ALTER TABLE reviews ADD COLUMN IF NOT EXISTS video_url TEXT;
