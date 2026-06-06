-- Неугадываемый share-токен для публичных страниц /p/<token>. Раньше шарилось
-- по инкрементному id (/p/123), что позволяло перебором выгрузить весь каталог
-- названий/городов/обложек. UUID убирает энумерацию. gen_random_uuid() —
-- встроенная в Postgres 13+ (расширение не нужно).
ALTER TABLE places ADD COLUMN IF NOT EXISTS share_token UUID;
UPDATE places SET share_token = gen_random_uuid() WHERE share_token IS NULL;
ALTER TABLE places ALTER COLUMN share_token SET DEFAULT gen_random_uuid();
ALTER TABLE places ALTER COLUMN share_token SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_places_share_token ON places(share_token);
