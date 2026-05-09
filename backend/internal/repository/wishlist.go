package repository

import (
	"database/sql"

	"github.com/aeva-eat/backend/internal/model"
)

type WishlistRepo struct {
	db *sql.DB
}

func NewWishlistRepo(db *sql.DB) *WishlistRepo {
	return &WishlistRepo{db: db}
}

func (r *WishlistRepo) Add(userID, placeID int) error {
	_, err := r.db.Exec(`
		INSERT INTO wishlists (user_id, place_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, userID, placeID)
	return err
}

// MarkStruck помечает запись wishlist как реализованную (визит произошёл).
// Idempotent: повторный вызов не двигает struck_at, если уже true.
// Возвращает true, если запись была найдена и обновлена впервые.
func (r *WishlistRepo) MarkStruck(userID, placeID int) (bool, error) {
	res, err := r.db.Exec(`
		UPDATE wishlists
		   SET is_struck = TRUE, struck_at = COALESCE(struck_at, now())
		 WHERE user_id = $1 AND place_id = $2 AND is_struck = FALSE
	`, userID, placeID)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// ListAll — общий wishlist круга (backend.md §Wishlist). Свежие сверху,
// зачёркнутые в конце (внутри своих групп — по дате).
func (r *WishlistRepo) ListAll() ([]model.WishlistEntry, error) {
	rows, err := r.db.Query(`
		SELECT
			w.user_id, u.username, u.avatar_url,
			w.is_struck, w.struck_at, w.created_at,
			p.id, p.name, p.address, p.city, p.lat, p.lng,
			p.cuisine_type_id, ct.name AS cuisine_type, p.website,
			p.created_by,
			COALESCE(
				p.image_url,
				(SELECT r.image_url FROM reviews r
				  WHERE r.place_id = p.id AND r.image_url IS NOT NULL
				  ORDER BY r.created_at ASC, r.id ASC LIMIT 1)
			) AS image_url,
			p.created_at, p.updated_at,
			COALESCE(rs.avg_food, 0), COALESCE(rs.avg_service, 0), COALESCE(rs.avg_vibe, 0),
			COALESCE(rs.review_count, 0),
			EXISTS(SELECT 1 FROM reviews rv WHERE rv.place_id = p.id AND rv.is_gem = true) AS is_gem_place
		FROM wishlists w
		JOIN users u  ON u.id = w.user_id
		JOIN places p ON p.id = w.place_id
		LEFT JOIN cuisine_types ct ON ct.id = p.cuisine_type_id
		LEFT JOIN (
			SELECT place_id,
				AVG(food_rating)::numeric(3,1) AS avg_food,
				AVG(service_rating)::numeric(3,1) AS avg_service,
				AVG(vibe_rating)::numeric(3,1) AS avg_vibe,
				COUNT(*) AS review_count
			FROM reviews GROUP BY place_id
		) rs ON rs.place_id = p.id
		ORDER BY w.is_struck ASC, w.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.WishlistEntry{}
	for rows.Next() {
		var e model.WishlistEntry
		var avgFood, avgService, avgVibe float64
		err := rows.Scan(
			&e.UserID, &e.Username, &e.AvatarURL,
			&e.IsStruck, &e.StruckAt, &e.CreatedAt,
			&e.Place.ID, &e.Place.Name, &e.Place.Address, &e.Place.City,
			&e.Place.Lat, &e.Place.Lng,
			&e.Place.CuisineTypeID, &e.Place.CuisineType, &e.Place.Website,
			&e.Place.CreatedBy, &e.Place.ImageURL,
			&e.Place.CreatedAt, &e.Place.UpdatedAt,
			&avgFood, &avgService, &avgVibe, &e.Place.ReviewCount,
			&e.Place.IsGemPlace,
		)
		if err != nil {
			return nil, err
		}
		if e.Place.ReviewCount > 0 {
			e.Place.AvgFood = &avgFood
			e.Place.AvgService = &avgService
			e.Place.AvgVibe = &avgVibe
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *WishlistRepo) Remove(userID, placeID int) error {
	_, err := r.db.Exec(`DELETE FROM wishlists WHERE user_id = $1 AND place_id = $2`, userID, placeID)
	return err
}

func (r *WishlistRepo) ListByUser(userID int) ([]model.Place, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.address, p.city, p.lat, p.lng,
			p.cuisine_type_id, ct.name AS cuisine_type, p.website,
			p.created_by, p.image_url, p.created_at, p.updated_at,
			COALESCE(rs.avg_food, 0), COALESCE(rs.avg_service, 0), COALESCE(rs.avg_vibe, 0),
			COALESCE(rs.review_count, 0)
		FROM wishlists w
		JOIN places p ON p.id = w.place_id
		LEFT JOIN cuisine_types ct ON ct.id = p.cuisine_type_id
		LEFT JOIN (
			SELECT place_id,
				AVG(food_rating)::numeric(3,1) AS avg_food,
				AVG(service_rating)::numeric(3,1) AS avg_service,
				AVG(vibe_rating)::numeric(3,1) AS avg_vibe,
				COUNT(*) AS review_count
			FROM reviews GROUP BY place_id
		) rs ON rs.place_id = p.id
		WHERE w.user_id = $1
		ORDER BY w.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var places []model.Place
	for rows.Next() {
		var p model.Place
		var avgFood, avgService, avgVibe float64
		err := rows.Scan(
			&p.ID, &p.Name, &p.Address, &p.City, &p.Lat, &p.Lng,
			&p.CuisineTypeID, &p.CuisineType, &p.Website,
			&p.CreatedBy, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt,
			&avgFood, &avgService, &avgVibe, &p.ReviewCount,
		)
		if err != nil {
			return nil, err
		}
		if p.ReviewCount > 0 {
			p.AvgFood = &avgFood
			p.AvgService = &avgService
			p.AvgVibe = &avgVibe
		}
		places = append(places, p)
	}
	return places, rows.Err()
}

func (r *WishlistRepo) IsWishlisted(userID, placeID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM wishlists WHERE user_id=$1 AND place_id=$2)`,
		userID, placeID,
	).Scan(&exists)
	return exists, err
}

func (r *WishlistRepo) GetUserWishlistIDs(userID int) ([]int, error) {
	rows, err := r.db.Query(`SELECT place_id FROM wishlists WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// Custom wishlist (free-text entries)

func (r *WishlistRepo) AddCustom(userID int, name string, note *string) (*model.WishlistCustom, error) {
	wc := &model.WishlistCustom{}
	err := r.db.QueryRow(`
		INSERT INTO wishlist_custom (user_id, name, note) VALUES ($1, $2, $3)
		RETURNING id, user_id, name, note, created_at
	`, userID, name, note).Scan(&wc.ID, &wc.UserID, &wc.Name, &wc.Note, &wc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return wc, nil
}

func (r *WishlistRepo) ListCustom(userID int) ([]model.WishlistCustom, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, name, note, created_at
		FROM wishlist_custom WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.WishlistCustom
	for rows.Next() {
		var wc model.WishlistCustom
		if err := rows.Scan(&wc.ID, &wc.UserID, &wc.Name, &wc.Note, &wc.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, wc)
	}
	return items, rows.Err()
}

func (r *WishlistRepo) DeleteCustom(userID, id int) error {
	_, err := r.db.Exec(`DELETE FROM wishlist_custom WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}
