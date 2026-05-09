package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aeva-eat/backend/internal/imageutil"
	"github.com/aeva-eat/backend/internal/middleware"
	"github.com/aeva-eat/backend/internal/model"
	"github.com/aeva-eat/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PlaceHandler struct {
	placeRepo *repository.PlaceRepo
}

func NewPlaceHandler(placeRepo *repository.PlaceRepo) *PlaceHandler {
	return &PlaceHandler{placeRepo: placeRepo}
}

type createPlaceRequest struct {
	Name          string   `json:"name"`
	Address       *string  `json:"address,omitempty"`
	City          *string  `json:"city,omitempty"`
	Lat           *float64 `json:"lat,omitempty"`
	Lng           *float64 `json:"lng,omitempty"`
	CuisineTypeID *int     `json:"cuisine_type_id,omitempty"`
	Website       *string  `json:"website,omitempty"`
	CategoryIDs   []int    `json:"category_ids,omitempty"`
}

func parsePlaceListFilter(q url.Values) repository.PlaceFilter {
	filter := repository.PlaceFilter{
		City:   q.Get("city"),
		Search: q.Get("search"),
		Sort:   q.Get("sort"),
	}

	if v := q.Get("cuisine_type_id"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if id, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				filter.CuisineTypeIDs = append(filter.CuisineTypeIDs, id)
			}
		}
	}
	if v := q.Get("category_id"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if id, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				filter.CategoryIDs = append(filter.CategoryIDs, id)
			}
		}
	}
	if v := q.Get("min_rating"); v != "" {
		if rating, err := strconv.ParseFloat(v, 64); err == nil {
			filter.MinRating = rating
		}
	}
	if v := q.Get("is_gem"); v == "true" {
		isGem := true
		filter.IsGem = &isGem
	}
	if v := q.Get("attended_by"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if id, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				filter.AttendedBy = append(filter.AttendedBy, id)
			}
		}
	}
	if v := q.Get("visit_from"); v != "" {
		filter.VisitFrom = v
	}
	if v := q.Get("visit_to"); v != "" {
		filter.VisitTo = v
	}
	// sort=rating_user:<id> — раскладываем на canonical "rating_user" + id.
	if strings.HasPrefix(filter.Sort, "rating_user:") {
		if id, err := strconv.Atoi(strings.TrimPrefix(filter.Sort, "rating_user:")); err == nil && id > 0 {
			filter.Sort = "rating_user"
			filter.SortRatingUserID = id
		} else {
			// невалидный — фоллбэк на дефолт
			filter.Sort = ""
		}
	}

	// Pagination
	limit := 20
	if v := q.Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil {
			if l == 0 {
				limit = 0
			} else if l > 0 && l <= 100 {
				limit = l
			}
		}
	}
	page := 1
	if v := q.Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}
	filter.Limit = limit
	filter.Offset = (page - 1) * limit
	return filter
}

func (h *PlaceHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := parsePlaceListFilter(r.URL.Query())

	result, err := h.placeRepo.List(filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list places"})
		return
	}
	if result.Places == nil {
		result.Places = []model.Place{}
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *PlaceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid place id"})
		return
	}

	place, err := h.placeRepo.GetByID(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "place not found"})
		return
	}
	writeJSON(w, http.StatusOK, place)
}

func (h *PlaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req createPlaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	place := &model.Place{
		Name:          req.Name,
		Address:       req.Address,
		City:          req.City,
		Lat:           req.Lat,
		Lng:           req.Lng,
		CuisineTypeID: req.CuisineTypeID,
		Website:       req.Website,
		CreatedBy:     &userID,
	}

	created, err := h.placeRepo.Create(place, req.CategoryIDs)
	if err != nil {
		if strings.Contains(err.Error(), "idx_places_name_city") {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "Такое заведение уже существует в этом городе"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create place"})
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *PlaceHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid place id"})
		return
	}

	ownerID, err := h.placeRepo.GetOwnerID(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "place not found"})
		return
	}
	if ownerID != userID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only edit your own places"})
		return
	}

	var req createPlaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	// Preserve existing image_url
	existing, _ := h.placeRepo.GetByID(id)
	var imageURL *string
	if existing != nil {
		imageURL = existing.ImageURL
	}

	place := &model.Place{
		ID:            id,
		Name:          req.Name,
		Address:       req.Address,
		City:          req.City,
		Lat:           req.Lat,
		Lng:           req.Lng,
		CuisineTypeID: req.CuisineTypeID,
		Website:       req.Website,
		ImageURL:      imageURL,
	}

	updated, err := h.placeRepo.Update(place, req.CategoryIDs)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update place"})
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *PlaceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid place id"})
		return
	}

	ownerID, err := h.placeRepo.GetOwnerID(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "place not found"})
		return
	}
	if ownerID != userID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only delete your own places"})
		return
	}

	if err := h.placeRepo.Delete(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete place"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *PlaceHandler) ListCities(w http.ResponseWriter, r *http.Request) {
	cities, err := h.placeRepo.ListCities()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list cities"})
		return
	}
	if cities == nil {
		cities = []string{}
	}
	writeJSON(w, http.StatusOK, cities)
}

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// Random — GET /api/random?city=&cuisine_type_id=&is_gem=true&exclude_visited_by=me
// (NEXT.md §B5). exclude_visited_by=me использует id текущего пользователя
// из JWT, иначе принимает явный user_id. 404 если ничего не подошло.
func (h *PlaceHandler) Random(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := parsePlaceListFilter(q)
	filter.Limit = 0
	filter.Offset = 0

	var excludeUser int
	if exclude := q.Get("exclude_visited_by"); exclude != "" {
		if exclude == "me" {
			if uid, ok := middleware.GetUserID(r); ok {
				excludeUser = uid
			}
		} else if id, err := strconv.Atoi(exclude); err == nil {
			excludeUser = id
		}
	}

	place, err := h.placeRepo.Random(filter, excludeUser)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to pick random place"})
		return
	}
	if place == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no place matches filters"})
		return
	}
	writeJSON(w, http.StatusOK, place)
}

func (h *PlaceHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid place id"})
		return
	}

	ownerID, err := h.placeRepo.GetOwnerID(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "place not found"})
		return
	}
	if ownerID != userID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only upload images for your own places"})
		return
	}

	// 5MB max
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file too large (max 5MB)"})
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "image field required"})
		return
	}
	defer file.Close()

	ct := header.Header.Get("Content-Type")
	if _, ok := allowedImageTypes[ct]; !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "only JPEG, PNG and WebP images are allowed"})
		return
	}

	// Always save as .jpg after processing
	filename := fmt.Sprintf("%s.jpg", uuid.New().String())

	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create uploads directory"})
		return
	}

	dstPath := filepath.Join(uploadsDir, filename)
	if err := imageutil.Process(file, dstPath); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to process image"})
		return
	}

	imageURL := "/uploads/" + filename

	// Get old image before updating
	oldPlace, _ := h.placeRepo.GetByID(id)

	if err := h.placeRepo.UpdateImageURL(id, imageURL); err != nil {
		os.Remove(dstPath)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update image"})
		return
	}

	// Remove old image if present
	if oldPlace != nil && oldPlace.ImageURL != nil && strings.HasPrefix(*oldPlace.ImageURL, "/uploads/") {
		old := filepath.Join(uploadsDir, filepath.Base(*oldPlace.ImageURL))
		os.Remove(old)
	}

	writeJSON(w, http.StatusOK, map[string]string{"image_url": imageURL})
}
