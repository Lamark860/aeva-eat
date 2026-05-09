-- 014_feed_unread.up.sql
--
-- Точка-индикатор "новостей" на табе Доска (NEXT.md §C2). Храним момент
-- последнего открытия Доски пользователем; unread-count считается на лету
-- как кол-во feed_events с occurred_at > last_seen_feed_at, исключая
-- собственные события (свои визиты не показываем как "новое").

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS last_seen_feed_at TIMESTAMPTZ;

-- Дефолт для существующих юзеров — текущий момент, чтобы они не увидели
-- большой счётчик "сотни новостей" при первой загрузке после деплоя.
UPDATE users SET last_seen_feed_at = now()
 WHERE last_seen_feed_at IS NULL;
