package handler

import (
	"net/http"
	"strconv"

	"github.com/aeva-eat/backend/internal/model"
	"github.com/aeva-eat/backend/internal/repository"
	"github.com/go-chi/chi/v5"
)

// AggregateHandler — роуты для путеводных страниц: /cities, /users/:id, /gems.
type AggregateHandler struct {
	agg   *repository.AggregateRepo
	place *repository.PlaceRepo
}

func NewAggregateHandler(agg *repository.AggregateRepo, place *repository.PlaceRepo) *AggregateHandler {
	return &AggregateHandler{agg: agg, place: place}
}

// loadPlaces — батч-чтение мест по id'шникам через GetByID. Сохраняет
// порядок входного списка. Тихо пропускает удалённые / недоступные id.
func (h *AggregateHandler) loadPlaces(ids []int) []model.Place {
	out := make([]model.Place, 0, len(ids))
	for _, id := range ids {
		p, err := h.place.GetByID(id)
		if err != nil || p == nil {
			continue
		}
		out = append(out, *p)
	}
	return out
}

// GET /api/cities — все города круга с агрегатами.
func (h *AggregateHandler) ListCities(w http.ResponseWriter, r *http.Request) {
	cities, err := h.agg.Cities()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list cities"})
		return
	}
	writeJSON(w, http.StatusOK, cities)
}

// GET /api/cities/{name}
func (h *AggregateHandler) GetCity(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	c, err := h.agg.City(name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load city"})
		return
	}
	if c == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "city not found"})
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// GET /api/cities/{name}/places — делегирует в PlaceRepo.List с City-фильтром.
func (h *AggregateHandler) ListCityPlaces(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	filter := parsePlaceListFilter(r.URL.Query())
	filter.City = name
	res, err := h.place.List(filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list places"})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// GET /api/cities/{name}/gems
func (h *AggregateHandler) ListCityGems(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	filter := parsePlaceListFilter(r.URL.Query())
	filter.City = name
	t := true
	filter.IsGem = &t
	res, err := h.place.List(filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list gems"})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// GET /api/users — список всех пользователей круга для полки «По друзьям».
func (h *AggregateHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.agg.UsersList()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list users"})
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// GET /api/users/{id} — публичный профиль.
func (h *AggregateHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}
	profile, err := h.agg.UserProfile(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load user"})
		return
	}
	if profile == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

// GET /api/users/{id}/places
func (h *AggregateHandler) ListUserPlaces(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}
	ids, err := h.agg.UserPlaceIDs(id, false)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list places"})
		return
	}
	writeJSON(w, http.StatusOK, h.loadPlaces(ids))
}

// GET /api/users/{id}/gems
func (h *AggregateHandler) ListUserGems(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}
	ids, err := h.agg.UserPlaceIDs(id, true)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list gems"})
		return
	}
	writeJSON(w, http.StatusOK, h.loadPlaces(ids))
}

// GET /api/users/{id}/cities
func (h *AggregateHandler) ListUserCities(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}
	cities, err := h.agg.UserCities(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list cities"})
		return
	}
	writeJSON(w, http.StatusOK, cities)
}

// GET /api/gems — hub жемчужин.
func (h *AggregateHandler) Gems(w http.ResponseWriter, r *http.Request) {
	ids, err := h.agg.GemPlaceIDs()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list gems"})
		return
	}
	byCity, err := h.agg.GemsByCity()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list gems by city"})
		return
	}
	byUser, err := h.agg.GemsByUser()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list gems by user"})
		return
	}

	hub := model.GemsHub{
		Places: h.loadPlaces(ids),
		Total:  len(ids),
		ByCity: byCity,
		ByUser: byUser,
	}
	writeJSON(w, http.StatusOK, hub)
}
