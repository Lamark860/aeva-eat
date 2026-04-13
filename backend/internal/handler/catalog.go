package handler

import (
	"net/http"

	"github.com/aeva-eat/backend/internal/repository"
)

type CatalogHandler struct {
	catalogRepo *repository.CatalogRepo
}

func NewCatalogHandler(catalogRepo *repository.CatalogRepo) *CatalogHandler {
	return &CatalogHandler{catalogRepo: catalogRepo}
}

func (h *CatalogHandler) ListCuisineTypes(w http.ResponseWriter, r *http.Request) {
	items, err := h.catalogRepo.ListCuisineTypes()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list cuisine types"})
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *CatalogHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	items, err := h.catalogRepo.ListCategories()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list categories"})
		return
	}
	writeJSON(w, http.StatusOK, items)
}
