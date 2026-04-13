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

type PlaceFilter struct {
	City          string
	CuisineTypeID int
	CategoryID    int
	MinRating     float64
	IsGem         *bool
	Search        string
	Sort          string
}

func (r *PlaceRepo) List(f PlaceFilter) ([]model.Place, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.address, p.city, p.lat, p.lng,
			p.cuisine_type_id, ct.name AS cuisine_type, p.website,
			p.created_by, p.image_url, p.created_at, p.updated_at,
			COALESCE(rs.avg_food, 0), COALESCE(rs.avg_service, 0), COALESCE(rs.avg_vibe, 0),
			COALESCE(rs.review_count, 0)
		FROM places p
		LEFT JOIN cuisine_types ct ON ct.id = p.cuisine_type_id
		LEFT JOIN place_categories pc ON pc.place_id = p.id
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
	if f.CuisineTypeID > 0 {
		conditions = append(conditions, fmt.Sprintf("p.cuisine_type_id = $%d", argIdx))
		args = append(args, f.CuisineTypeID)
		argIdx++
	}
	if f.CategoryID > 0 {
		conditions = append(conditions, fmt.Sprintf("pc.category_id = $%d", argIdx))
		args = append(args, f.CategoryID)
		argIdx++
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
	if f.Search != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(p.name) LIKE LOWER($%d)", argIdx))
		args = append(args, "%"+f.Search+"%")
		argIdx++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	switch f.Sort {
	case "rating":
		query += " ORDER BY (COALESCE(rs.avg_food,0) + COALESCE(rs.avg_service,0) + COALESCE(rs.avg_vibe,0)) DESC"
	case "name":
		query += " ORDER BY p.name"
	default:
		query += " ORDER BY p.created_at DESC"
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
	return places, rows.Err()
}

func (r *PlaceRepo) GetByID(id int) (*model.Place, error) {
	p := &model.Place{}
	var avgFood, avgService, avgVibe float64
	err := r.db.QueryRow(`
		SELECT p.id, p.name, p.address, p.city, p.lat, p.lng,
			p.cuisine_type_id, ct.name AS cuisine_type, p.website,
			p.created_by, p.image_url, p.created_at, p.updated_at,
			COALESCE(rs.avg_food, 0), COALESCE(rs.avg_service, 0), COALESCE(rs.avg_vibe, 0),
			COALESCE(rs.review_count, 0)
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
