-- 013_notes.up.sql
--
-- Записки от руки на доске (backend.md §notes). Записка — отдельный класс
-- артефакта в ленте: текст без визита, опционально привязан к месту/городу,
-- с цветом бумаги/тейпа и флагом "зачёркнуто" (для wishlist-перечёркивания
-- в будущем).

CREATE TABLE IF NOT EXISTS notes (
    id          BIGSERIAL PRIMARY KEY,
    author_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text        TEXT NOT NULL,
    place_id    BIGINT REFERENCES places(id) ON DELETE SET NULL,
    city        VARCHAR(255),
    paper_color VARCHAR(20),
    tape_color  VARCHAR(20),
    is_struck   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notes_author  ON notes(author_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notes_created ON notes(created_at DESC);

-- feed_events VIEW — единая хронология ленты (backend.md §feed_events).
-- Сейчас покрывает только review_added и note_added. Остальные типы
-- (wishlist_*, gem_found, kruzhok_added) добавим, когда появятся таблицы /
-- триггеры. VIEW идемпотентна — пересоздаём CREATE OR REPLACE.

CREATE OR REPLACE VIEW feed_events AS
    SELECT 'review_added'::text                                  AS kind,
           rv.id                                                  AS event_id,
           rv.created_at                                          AS occurred_at,
           rv.place_id                                            AS place_id,
           (SELECT ra.user_id FROM review_authors ra
              WHERE ra.review_id = rv.id ORDER BY ra.user_id LIMIT 1) AS author_id,
           rv.id                                                  AS review_id,
           NULL::BIGINT                                           AS note_id
      FROM reviews rv
    UNION ALL
    SELECT 'note_added'::text,
           n.id,
           n.created_at,
           n.place_id,
           n.author_id,
           NULL::BIGINT,
           n.id
      FROM notes n;
