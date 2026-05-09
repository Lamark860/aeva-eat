package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/aeva-eat/backend/internal/model"
)

type PlaceRepo struct {
	db *sql.DB
}

func NewPlaceRepo(db *sql.DB) *PlaceRepo {
	return &PlaceRepo{db: db}
}

// escapeLike escapes characters that have special meaning inside an SQL LIKE pattern
// when the query uses ESCAPE '\'. Order matters: backslash first.
func escapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}

type PlaceFilter struct {
	City           string
	CuisineTypeIDs []int
	CategoryIDs    []int
	MinRating      float64
	IsGem          *bool
	Search         string
	Sort           string
	Limit          int
	Offset         int
}

type PlaceListResult struct {
	Places []model.Place `json:"places"`
	Total  int           `json:"total"`
}

func (r *PlaceRepo) List(f PlaceFilter) (*PlaceListResult, error) {
	// place_categories is M2M, so we used to JOIN it and rely on DISTINCT to dedupe.
	// That broke sort by rating: `ORDER BY (rs.avg_food IS NULL)` violates the
	// Postgres rule that ORDER BY expressions must appear in SELECT under DISTINCT.
	// Switching the category filter to EXISTS removes the need for both the JOIN
	// and DISTINCT; ct and rs are 1:1 with p so no duplicates remain.
	baseFrom := `
		FROM places p
		LEFT JOIN cuisine_types ct ON ct.id = p.cuisine_type_id
		LEFT JOIN (
			SELECT place_id,
				AVG(food_rating)::numeric(3,1) AS avg_food,
				AVG(service_rating)::numeric(3,1) AS avg_service,
				AVG(vibe_rating)::numeric(3,1) AS avg_vibe,
				COUNT(*) AS review_count
			FROM reviews GROUP BY place_id
		) rs ON rs.place_id = p.id
	`

	var conditions []string
	var args []interface{}
	argIdx := 1

	if f.City != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(p.city) = LOWER($%d)", argIdx))
		args = append(args, f.City)
		argIdx++
	}
	if len(f.CuisineTypeIDs) > 0 {
		placeholders := make([]string, len(f.CuisineTypeIDs))
		for i, id := range f.CuisineTypeIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIdx)
			args = append(args, id)
			argIdx++
		}
		conditions = append(conditions, fmt.Sprintf("p.cuisine_type_id IN (%s)", strings.Join(placeholders, ",")))
	}
	if len(f.CategoryIDs) > 0 {
		placeholders := make([]string, len(f.CategoryIDs))
		for i, id := range f.CategoryIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIdx)
			args = append(args, id)
			argIdx++
		}
		conditions = append(conditions, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM place_categories pc WHERE pc.place_id = p.id AND pc.category_id IN (%s))",
			strings.Join(placeholders, ",")))
	}
	if f.MinRating > 0 {
		conditions = append(conditions, fmt.Sprintf(
			"(COALESCE(rs.avg_food,0) + COALESCE(rs.avg_service,0) + COALESCE(rs.avg_vibe,0)) / 3.0 >= $%d", argIdx))
		args = append(args, f.MinRating)
		argIdx++
	}
	if f.IsGem != nil && *f.IsGem {
		conditions = append(conditions, `EXISTS (SELECT 1 FROM reviews rv WHERE rv.place_id = p.id AND rv.is_gem = true)`)
	}
	if s := strings.TrimSpace(f.Search); s != "" {
		// Escape LIKE wildcards so user-typed % and _ match literally, not as wildcards.
		conditions = append(conditions, fmt.Sprintf("LOWER(p.name) LIKE LOWER($%d) ESCAPE '\\'", argIdx))
		args = append(args, "%"+escapeLike(s)+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseFrom + whereClause
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count places: %w", err)
	}

	// Main query
	query := `
		SELECT p.id, p.name, p.address, p.city, p.lat, p.lng,
			p.cuisine_type_id, ct.name AS cuisine_type, p.website,
			p.created_by,
			COALESCE(
				p.image_url,
				(SELECT r.image_url
				   FROM reviews r
				  WHERE r.place_id = p.id AND r.image_url IS NOT NULL
				  ORDER BY r.created_at ASC, r.id ASC
				  LIMIT 1)
			) AS image_url,
			p.created_at, p.updated_at,
			COALESCE(rs.avg_food, 0), COALESCE(rs.avg_service, 0), COALESCE(rs.avg_vibe, 0),
			COALESCE(rs.review_count, 0),
			EXISTS(SELECT 1 FROM reviews rv WHERE rv.place_id = p.id AND rv.is_gem = true) AS is_gem_place
	` + baseFrom + whereClause

	// All ORDER BY clauses end with `p.id DESC` as a stable tiebreaker — without it,
	// rows with equal sort keys return in undefined order and pagination becomes unstable
	// (the same place can appear on two pages, or vanish between pages).
	switch f.Sort {
	case "rating":
		query += ` ORDER BY (rs.avg_food IS NULL) ASC, (COALESCE(rs.avg_food,0) + COALESCE(rs.avg_service,0) + COALESCE(rs.avg_vibe,0)) DESC, p.id DESC`
	case "rating_asc":
		query += ` ORDER BY (rs.avg_food IS NULL) ASC, (COALESCE(rs.avg_food,0) + COALESCE(rs.avg_service,0) + COALESCE(rs.avg_vibe,0)) ASC, p.id DESC`
	case "name":
		query += " ORDER BY p.name, p.id DESC"
	default:
		query += " ORDER BY p.created_at DESC, p.id DESC"
	}

	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
		args = append(args, f.Limit, f.Offset)
		argIdx += 2
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query places: %w", err)
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
			&p.IsGemPlace,
		)
		if err != nil {
			return nil, fmt.Errorf("scan place: %w", err)
		}
		if p.ReviewCount > 0 {
			p.AvgFood = &avgFood
			p.AvgService = &avgService
			p.AvgVibe = &avgVibe
		}
		places = append(places, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load reviewers for each place
	for i := range places {
		reviewers, err := r.getReviewers(places[i].ID)
		if err != nil {
			return nil, err
		}
		places[i].Reviewers = reviewers
	}

	return &PlaceListResult{Places: places, Total: total}, nil
}

