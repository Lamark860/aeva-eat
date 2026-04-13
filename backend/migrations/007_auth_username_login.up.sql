-- 007: Login by username, drop email requirement, add display_name

-- Add display_name (optional, non-unique)
ALTER TABLE users ADD COLUMN IF NOT EXISTS display_name VARCHAR(100);

-- Make email optional (drop NOT NULL, keep column for legacy data)
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;
ALTER TABLE users ALTER COLUMN email SET DEFAULT NULL;

-- Update existing users: set display_name = username where null
UPDATE users SET display_name = username WHERE display_name IS NULL;
