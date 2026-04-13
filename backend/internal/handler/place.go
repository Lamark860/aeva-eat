package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

func (h *PlaceHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
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

	places, err := h.placeRepo.List(filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list places"})
		return
	}
	if places == nil {
		places = []model.Place{}
	}
	writeJSON(w, http.StatusOK, places)
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

	place := &model.Place{
		ID:            id,
		Name:          req.Name,
		Address:       req.Address,
		City:          req.City,
		Lat:           req.Lat,
		Lng:           req.Lng,
		CuisineTypeID: req.CuisineTypeID,
		Website:       req.Website,
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
	ext, ok := allowedImageTypes[ct]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "only JPEG, PNG and WebP images are allowed"})
		return
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create uploads directory"})
		return
	}

	dst, err := os.Create(filepath.Join(uploadsDir, filename))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save image"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to write image"})
		return
	}

	imageURL := "/uploads/" + filename

	// Get old image before updating
	oldPlace, _ := h.placeRepo.GetByID(id)

	if err := h.placeRepo.UpdateImageURL(id, imageURL); err != nil {
		os.Remove(filepath.Join(uploadsDir, filename))
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
