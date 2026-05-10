package repository

import (
	"database/sql"

	"github.com/aeva-eat/backend/internal/model"
	"github.com/lib/pq"
)

type NoteRepo struct {
	db *sql.DB
}

func NewNoteRepo(db *sql.DB) *NoteRepo {
	return &NoteRepo{db: db}
}

const noteSelectColumns = `
	n.id, n.author_id, n.text, n.place_id, n.city,
	n.paper_color, n.tape_color, n.is_struck,
	n.created_at, n.updated_at,
	u.username, u.avatar_url, u.role, u.created_at, u.updated_at,
	p.name AS place_name
`

const noteJoins = `
	FROM notes n
	JOIN users u ON u.id = n.author_id
	LEFT JOIN places p ON p.id = n.place_id
`

func scanNote(row interface {
	Scan(dest ...any) error
}) (*model.Note, error) {
	n := &model.Note{Author: &model.User{}}
	err := row.Scan(
		&n.ID, &n.AuthorID, &n.Text, &n.PlaceID, &n.City,
		&n.PaperColor, &n.TapeColor, &n.IsStruck,
		&n.CreatedAt, &n.UpdatedAt,
		&n.Author.Username, &n.Author.AvatarURL, &n.Author.Role,
		&n.Author.CreatedAt, &n.Author.UpdatedAt,
		&n.PlaceName,
	)
	if err != nil {
		return nil, err
	}
	n.Author.ID = n.AuthorID
	return n, nil
}

