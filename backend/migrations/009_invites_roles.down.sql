-- 009 down: revert invite-only system

DROP TABLE IF EXISTS invites;
ALTER TABLE users DROP COLUMN IF EXISTS role;
