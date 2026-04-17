-- 009: Invite-only system + user roles

-- Add role column to users (superuser can manage invites and approve users)
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'user';

-- Mark first user (id=1) as superuser, or seed users alina/lamark
UPDATE users SET role = 'superuser' WHERE username IN ('alina', 'lamark');

-- Invite codes table
CREATE TABLE IF NOT EXISTS invites (
    id          SERIAL PRIMARY KEY,
    code        VARCHAR(32) NOT NULL UNIQUE,
    created_by  INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    used_by     INT REFERENCES users(id) ON DELETE SET NULL,
    used_at     TIMESTAMPTZ,
    expires_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_invites_code ON invites(code);
