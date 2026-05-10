package repository

import (
	"database/sql"

	"github.com/aeva-eat/backend/internal/model"
)

// AggregateRepo — sql-агрегации для /api/cities, /api/users/:id, /api/gems.
// Все запросы read-only, не нуждаются в транзакциях.
type AggregateRepo struct {
	db *sql.DB
}

func NewAggregateRepo(db *sql.DB) *AggregateRepo {
	return &AggregateRepo{db: db}
}

// Cities — список городов круга с агрегатами. Используется в полке "По
// городам" на Найти и на /cities/:name.
func (r *AggregateRepo) Cities() ([]model.CityAggregate, error) {
	rows, err := r.db.Query(`
		SELECT
			p.city                                                           AS city,
			COUNT(DISTINCT p.id)                                             AS count,
			COUNT(DISTINCT CASE WHEN gem_marks.place_id IS NOT NULL THEN p.id END) AS gem_count,
			COUNT(DISTINCT ra.user_id)                                       AS contributor_count
		FROM places p
		LEFT JOIN reviews rv ON rv.place_id = p.id
		LEFT JOIN review_authors ra ON ra.review_id = rv.id
		LEFT JOIN (
			SELECT DISTINCT place_id FROM reviews WHERE is_gem = true
		) gem_marks ON gem_marks.place_id = p.id
		WHERE p.city IS NOT NULL AND p.city <> ''
		GROUP BY p.city
		ORDER BY count DESC, p.city ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.CityAggregate{}
	for rows.Next() {
		var c model.CityAggregate
		if err := rows.Scan(&c.City, &c.Count, &c.GemCount, &c.ContributorCount); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// City — один город по имени (case-insensitive). Возвращает nil, nil если
// в этом городе нет мест.
func (r *AggregateRepo) City(name string) (*model.CityAggregate, error) {
	row := r.db.QueryRow(`
		SELECT
			MIN(p.city)                                                       AS city,
			COUNT(DISTINCT p.id)                                              AS count,
			COUNT(DISTINCT CASE WHEN gem_marks.place_id IS NOT NULL THEN p.id END) AS gem_count,
			COUNT(DISTINCT ra.user_id)                                        AS contributor_count
		FROM places p
		LEFT JOIN reviews rv ON rv.place_id = p.id
		LEFT JOIN review_authors ra ON ra.review_id = rv.id
		LEFT JOIN (
			SELECT DISTINCT place_id FROM reviews WHERE is_gem = true
		) gem_marks ON gem_marks.place_id = p.id
		WHERE LOWER(p.city) = LOWER($1)
	`, name)

	var c model.CityAggregate
	var maybeName sql.NullString
	if err := row.Scan(&maybeName, &c.Count, &c.GemCount, &c.ContributorCount); err != nil {
		return nil, err
	}
	if !maybeName.Valid || c.Count == 0 {
		return nil, nil
	}
	c.City = maybeName.String
	return &c, nil
}

// UsersList — все пользователи круга для полки "По друзьям" (DESIGN-DECISIONS
// §F1) и страницы /people/. Сортировка по активности (review_count desc) —
// сверху самые вовлечённые.
func (r *AggregateRepo) UsersList() ([]model.UserProfile, error) {
	rows, err := r.db.Query(`
		SELECT
			u.id, u.username, u.display_name, u.avatar_url,
			COALESCE(place_stats.place_count, 0) AS place_count,
			COALESCE(place_stats.gem_count, 0)   AS gem_count,
			COALESCE(city_stats.city_count, 0)   AS city_count,
			COALESCE(rev_stats.review_count, 0)  AS review_count
		FROM users u
		LEFT JOIN (
			SELECT ra.user_id,
				COUNT(DISTINCT rv.place_id) AS place_count,
				COUNT(DISTINCT CASE WHEN rv.is_gem THEN rv.place_id END) AS gem_count
			FROM reviews rv
			JOIN review_authors ra ON ra.review_id = rv.id
			GROUP BY ra.user_id
		) place_stats ON place_stats.user_id = u.id
		LEFT JOIN (
			SELECT ra.user_id, COUNT(DISTINCT p.city) AS city_count
			FROM reviews rv
			JOIN review_authors ra ON ra.review_id = rv.id
			JOIN places p ON p.id = rv.place_id
			WHERE p.city IS NOT NULL AND p.city <> ''
			GROUP BY ra.user_id
		) city_stats ON city_stats.user_id = u.id
		LEFT JOIN (
			SELECT user_id, COUNT(*) AS review_count
			FROM review_authors GROUP BY user_id
		) rev_stats ON rev_stats.user_id = u.id
		ORDER BY review_count DESC, u.username ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.UserProfile{}
	for rows.Next() {
		var p model.UserProfile
		if err := rows.Scan(
			&p.ID, &p.Username, &p.DisplayName, &p.AvatarURL,
			&p.PlaceCount, &p.GemCount, &p.CityCount, &p.ReviewCount,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// UserProfile — публичная статистика пользователя. Возвращает nil, nil если
// пользователя нет.
func (r *AggregateRepo) UserProfile(userID int) (*model.UserProfile, error) {
	row := r.db.QueryRow(`
		SELECT
			u.id, u.username, u.display_name, u.avatar_url,
			(SELECT COUNT(DISTINCT rv.place_id)
			   FROM reviews rv
			   JOIN review_authors ra ON ra.review_id = rv.id
			   WHERE ra.user_id = u.id)                          AS place_count,
			(SELECT COUNT(DISTINCT rv.place_id)
			   FROM reviews rv
			   JOIN review_authors ra ON ra.review_id = rv.id
			   WHERE ra.user_id = u.id AND rv.is_gem = true)     AS gem_count,
			(SELECT COUNT(DISTINCT p.city)
			   FROM reviews rv
			   JOIN review_authors ra ON ra.review_id = rv.id
			   JOIN places p ON p.id = rv.place_id
			   WHERE ra.user_id = u.id AND p.city IS NOT NULL AND p.city <> '') AS city_count,
			(SELECT COUNT(*)
			   FROM review_authors ra
			   WHERE ra.user_id = u.id)                          AS review_count
		FROM users u
		WHERE u.id = $1
	`, userID)

	var p model.UserProfile
	if err := row.Scan(
		&p.ID, &p.Username, &p.DisplayName, &p.AvatarURL,
		&p.PlaceCount, &p.GemCount, &p.CityCount, &p.ReviewCount,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// DESIGN-DECISIONS Q8: «любит грузинскую — 11 раз» — берём самую
	// частую кухню по визитам пользователя. Считаем по review_authors,
	// чтобы соавторы тоже получили эту кухню в свой профиль.
	cuisineRow := r.db.QueryRow(`
		SELECT ct.name, COUNT(*)
		FROM reviews rv
		JOIN review_authors ra ON ra.review_id = rv.id
		JOIN places p ON p.id = rv.place_id
		JOIN cuisine_types ct ON ct.id = p.cuisine_type_id
		WHERE ra.user_id = $1
		GROUP BY ct.id, ct.name
		ORDER BY COUNT(*) DESC, ct.name ASC
		LIMIT 1
	`, userID)
	var cuisine string
	var cuisineCount int
	if err := cuisineRow.Scan(&cuisine, &cuisineCount); err == nil {
		p.FavoriteCuisine = &cuisine
		p.FavoriteCuisineCount = cuisineCount
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	return &p, nil
}

// UserPlaceIDs — id мест, в которых пользователь был. Не возвращает сами
// объекты — далее их грузит обычный PlaceRepo с join'ом фото/авторов.
func (r *AggregateRepo) UserPlaceIDs(userID int, gemsOnly bool) ([]int, error) {
	q := `
		SELECT DISTINCT rv.place_id
		FROM reviews rv
		JOIN review_authors ra ON ra.review_id = rv.id
		WHERE ra.user_id = $1
	`
	if gemsOnly {
		q += " AND rv.is_gem = true"
	}
	q += " ORDER BY rv.place_id DESC"

	rows, err := r.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// UserCities — города пользователя со счётчиком визитов в каждый.
func (r *AggregateRepo) UserCities(userID int) ([]model.UserCity, error) {
	rows, err := r.db.Query(`
		SELECT p.city, COUNT(DISTINCT p.id) AS count
		FROM reviews rv
		JOIN review_authors ra ON ra.review_id = rv.id
		JOIN places p ON p.id = rv.place_id
		WHERE ra.user_id = $1 AND p.city IS NOT NULL AND p.city <> ''
		GROUP BY p.city
		ORDER BY count DESC, p.city ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.UserCity{}
	for rows.Next() {
		var c model.UserCity
		if err := rows.Scan(&c.City, &c.Count); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// GemPlaceIDs — id всех мест-жемчужин, в порядке свежести (по самому
// свежему gem-отзыву).
func (r *AggregateRepo) GemPlaceIDs() ([]int, error) {
	rows, err := r.db.Query(`
		SELECT rv.place_id
		FROM reviews rv
		WHERE rv.is_gem = true
		GROUP BY rv.place_id
		ORDER BY MAX(rv.created_at) DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// GemsByCity — города в порядке количества жемчужин в каждом.
func (r *AggregateRepo) GemsByCity() ([]model.CityAggregate, error) {
	rows, err := r.db.Query(`
		SELECT p.city,
			COUNT(DISTINCT p.id)                AS count,
			COUNT(DISTINCT p.id)                AS gem_count,
			COUNT(DISTINCT ra.user_id)          AS contributor_count
		FROM places p
		JOIN reviews rv ON rv.place_id = p.id AND rv.is_gem = true
		JOIN review_authors ra ON ra.review_id = rv.id
		WHERE p.city IS NOT NULL AND p.city <> ''
		GROUP BY p.city
		ORDER BY gem_count DESC, p.city ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.CityAggregate{}
	for rows.Next() {
		var c model.CityAggregate
		if err := rows.Scan(&c.City, &c.Count, &c.GemCount, &c.ContributorCount); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// GemsByUser — авторы, отметившие наибольшее число жемчужин.
func (r *AggregateRepo) GemsByUser() ([]model.UserGemCount, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.username, u.avatar_url, COUNT(DISTINCT rv.place_id) AS gem_count
		FROM users u
		JOIN review_authors ra ON ra.user_id = u.id
		JOIN reviews rv ON rv.id = ra.review_id AND rv.is_gem = true
		GROUP BY u.id, u.username, u.avatar_url
		ORDER BY gem_count DESC, u.username ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.UserGemCount{}
	for rows.Next() {
		var c model.UserGemCount
		if err := rows.Scan(&c.UserID, &c.Username, &c.AvatarURL, &c.Count); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
