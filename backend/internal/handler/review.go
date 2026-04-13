package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aeva-eat/backend/internal/middleware"
	"github.com/aeva-eat/backend/internal/model"
	"github.com/aeva-eat/backend/internal/repository"
	"github.com/go-chi/chi/v5"
)

type ReviewHandler struct {
	reviewRepo *repository.ReviewRepo
}

func NewReviewHandler(reviewRepo *repository.ReviewRepo) *ReviewHandler {
	return &ReviewHandler{reviewRepo: reviewRepo}
}

type createReviewRequest struct {
	FoodRating    int     `json:"food_rating"`
	ServiceRating int     `json:"service_rating"`
	VibeRating    int     `json:"vibe_rating"`
	IsGem         bool    `json:"is_gem"`
	Comment       *string `json:"comment,omitempty"`
	VisitedAt     *string `json:"visited_at,omitempty"`
	AuthorIDs     []int   `json:"author_ids,omitempty"`
}

func (h *ReviewHandler) ListByPlace(w http.ResponseWriter, r *http.Request) {
	placeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid place id"})
		return
	}

	reviews, err := h.reviewRepo.ListByPlace(placeID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list reviews"})
		return
	}
	if reviews == nil {
		reviews = []model.Review{}
	}
	writeJSON(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}

	reviews, err := h.reviewRepo.ListByUser(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list reviews"})
		return
	}
	if reviews == nil {
		reviews = []model.Review{}
	}
	writeJSON(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	placeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid place id"})
		return
	}

	var req createReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := validateRatings(req.FoodRating, req.ServiceRating, req.VibeRating); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	rv := &model.Review{
		PlaceID:       placeID,
		FoodRating:    req.FoodRating,
		ServiceRating: req.ServiceRating,
		VibeRating:    req.VibeRating,
		IsGem:         req.IsGem,
		Comment:       req.Comment,
		VisitedAt:     req.VisitedAt,
	}

	authorIDs := req.AuthorIDs
	hasCurrentUser := false
	for _, id := range authorIDs {
		if id == userID {
			hasCurrentUser = true
			break
		}
	}
	if !hasCurrentUser {
		authorIDs = append([]int{userID}, authorIDs...)
	}

	created, err := h.reviewRepo.Create(rv, authorIDs)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create review"})
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	reviewID, err := strconv.Atoi(chi.URLParam(r, "rid"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid review id"})
		return
	}

	isAuthor, err := h.reviewRepo.IsAuthor(reviewID, userID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "review not found"})
		return
	}
	if !isAuthor {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only edit your own reviews"})
		return
	}

	var req createReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := validateRatings(req.FoodRating, req.ServiceRating, req.VibeRating); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	rv := &model.Review{
		ID:            reviewID,
		FoodRating:    req.FoodRating,
		ServiceRating: req.ServiceRating,
		VibeRating:    req.VibeRating,
		IsGem:         req.IsGem,
		Comment:       req.Comment,
		VisitedAt:     req.VisitedAt,
	}

	updated, err := h.reviewRepo.Update(rv)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update review"})
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	reviewID, err := strconv.Atoi(chi.URLParam(r, "rid"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid review id"})
		return
	}

	isAuthor, err := h.reviewRepo.IsAuthor(reviewID, userID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "review not found"})
		return
	}
	if !isAuthor {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only delete your own reviews"})
		return
	}

	if err := h.reviewRepo.Delete(reviewID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete review"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func validateRatings(food, service, vibe int) error {
	for _, v := range []struct {
		name string
		val  int
	}{
		{"food_rating", food},
		{"service_rating", service},
		{"vibe_rating", vibe},
	} {
		if v.val < 1 || v.val > 10 {
			return fmt.Errorf("%s must be between 1 and 10", v.name)
		}
	}
	return nil
}
