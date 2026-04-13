-- 004_wishlist.up.sql — planned places (wishlist)

CREATE TABLE IF NOT EXISTS wishlists (
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    place_id   INT NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, place_id)
);

CREATE INDEX IF NOT EXISTS idx_wishlists_user ON wishlists(user_id);
