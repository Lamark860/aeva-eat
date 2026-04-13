-- 006_review_photos_custom_wishlist.up.sql

-- Review photos (image attached to a review)
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS image_url TEXT;

-- Custom wishlist (free-text places user wants to visit, not in system)
CREATE TABLE IF NOT EXISTS wishlist_custom (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    note       TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_wishlist_custom_user ON wishlist_custom(user_id);