func (r *PlaceRepo) GetByID(id int) (*model.Place, error) {
	p := &model.Place{}
	var avgFood, avgService, avgVibe float64
	err := r.db.QueryRow(`
		SELECT p.id, p.name, p.address, p.city, p.lat, p.lng,
			p.cuisine_type_id, ct.name AS cuisine_type, p.website,
			p.created_by,
			COALESCE(
				p.image_url,
				(SELECT r.image_url
				   FROM reviews r
				  WHERE r.place_id = p.id AND r.image_url IS NOT NULL
				  ORDER BY r.created_at ASC, r.id ASC
				  LIMIT 1)
			) AS image_url,
			p.created_at, p.updated_at,
			COALESCE(rs.avg_food, 0), COALESCE(rs.avg_service, 0), COALESCE(rs.avg_vibe, 0),
			COALESCE(rs.review_count, 0),
			EXISTS(SELECT 1 FROM reviews rv WHERE rv.place_id = p.id AND rv.is_gem = true) AS is_gem_place
		FROM places p
		LEFT JOIN cuisine_types ct ON ct.id = p.cuisine_type_id
		LEFT JOIN (
			SELECT place_id,
				AVG(food_rating)::numeric(3,1) AS avg_food,
				AVG(service_rating)::numeric(3,1) AS avg_service,
				AVG(vibe_rating)::numeric(3,1) AS avg_vibe,
				COUNT(*) AS review_count
			FROM reviews GROUP BY place_id
		) rs ON rs.place_id = p.id
		WHERE p.id = $1
	`, id).Scan(
		&p.ID, &p.Name, &p.Address, &p.City, &p.Lat, &p.Lng,
		&p.CuisineTypeID, &p.CuisineType, &p.Website,
		&p.CreatedBy, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt,
		&avgFood, &avgService, &avgVibe, &p.ReviewCount,
		&p.IsGemPlace,
	)
	if err != nil {
		return nil, err
	}
	if p.ReviewCount > 0 {
		p.AvgFood = &avgFood
		p.AvgService = &avgService
		p.AvgVibe = &avgVibe
	}

	// Load categories
	rows, err := r.db.Query(
		`SELECT c.name FROM categories c JOIN place_categories pc ON pc.category_id = c.id WHERE pc.place_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		p.Categories = append(p.Categories, name)
	}

	// Load reviewers
	reviewers, err := r.getReviewers(id)
	if err != nil {
		return nil, err
	}
	p.Reviewers = reviewers

	return p, rows.Err()
}

func (r *PlaceRepo) Create(p *model.Place, categoryIDs []int) (*model.Place, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
		INSERT INTO places (name, address, city, lat, lng, cuisine_type_id, website, image_url, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`, p.Name, p.Address, p.City, p.Lat, p.Lng, p.CuisineTypeID, p.Website, p.ImageURL, p.CreatedBy,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	for _, catID := range categoryIDs {
		_, err := tx.Exec(`INSERT INTO place_categories (place_id, category_id) VALUES ($1, $2)`, p.ID, catID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(p.ID)
}

func (r *PlaceRepo) Update(p *model.Place, categoryIDs []int) (*model.Place, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE places SET name=$1, address=$2, city=$3, lat=$4, lng=$5,
			cuisine_type_id=$6, website=$7, image_url=$8, updated_at=now()
		WHERE id=$9
	`, p.Name, p.Address, p.City, p.Lat, p.Lng, p.CuisineTypeID, p.Website, p.ImageURL, p.ID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`DELETE FROM place_categories WHERE place_id = $1`, p.ID)
	if err != nil {
		return nil, err
	}
	for _, catID := range categoryIDs {
		_, err := tx.Exec(`INSERT INTO place_categories (place_id, category_id) VALUES ($1, $2)`, p.ID, catID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(p.ID)
}

func (r *PlaceRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM places WHERE id = $1`, id)
	return err
}

func (r *PlaceRepo) GetOwnerID(placeID int) (int, error) {
	var ownerID int
	err := r.db.QueryRow(`SELECT COALESCE(created_by, 0) FROM places WHERE id = $1`, placeID).Scan(&ownerID)
	return ownerID, err
}
func (r *PlaceRepo) UpdateImageURL(placeID int, imageURL string) error {
	_, err := r.db.Exec(`UPDATE places SET image_url = $1, updated_at = now() WHERE id = $2`, imageURL, placeID)
	return err
}

func (r *PlaceRepo) ListCities() ([]string, error) {
	rows, err := r.db.Query(`SELECT DISTINCT city FROM places WHERE city IS NOT NULL AND city != '' ORDER BY city`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	return cities, rows.Err()
}

func (r *PlaceRepo) getReviewers(placeID int) ([]model.Reviewer, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT u.id, u.username, u.avatar_url
		FROM users u
		JOIN review_authors ra ON ra.user_id = u.id
		JOIN reviews rv ON rv.id = ra.review_id
		WHERE rv.place_id = $1
		ORDER BY u.username
	`, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []model.Reviewer
	for rows.Next() {
		var rev model.Reviewer
		if err := rows.Scan(&rev.ID, &rev.Username, &rev.AvatarURL); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, rev)
	}
	return reviewers, rows.Err()
}
