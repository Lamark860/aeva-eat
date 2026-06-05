package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aeva-eat/backend/internal/model"
	"github.com/lib/pq"
)

type PlaceRepo struct {
	db *sql.DB
}

// DuplicatePlaceError возвращается Create/Update, когда место с такой же
// идентичностью (name + address + city, см. migrations/016) уже существует.
// Existing хранит конфликтующую запись (может быть nil, если её не удалось
// перечитать) — фронт показывает «это оно? → перейти и оставить отзыв».
type DuplicatePlaceError struct {
	Existing *model.Place
}

func (e *DuplicatePlaceError) Error() string { return "place already exists" }

// isUniqueViolation — Postgres unique_violation (23505). Не завязано на имя
// индекса, поэтому переименование индекса ничего не ломает.
func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
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
	// B5 — расширенные фильтры (backend.md §/api/places).
	// AttendedBy: места, в которых хотя бы один из этих пользователей был.
	// VisitFrom/To: окно по дате визита (RFC3339-date YYYY-MM-DD).
	// SortRatingUserID: при Sort=="rating_user" сортируем по среднему
	// рейтингу конкретного пользователя (а не среднего по кругу).
	AttendedBy       []int
	VisitFrom        string
	VisitTo          string
	SortRatingUserID int
	Limit            int
	Offset           int
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

	// Базовое условие: показываем только активные (не soft-deleted) места.
	conditions := []string{"p.deleted_at IS NULL"}
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
		// Ищем по названию, городу И кухне (плейсхолдер обещает «место, кухня, город»).
		// $argIdx переиспользуется во всех трёх — это валидно, параметр один.
		conditions = append(conditions, fmt.Sprintf(
			`(LOWER(p.name) LIKE LOWER($%[1]d) ESCAPE '\'
			  OR LOWER(COALESCE(p.city, '')) LIKE LOWER($%[1]d) ESCAPE '\'
			  OR LOWER(COALESCE(ct.name, '')) LIKE LOWER($%[1]d) ESCAPE '\')`, argIdx))
		args = append(args, "%"+escapeLike(s)+"%")
		argIdx++
	}
	if len(f.AttendedBy) > 0 {
		placeholders := make([]string, len(f.AttendedBy))
		for i, id := range f.AttendedBy {
			placeholders[i] = fmt.Sprintf("$%d", argIdx)
			args = append(args, id)
			argIdx++
		}
		conditions = append(conditions, fmt.Sprintf(
			`EXISTS (SELECT 1 FROM reviews rv JOIN review_authors ra ON ra.review_id = rv.id
			           WHERE rv.place_id = p.id AND ra.user_id IN (%s))`,
			strings.Join(placeholders, ",")))
	}
	if f.VisitFrom != "" {
		conditions = append(conditions, fmt.Sprintf(
			`EXISTS (SELECT 1 FROM reviews rv WHERE rv.place_id = p.id
			          AND rv.visited_at IS NOT NULL AND rv.visited_at >= $%d)`, argIdx))
		args = append(args, f.VisitFrom)
		argIdx++
	}
	if f.VisitTo != "" {
		conditions = append(conditions, fmt.Sprintf(
			`EXISTS (SELECT 1 FROM reviews rv WHERE rv.place_id = p.id
			          AND rv.visited_at IS NOT NULL AND rv.visited_at <= $%d)`, argIdx))
		args = append(args, f.VisitTo)
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
			EXISTS(SELECT 1 FROM reviews rv WHERE rv.place_id = p.id AND rv.is_gem = true) AS is_gem_place,
			EXISTS(SELECT 1 FROM reviews rv WHERE rv.place_id = p.id AND rv.video_url IS NOT NULL AND rv.video_url <> '') AS has_video,
			(SELECT rv.video_url FROM reviews rv
			   WHERE rv.place_id = p.id AND rv.video_url IS NOT NULL AND rv.video_url <> ''
			   ORDER BY rv.created_at DESC, rv.id DESC LIMIT 1) AS video_url,
			COALESCE(
				(SELECT array_agg(rv.video_url ORDER BY rv.created_at DESC, rv.id DESC)
				   FROM reviews rv WHERE rv.place_id = p.id
				   AND rv.video_url IS NOT NULL AND rv.video_url <> ''),
				ARRAY[]::TEXT[]
			) AS videos,
			(SELECT rv.comment FROM reviews rv
			   WHERE rv.place_id = p.id AND rv.comment IS NOT NULL
			   AND LENGTH(rv.comment) >= 30
			   ORDER BY LENGTH(rv.comment) DESC, rv.created_at ASC, rv.id ASC
			   LIMIT 1) AS top_review_comment
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
	case "rating_user":
		// Сортировка по среднему рейтингу конкретного пользователя.
		// Подзапрос вытаскивает avg(всех 3-х рейтингов) для $userID;
		// места без его отзыва идут в конце (NULLS LAST).
		query += fmt.Sprintf(`
			ORDER BY (
				SELECT (rv.food_rating + rv.service_rating + rv.vibe_rating) / 3.0
				FROM reviews rv
				JOIN review_authors ra ON ra.review_id = rv.id
				WHERE rv.place_id = p.id AND ra.user_id = $%d
				ORDER BY rv.created_at DESC LIMIT 1
			) DESC NULLS LAST, p.id DESC`, argIdx)
		args = append(args, f.SortRatingUserID)
		argIdx++
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
		var videos pq.StringArray
		err := rows.Scan(
			&p.ID, &p.Name, &p.Address, &p.City, &p.Lat, &p.Lng,
			&p.CuisineTypeID, &p.CuisineType, &p.Website,
			&p.CreatedBy, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt,
			&avgFood, &avgService, &avgVibe, &p.ReviewCount,
			&p.IsGemPlace, &p.HasVideo, &p.VideoURL, &videos,
			&p.TopReviewComment,
		)
		if err != nil {
			return nil, fmt.Errorf("scan place: %w", err)
		}
		p.Videos = []string(videos)
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

	// Reviewers + feedPhotos одним батч-запросом каждый (вместо N+1 цикла по
	// местам) — главный листинг «Найти» был самым горячим путём.
	if len(places) > 0 {
		ids := make([]int, len(places))
		for i := range places {
			ids[i] = places[i].ID
		}
		reviewersByPlace, err := r.getReviewersBatch(ids)
		if err != nil {
			return nil, err
		}
		photosByPlace, err := r.getFeedPhotosBatch(ids)
		if err != nil {
			return nil, err
		}
		for i := range places {
			places[i].Reviewers = reviewersByPlace[places[i].ID]
			places[i].FeedPhotos = photosByPlace[places[i].ID]
		}
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
			EXISTS(SELECT 1 FROM reviews rv WHERE rv.place_id = p.id AND rv.is_gem = true) AS is_gem_place,
			EXISTS(SELECT 1 FROM reviews rv WHERE rv.place_id = p.id AND rv.video_url IS NOT NULL AND rv.video_url <> '') AS has_video,
			(SELECT rv.video_url FROM reviews rv
			   WHERE rv.place_id = p.id AND rv.video_url IS NOT NULL AND rv.video_url <> ''
			   ORDER BY rv.created_at DESC, rv.id DESC LIMIT 1) AS video_url,
			COALESCE(
				(SELECT array_agg(rv.video_url ORDER BY rv.created_at DESC, rv.id DESC)
				   FROM reviews rv WHERE rv.place_id = p.id
				   AND rv.video_url IS NOT NULL AND rv.video_url <> ''),
				ARRAY[]::TEXT[]
			) AS videos,
			(SELECT rv.comment FROM reviews rv
			   WHERE rv.place_id = p.id AND rv.comment IS NOT NULL
			   AND LENGTH(rv.comment) >= 30
			   ORDER BY LENGTH(rv.comment) DESC, rv.created_at ASC, rv.id ASC
			   LIMIT 1) AS top_review_comment
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
		&p.IsGemPlace, &p.HasVideo, &p.VideoURL, (*pq.StringArray)(&p.Videos),
		&p.TopReviewComment,
	)
	if err != nil {
		return nil, err
	}
	if p.ReviewCount > 0 {
		p.AvgFood = &avgFood
		p.AvgService = &avgService
		p.AvgVibe = &avgVibe
	}

	// deleted_at — для soft-delete: loadPlaces/share проверяют это поле, чтобы
	// не показывать архивированные места (отдельным лёгким запросом, чтобы не
	// трогать общий SELECT/Scan, идентичный в List).
	_ = r.db.QueryRow(`SELECT deleted_at FROM places WHERE id = $1`, id).Scan(&p.DeletedAt)

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

	photos, err := r.getFeedPhotos(id)
	if err != nil {
		return nil, err
	}
	p.FeedPhotos = photos

	// Q2 — gem-status + attendance + ratings-per-user живут только на детальной
	// карточке места. Каждое поле — отдельный SQL ради ясности; payload growth
	// маленький, а GetByID вызывается редко.
	gemStatus, err := r.getGemStatus(id)
	if err != nil {
		return nil, err
	}
	p.GemStatus = gemStatus

	attendance, err := r.getAttendance(id)
	if err != nil {
		return nil, err
	}
	p.Attendance = attendance

	ratings, err := r.getRatingsPerUser(id)
	if err != nil {
		return nil, err
	}
	p.RatingsPerUser = ratings

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
		// Идентичность совпала с уже существующим местом — отдаём его наверх,
		// чтобы UI предложил «перейти и оставить отзыв», а не глухую ошибку.
		if isUniqueViolation(err) {
			existing, ferr := r.findByIdentity(p.Name, p.Address, p.City)
			if ferr != nil {
				existing = nil
			}
			return nil, &DuplicatePlaceError{Existing: existing}
		}
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

// findByIdentity находит место по тому же ключу, что и uniqueный idx_places_identity
// (LOWER(name) + LOWER(address) + LOWER(city)). Возвращает полную карточку.
func (r *PlaceRepo) findByIdentity(name string, address, city *string) (*model.Place, error) {
	var id int
	err := r.db.QueryRow(`
		SELECT id FROM places
		WHERE LOWER(name) = LOWER($1)
		  AND LOWER(COALESCE(address, '')) = LOWER(COALESCE($2, ''))
		  AND LOWER(COALESCE(city, '')) = LOWER(COALESCE($3, ''))
		LIMIT 1
	`, name, address, city).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
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
		// Переименование в уже существующее место (тот же name+address+city).
		if isUniqueViolation(err) {
			return nil, &DuplicatePlaceError{}
		}
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

// Delete — физическое удаление (каскадит отзывы/фото). Используется только для
// будущего superuser-purge; обычное удаление мест теперь мягкое (SoftDelete).
func (r *PlaceRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM places WHERE id = $1`, id)
	return err
}

// SoftDelete архивирует место: скрывает из листингов/карты/поиска/рандома, но
// отзывы/фото/видео круга сохраняются. No-op, если уже удалено.
func (r *PlaceRepo) SoftDelete(id int) error {
	_, err := r.db.Exec(`UPDATE places SET deleted_at = now(), updated_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

// Restore снимает мягкое удаление (superuser).
func (r *PlaceRepo) Restore(id int) error {
	_, err := r.db.Exec(`UPDATE places SET deleted_at = NULL, updated_at = now() WHERE id = $1`, id)
	return err
}

// CollectUploadPaths собирает все URL загруженных файлов места (обложка места +
// обложки/видео/фото всех его отзывов), чтобы хендлер мог удалить их с диска
// перед каскадным удалением строк. Возвращает значения вида "/uploads/...".
func (r *PlaceRepo) CollectUploadPaths(placeID int) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT image_url FROM places  WHERE id = $1       AND image_url IS NOT NULL AND image_url <> ''
		UNION ALL
		SELECT image_url FROM reviews WHERE place_id = $1 AND image_url IS NOT NULL AND image_url <> ''
		UNION ALL
		SELECT video_url FROM reviews WHERE place_id = $1 AND video_url IS NOT NULL AND video_url <> ''
		UNION ALL
		SELECT rp.url FROM review_photos rp
		  JOIN reviews rv ON rv.id = rp.review_id
		 WHERE rv.place_id = $1 AND rp.url IS NOT NULL AND rp.url <> ''
	`, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		paths = append(paths, p)
	}
	return paths, rows.Err()
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

// CanonicalCity возвращает уже существующее написание города (case-insensitive
// совпадение), чтобы «москва»/«Москва»/«МОСКВА» схлопывались в одну запись на
// полке «По городам» (которая группирует по точной строке p.city). Если такого
// города ещё нет — возвращает вход как есть. Вход ожидается уже TrimSpace'нутым.
func (r *PlaceRepo) CanonicalCity(city string) string {
	if city == "" {
		return ""
	}
	var existing string
	err := r.db.QueryRow(
		`SELECT city FROM places WHERE LOWER(city) = LOWER($1) AND city IS NOT NULL
		 ORDER BY id ASC LIMIT 1`, city,
	).Scan(&existing)
	if err == nil && existing != "" {
		return existing
	}
	return city
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

// Random — случайное место из подходящих под фильтр. Если excludeVisitedBy
// > 0, исключает места, где этот пользователь уже оставил отзыв (B5: «мне
// повезёт» с exclude_visited_by). Возвращает nil, nil если ничего не нашлось.
func (r *PlaceRepo) Random(f PlaceFilter, excludeVisitedBy int) (*model.Place, error) {
	conditions := []string{"p.deleted_at IS NULL"}
	var args []any
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
	if f.IsGem != nil && *f.IsGem {
		conditions = append(conditions, `EXISTS (SELECT 1 FROM reviews rv WHERE rv.place_id = p.id AND rv.is_gem = true)`)
	}
	if excludeVisitedBy > 0 {
		conditions = append(conditions, fmt.Sprintf(
			`NOT EXISTS (SELECT 1 FROM reviews rv JOIN review_authors ra ON ra.review_id = rv.id
			              WHERE rv.place_id = p.id AND ra.user_id = $%d)`, argIdx))
		args = append(args, excludeVisitedBy)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	var id int
	err := r.db.QueryRow(`SELECT p.id FROM places p`+where+` ORDER BY random() LIMIT 1`, args...).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// getFeedPhotos пулит фото по ВСЕМ отзывам места — стопка на ArtifactCard
// представляет все фото круга (DESIGN-DECISIONS §L1: один полароид на place,
// авторы стекаются → их фото тоже стекаются). Сортировка: сперва свежие
// отзывы, внутри отзыва — по position. Лимит 5 для визуальной плотности.
func (r *PlaceRepo) getFeedPhotos(placeID int) ([]model.ReviewPhoto, error) {
	rows, err := r.db.Query(`
		SELECT rp.id, rp.url, rp.position
		FROM review_photos rp
		JOIN reviews rv ON rv.id = rp.review_id
		WHERE rv.place_id = $1
		ORDER BY rv.created_at DESC, rv.id DESC, rp.position ASC, rp.id ASC
		LIMIT 5
	`, placeID)
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

// getGemStatus — Q2: «отметила Аня · 12 марта (+ Серёжа, Миша)». Список собран
// в порядке первой gem-отметки каждого автора; first_marked_at — самая ранняя.
func (r *PlaceRepo) getGemStatus(placeID int) (*model.GemStatus, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.username, u.avatar_url, MIN(rv.created_at) AS first_at
		FROM reviews rv
		JOIN review_authors ra ON ra.review_id = rv.id
		JOIN users u ON u.id = ra.user_id
		WHERE rv.place_id = $1 AND rv.is_gem = true
		GROUP BY u.id, u.username, u.avatar_url
		ORDER BY first_at ASC, u.username ASC
	`, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gs := &model.GemStatus{MarkedBy: []model.Reviewer{}}
	first := true
	for rows.Next() {
		var rev model.Reviewer
		var firstAt time.Time
		if err := rows.Scan(&rev.ID, &rev.Username, &rev.AvatarURL, &firstAt); err != nil {
			return nil, err
		}
		if first {
			gs.FirstMarkedAt = firstAt
			first = false
		}
		gs.MarkedBy = append(gs.MarkedBy, rev)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(gs.MarkedBy) == 0 {
		return nil, nil
	}
	return gs, nil
}

// getAttendance — список «кто был и сколько раз». Под общим тикетом рисуется
// ряд `аватарка · ×N` (N рукописно).
func (r *PlaceRepo) getAttendance(placeID int) ([]model.Attendance, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.username, u.avatar_url, COUNT(rv.id) AS visit_count
		FROM users u
		JOIN review_authors ra ON ra.user_id = u.id
		JOIN reviews rv ON rv.id = ra.review_id
		WHERE rv.place_id = $1
		GROUP BY u.id, u.username, u.avatar_url
		ORDER BY visit_count DESC, u.username ASC
	`, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.Attendance{}
	for rows.Next() {
		var a model.Attendance
		if err := rows.Scan(&a.UserID, &a.Username, &a.AvatarURL, &a.VisitCount); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// getRatingsPerUser — средние оценки каждого пользователя по этому месту.
// Используется только бэком (sort=rating_user:N) и для будущей фичи
// мини-сравнения; UI таблицу не рисует.
func (r *PlaceRepo) getRatingsPerUser(placeID int) ([]model.RatingsPerUser, error) {
	rows, err := r.db.Query(`
		SELECT ra.user_id,
			AVG(rv.food_rating)::numeric(3,1)    AS food,
			AVG(rv.service_rating)::numeric(3,1) AS service,
			AVG(rv.vibe_rating)::numeric(3,1)    AS vibe
		FROM reviews rv
		JOIN review_authors ra ON ra.review_id = rv.id
		WHERE rv.place_id = $1
		GROUP BY ra.user_id
		ORDER BY ra.user_id
	`, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.RatingsPerUser{}
	for rows.Next() {
		var r model.RatingsPerUser
		var food, service, vibe float64
		if err := rows.Scan(&r.UserID, &food, &service, &vibe); err != nil {
			return nil, err
		}
		r.Food = &food
		r.Service = &service
		r.Vibe = &vibe
		out = append(out, r)
	}
	return out, rows.Err()
}

// getReviewersBatch — рецензенты сразу для набора мест (один запрос вместо N).
func (r *PlaceRepo) getReviewersBatch(placeIDs []int) (map[int][]model.Reviewer, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT rv.place_id, u.id, u.username, u.avatar_url
		FROM users u
		JOIN review_authors ra ON ra.user_id = u.id
		JOIN reviews rv ON rv.id = ra.review_id
		WHERE rv.place_id = ANY($1)
		ORDER BY rv.place_id, u.username
	`, pq.Array(placeIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := map[int][]model.Reviewer{}
	for rows.Next() {
		var pid int
		var rev model.Reviewer
		if err := rows.Scan(&pid, &rev.ID, &rev.Username, &rev.AvatarURL); err != nil {
			return nil, err
		}
		m[pid] = append(m[pid], rev)
	}
	return m, rows.Err()
}

// getFeedPhotosBatch — до 5 фото на место сразу для набора мест. Per-place
// LIMIT 5 реализован оконной функцией ROW_NUMBER (тот же порядок, что в
// getFeedPhotos): свежие отзывы → внутри по position.
func (r *PlaceRepo) getFeedPhotosBatch(placeIDs []int) (map[int][]model.ReviewPhoto, error) {
	rows, err := r.db.Query(`
		SELECT place_id, id, url, position FROM (
			SELECT rv.place_id, rp.id, rp.url, rp.position,
				ROW_NUMBER() OVER (
					PARTITION BY rv.place_id
					ORDER BY rv.created_at DESC, rv.id DESC, rp.position ASC, rp.id ASC
				) AS rn
			FROM review_photos rp
			JOIN reviews rv ON rv.id = rp.review_id
			WHERE rv.place_id = ANY($1)
		) t
		WHERE rn <= 5
		ORDER BY place_id
	`, pq.Array(placeIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := map[int][]model.ReviewPhoto{}
	for rows.Next() {
		var pid int
		var p model.ReviewPhoto
		if err := rows.Scan(&pid, &p.ID, &p.URL, &p.Position); err != nil {
			return nil, err
		}
		m[pid] = append(m[pid], p)
	}
	return m, rows.Err()
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
