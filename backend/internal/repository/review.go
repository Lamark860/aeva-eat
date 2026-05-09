package repository

import (
	"database/sql"
	"fmt"

	"github.com/aeva-eat/backend/internal/model"
)

type ReviewRepo struct {
	db *sql.DB
}

func NewReviewRepo(db *sql.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) ListByPlace(placeID int) ([]model.Review, error) {
	rows, err := r.db.Query(`
		SELECT rv.id, rv.place_id, rv.food_rating, rv.service_rating, rv.vibe_rating,
			rv.is_gem, rv.comment, rv.image_url, rv.video_url, rv.visited_at, rv.created_at, rv.updated_at
		FROM reviews rv
		WHERE rv.place_id = $1
		ORDER BY rv.created_at DESC
	`, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var rv model.Review
		if err := rows.Scan(&rv.ID, &rv.PlaceID, &rv.FoodRating, &rv.ServiceRating,
			&rv.VibeRating, &rv.IsGem, &rv.Comment, &rv.ImageURL, &rv.VideoURL, &rv.VisitedAt,
			&rv.CreatedAt, &rv.UpdatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range reviews {
		authors, err := r.getAuthors(reviews[i].ID)
		if err != nil {
			return nil, err
		}
		reviews[i].Authors = authors

		photos, err := r.ListPhotos(reviews[i].ID)
		if err != nil {
			return nil, err
		}
		reviews[i].Photos = photos
	}

	return reviews, nil
}

func (r *ReviewRepo) ListByUser(userID int) ([]model.Review, error) {
	rows, err := r.db.Query(`
		SELECT rv.id, rv.place_id, COALESCE(p.name, '') AS place_name,
			rv.food_rating, rv.service_rating, rv.vibe_rating,
			rv.is_gem, rv.comment, rv.image_url, rv.video_url, rv.visited_at, rv.created_at, rv.updated_at
		FROM reviews rv
		JOIN review_authors ra ON ra.review_id = rv.id
		LEFT JOIN places p ON p.id = rv.place_id
		WHERE ra.user_id = $1
		ORDER BY rv.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var rv model.Review
		if err := rows.Scan(&rv.ID, &rv.PlaceID, &rv.PlaceName, &rv.FoodRating, &rv.ServiceRating,
			&rv.VibeRating, &rv.IsGem, &rv.Comment, &rv.ImageURL, &rv.VideoURL, &rv.VisitedAt,
			&rv.CreatedAt, &rv.UpdatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range reviews {
		authors, err := r.getAuthors(reviews[i].ID)
		if err != nil {
			return nil, err
		}
		reviews[i].Authors = authors

		photos, err := r.ListPhotos(reviews[i].ID)
		if err != nil {
			return nil, err
		}
		reviews[i].Photos = photos
	}

	return reviews, nil
}

func (r *ReviewRepo) GetByID(id int) (*model.Review, error) {
	rv := &model.Review{}
	err := r.db.QueryRow(`
		SELECT id, place_id, food_rating, service_rating, vibe_rating,
			is_gem, comment, image_url, video_url, visited_at, created_at, updated_at
		FROM reviews WHERE id = $1
	`, id).Scan(&rv.ID, &rv.PlaceID, &rv.FoodRating, &rv.ServiceRating,
		&rv.VibeRating, &rv.IsGem, &rv.Comment, &rv.ImageURL, &rv.VideoURL, &rv.VisitedAt,
		&rv.CreatedAt, &rv.UpdatedAt)
	if err != nil {
		return nil, err
	}

	authors, err := r.getAuthors(rv.ID)
	if err != nil {
		return nil, err
	}
	rv.Authors = authors

	photos, err := r.ListPhotos(rv.ID)
	if err != nil {
		return nil, err
	}
	rv.Photos = photos

	return rv, nil
}

func (r *ReviewRepo) Create(rv *model.Review, authorIDs []int) (*model.Review, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
		INSERT INTO reviews (place_id, food_rating, service_rating, vibe_rating, is_gem, comment, image_url, video_url, visited_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`, rv.PlaceID, rv.FoodRating, rv.ServiceRating, rv.VibeRating,
		rv.IsGem, rv.Comment, rv.ImageURL, rv.VideoURL, rv.VisitedAt,
	).Scan(&rv.ID, &rv.CreatedAt, &rv.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert review: %w", err)
	}

	for _, uid := range authorIDs {
		_, err := tx.Exec(`INSERT INTO review_authors (review_id, user_id) VALUES ($1, $2)`, rv.ID, uid)
		if err != nil {
			return nil, fmt.Errorf("insert author %d: %w", uid, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(rv.ID)
}

func (r *ReviewRepo) Update(rv *model.Review) (*model.Review, error) {
	_, err := r.db.Exec(`
		UPDATE reviews SET food_rating=$1, service_rating=$2, vibe_rating=$3,
			is_gem=$4, comment=$5, image_url=$6, video_url=$7, visited_at=$8, updated_at=now()
		WHERE id=$9
	`, rv.FoodRating, rv.ServiceRating, rv.VibeRating,
		rv.IsGem, rv.Comment, rv.ImageURL, rv.VideoURL, rv.VisitedAt, rv.ID)
	if err != nil {
		return nil, err
	}
	return r.GetByID(rv.ID)
}

func (r *ReviewRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM reviews WHERE id = $1`, id)
	return err
}

func (r *ReviewRepo) IsAuthor(reviewID, userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM review_authors WHERE review_id=$1 AND user_id=$2)`,
		reviewID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *ReviewRepo) UpdateImageURL(reviewID int, imageURL string) error {
	_, err := r.db.Exec(`UPDATE reviews SET image_url = $1, updated_at = now() WHERE id = $2`, imageURL, reviewID)
	return err
}

// ListPhotos возвращает фото отзыва, отсортированные по position (затем id).
func (r *ReviewRepo) ListPhotos(reviewID int) ([]model.ReviewPhoto, error) {
	rows, err := r.db.Query(`
		SELECT id, url, position
		FROM review_photos
		WHERE review_id = $1
		ORDER BY position ASC, id ASC
	`, reviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photos := []model.ReviewPhoto{}
	for rows.Next() {
		var p model.ReviewPhoto
		if err := rows.Scan(&p.ID, &p.URL, &p.Position); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

// CountPhotos — текущее число фото для лимита 5.
func (r *ReviewRepo) CountPhotos(reviewID int) (int, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM review_photos WHERE review_id = $1`, reviewID).Scan(&n)
	return n, err
}

// AddPhoto добавляет фото в конец стопки. Возвращает созданную запись.
// Если это первое фото и у отзыва нет image_url — синхронизирует image_url
// для backwards-compat (cover-fallback в /api/places).
func (r *ReviewRepo) AddPhoto(reviewID int, url string) (*model.ReviewPhoto, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var nextPos int
	err = tx.QueryRow(
		`SELECT COALESCE(MAX(position), -1) + 1 FROM review_photos WHERE review_id = $1`,
		reviewID,
	).Scan(&nextPos)
	if err != nil {
		return nil, err
	}

	p := &model.ReviewPhoto{URL: url, Position: nextPos}
	err = tx.QueryRow(
		`INSERT INTO review_photos (review_id, url, position) VALUES ($1, $2, $3) RETURNING id`,
		reviewID, url, nextPos,
	).Scan(&p.ID)
	if err != nil {
		return nil, err
	}

	if nextPos == 0 {
		_, err = tx.Exec(
			`UPDATE reviews SET image_url = $1, updated_at = now()
			 WHERE id = $2 AND (image_url IS NULL OR image_url = '')`,
			url, reviewID,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return p, nil
}

// GetPhoto возвращает запись по id (нужен handler'у — узнать review_id и url
// для проверки авторства и удаления файла с диска).
func (r *ReviewRepo) GetPhoto(photoID int) (*model.ReviewPhoto, int, error) {
	var p model.ReviewPhoto
	var reviewID int
	err := r.db.QueryRow(
		`SELECT id, url, position, review_id FROM review_photos WHERE id = $1`,
		photoID,
	).Scan(&p.ID, &p.URL, &p.Position, &reviewID)
	if err != nil {
		return nil, 0, err
	}
	return &p, reviewID, nil
}

// DeletePhoto удаляет одно фото. Если это было фото на позиции 0,
// синкает image_url отзыва на следующее (или NULL, если стопка пуста).
func (r *ReviewRepo) DeletePhoto(photoID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var reviewID, pos int
	err = tx.QueryRow(
		`DELETE FROM review_photos WHERE id = $1 RETURNING review_id, position`,
		photoID,
	).Scan(&reviewID, &pos)
	if err != nil {
		return err
	}

	if pos == 0 {
		var nextURL sql.NullString
		err = tx.QueryRow(
			`SELECT url FROM review_photos WHERE review_id = $1 ORDER BY position ASC, id ASC LIMIT 1`,
			reviewID,
		).Scan(&nextURL)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if nextURL.Valid {
			_, err = tx.Exec(`UPDATE reviews SET image_url = $1, updated_at = now() WHERE id = $2`, nextURL.String, reviewID)
		} else {
			_, err = tx.Exec(`UPDATE reviews SET image_url = NULL, updated_at = now() WHERE id = $1`, reviewID)
		}
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ReviewRepo) UpdateVideoURL(reviewID int, videoURL string) error {
	_, err := r.db.Exec(`UPDATE reviews SET video_url = $1, updated_at = now() WHERE id = $2`, videoURL, reviewID)
	return err
}

func (r *ReviewRepo) getAuthors(reviewID int) ([]model.User, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.username, u.avatar_url, u.role, u.created_at, u.updated_at
		FROM users u
		JOIN review_authors ra ON ra.user_id = u.id
		WHERE ra.review_id = $1
	`, reviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		authors = append(authors, u)
	}
	return authors, rows.Err()
}