// List — все записки круга, свежие сверху.
func (r *NoteRepo) List() ([]model.Note, error) {
	rows, err := r.db.Query("SELECT " + noteSelectColumns + noteJoins +
		" ORDER BY n.created_at DESC, n.id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []model.Note{}
	for rows.Next() {
		n, err := scanNote(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, *n)
	}
	return notes, rows.Err()
}

// ListByAuthor — записки конкретного пользователя (для профиля).
func (r *NoteRepo) ListByAuthor(authorID int) ([]model.Note, error) {
	rows, err := r.db.Query("SELECT "+noteSelectColumns+noteJoins+
		" WHERE n.author_id = $1 ORDER BY n.created_at DESC, n.id DESC", authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []model.Note{}
	for rows.Next() {
		n, err := scanNote(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, *n)
	}
	return notes, rows.Err()
}

func (r *NoteRepo) GetByID(id int) (*model.Note, error) {
	row := r.db.QueryRow("SELECT "+noteSelectColumns+noteJoins+" WHERE n.id = $1", id)
	return scanNote(row)
}

func (r *NoteRepo) Create(n *model.Note) (*model.Note, error) {
	err := r.db.QueryRow(`
		INSERT INTO notes (author_id, text, place_id, city, paper_color, tape_color)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, n.AuthorID, n.Text, n.PlaceID, n.City, n.PaperColor, n.TapeColor).
		Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return r.GetByID(n.ID)
}

func (r *NoteRepo) Update(n *model.Note) (*model.Note, error) {
	_, err := r.db.Exec(`
		UPDATE notes SET text=$1, place_id=$2, city=$3,
			paper_color=$4, tape_color=$5, updated_at=now()
		WHERE id=$6
	`, n.Text, n.PlaceID, n.City, n.PaperColor, n.TapeColor, n.ID)
	if err != nil {
		return nil, err
	}
	return r.GetByID(n.ID)
}

func (r *NoteRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM notes WHERE id = $1", id)
	return err
}

func (r *NoteRepo) SetStruck(id int, struck bool) (*model.Note, error) {
	_, err := r.db.Exec(
		"UPDATE notes SET is_struck=$1, updated_at=now() WHERE id=$2",
		struck, id,
	)
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// IsAuthor — проверка прав на изменение записки.
func (r *NoteRepo) IsAuthor(noteID, userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM notes WHERE id=$1 AND author_id=$2)",
		noteID, userID,
	).Scan(&exists)
	return exists, err
}

// FeedEventsRepo — отдельный репо для VIEW feed_events.
type FeedEventsRepo struct {
	db *sql.DB
}

func NewFeedEventsRepo(db *sql.DB) *FeedEventsRepo {
	return &FeedEventsRepo{db: db}
}

// Weeks — список недель с count событий и gem_count (backend.md §Лента).
// Обрезаем до 26 недель (полгода) — глубже доска визуально не нужна.
func (r *FeedEventsRepo) Weeks() ([]model.FeedWeek, error) {
	rows, err := r.db.Query(`
		WITH events AS (
			SELECT
				date_trunc('week', e.occurred_at)::date AS week_start,
				e.kind,
				e.event_id,
				e.review_id
			FROM feed_events e
		)
		SELECT
			to_char(week_start, 'IYYY-"W"IW') AS key,
			week_start,
			COUNT(*)                          AS count,
			COUNT(*) FILTER (
				WHERE kind = 'review_added'
				  AND review_id IN (SELECT id FROM reviews WHERE is_gem = true)
			)                                 AS gem_count
		FROM events
		GROUP BY week_start
		ORDER BY week_start DESC
		LIMIT 26
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	weeks := []model.FeedWeek{}
	for rows.Next() {
		var w model.FeedWeek
		if err := rows.Scan(&w.Key, &w.WeekStart, &w.Count, &w.GemCount); err != nil {
			return nil, err
		}
		weeks = append(weeks, w)
	}
	return weeks, rows.Err()
}

// UnreadCount — число событий после last_seen_feed_at пользователя,
// исключая его собственные (свои визиты не считаем как "новое").
func (r *FeedEventsRepo) UnreadCount(userID int) (int, error) {
	var n int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM feed_events e
		WHERE e.occurred_at > COALESCE(
			(SELECT last_seen_feed_at FROM users WHERE id = $1),
			'epoch'::timestamptz
		)
		AND (e.author_id IS NULL OR e.author_id <> $1)
	`, userID).Scan(&n)
	return n, err
}

// MarkSeen — обновить last_seen_feed_at пользователя на now().
func (r *FeedEventsRepo) MarkSeen(userID int) error {
	_, err := r.db.Exec(
		`UPDATE users SET last_seen_feed_at = now() WHERE id = $1`,
		userID,
	)
	return err
}

// List возвращает события ленты, свежие сверху. Пагинация по timestamp
// (cursor) — опциональна; пока без неё, лимитируем на server-side.
// Q4 — каждое событие тащит attendees: для review_added берём всех co-authors,
// для note_added это единственный author. Фронт группирует по
// (place_id, date(occurred_at), attendees) — это и закрывает «раз в месяц
// в одно и то же место» как отдельные артефакты, а не одну плитку.
func (r *FeedEventsRepo) List(limit int) ([]model.FeedEvent, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := r.db.Query(`
		SELECT e.kind, e.event_id, e.occurred_at, e.place_id, e.author_id, e.review_id, e.note_id,
			COALESCE(
				(SELECT array_agg(ra.user_id ORDER BY ra.user_id)
					FROM review_authors ra WHERE ra.review_id = e.review_id),
				CASE WHEN e.author_id IS NOT NULL
					THEN ARRAY[e.author_id]::BIGINT[]
					ELSE ARRAY[]::BIGINT[]
				END
			) AS attendees,
			(SELECT NULLIF(rv.video_url, '') FROM reviews rv WHERE rv.id = e.review_id) AS video_url
		FROM feed_events e
		ORDER BY e.occurred_at DESC, e.event_id DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []model.FeedEvent{}
	for rows.Next() {
		var e model.FeedEvent
		var attendees pq.Int64Array
		if err := rows.Scan(&e.Kind, &e.EventID, &e.OccurredAt,
			&e.PlaceID, &e.AuthorID, &e.ReviewID, &e.NoteID, &attendees,
			&e.VideoURL); err != nil {
			return nil, err
		}
		e.Attendees = make([]int, len(attendees))
		for i, v := range attendees {
			e.Attendees[i] = int(v)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
