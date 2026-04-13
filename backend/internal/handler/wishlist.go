package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aeva-eat/backend/internal/middleware"
	"github.com/aeva-eat/backend/internal/model"
	"github.com/aeva-eat/backend/internal/repository"
	"github.com/go-chi/chi/v5"
)

type WishlistHandler struct {
	wishlistRepo *repository.WishlistRepo
}

func NewWishlistHandler(wishlistRepo *repository.WishlistRepo) *WishlistHandler {
	return &WishlistHandler{wishlistRepo: wishlistRepo}
}

func (h *WishlistHandler) Add(w http.ResponseWriter, r *http.Request) {
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

	if err := h.wishlistRepo.Add(userID, placeID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to add to wishlist"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "added"})
}

func (h *WishlistHandler) Remove(w http.ResponseWriter, r *http.Request) {
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

	if err := h.wishlistRepo.Remove(userID, placeID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to remove from wishlist"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

func (h *WishlistHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	places, err := h.wishlistRepo.ListByUser(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list wishlist"})
		return
	}
	if places == nil {
		places = []model.Place{}
	}
	writeJSON(w, http.StatusOK, places)
}

func (h *WishlistHandler) MyIDs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	ids, err := h.wishlistRepo.GetUserWishlistIDs(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get wishlist ids"})
		return
	}
	if ids == nil {
		ids = []int{}
	}
	writeJSON(w, http.StatusOK, ids)
}

// Custom wishlist (free-text entries)

type customWishlistRequest struct {
	Name string  `json:"name"`
	Note *string `json:"note,omitempty"`
}

func (h *WishlistHandler) AddCustom(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req customWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	wc, err := h.wishlistRepo.AddCustom(userID, req.Name, req.Note)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to add custom wishlist"})
		return
	}
	writeJSON(w, http.StatusCreated, wc)
}

func (h *WishlistHandler) ListCustom(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	items, err := h.wishlistRepo.ListCustom(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list custom wishlist"})
		return
	}
	if items == nil {
		items = []model.WishlistCustom{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *WishlistHandler) DeleteCustom(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if err := h.wishlistRepo.DeleteCustom(userID, id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete custom wishlist"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
