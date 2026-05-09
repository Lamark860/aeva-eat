package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/aeva-eat/backend/internal/middleware"
	"github.com/aeva-eat/backend/internal/model"
	"github.com/aeva-eat/backend/internal/repository"
	"github.com/go-chi/chi/v5"
)

type NoteHandler struct {
	notes *repository.NoteRepo
	feed  *repository.FeedEventsRepo
}

func NewNoteHandler(notes *repository.NoteRepo, feed *repository.FeedEventsRepo) *NoteHandler {
	return &NoteHandler{notes: notes, feed: feed}
}

type noteRequest struct {
	Text       string  `json:"text"`
	PlaceID    *int    `json:"place_id,omitempty"`
	City       *string `json:"city,omitempty"`
	PaperColor *string `json:"paper_color,omitempty"`
	TapeColor  *string `json:"tape_color,omitempty"`
}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
	authorParam := r.URL.Query().Get("author_id")
	var notes []model.Note
	var err error
	if authorParam != "" {
		authorID, perr := strconv.Atoi(authorParam)
		if perr != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid author_id"})
			return
		}
		notes, err = h.notes.ListByAuthor(authorID)
	} else {
		notes, err = h.notes.List()
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list notes"})
		return
	}
	if notes == nil {
		notes = []model.Note{}
	}
	writeJSON(w, http.StatusOK, notes)
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req noteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	text := strings.TrimSpace(req.Text)
	if text == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "text is required"})
		return
	}
	if len(text) > 2000 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "text too long (max 2000)"})
		return
	}

	n := &model.Note{
		AuthorID:   userID,
		Text:       text,
		PlaceID:    req.PlaceID,
		City:       req.City,
		PaperColor: req.PaperColor,
		TapeColor:  req.TapeColor,
	}
	created, err := h.notes.Create(n)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create note"})
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid note id"})
		return
	}

	isAuthor, err := h.notes.IsAuthor(id, userID)
	if err != nil || !isAuthor {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only edit your own notes"})
		return
	}

	var req noteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	text := strings.TrimSpace(req.Text)
	if text == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "text is required"})
		return
	}

	n := &model.Note{
		ID:         id,
		Text:       text,
		PlaceID:    req.PlaceID,
		City:       req.City,
		PaperColor: req.PaperColor,
		TapeColor:  req.TapeColor,
	}
	updated, err := h.notes.Update(n)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update note"})
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid note id"})
		return
	}

	isAuthor, err := h.notes.IsAuthor(id, userID)
	if err != nil || !isAuthor {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only delete your own notes"})
		return
	}

	if err := h.notes.Delete(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete note"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *NoteHandler) Strike(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid note id"})
		return
	}

	isAuthor, err := h.notes.IsAuthor(id, userID)
	if err != nil || !isAuthor {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only strike your own notes"})
		return
	}

	updated, err := h.notes.SetStruck(id, true)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to strike note"})
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// Feed — единая лента всех событий круга. На MVP без cursor-пагинации.
func (h *NoteHandler) Feed(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	events, err := h.feed.List(limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list feed"})
		return
	}
	if events == nil {
		events = []model.FeedEvent{}
	}
	writeJSON(w, http.StatusOK, events)
}

// Weeks — агрегат по неделям для CollapsedStrip на Доске (backend.md
// §Лента/Доска: GET /api/feed/weeks).
func (h *NoteHandler) Weeks(w http.ResponseWriter, r *http.Request) {
	weeks, err := h.feed.Weeks()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list weeks"})
		return
	}
	writeJSON(w, http.StatusOK, weeks)
}

// UnreadCount — точка на BottomTabBar (NEXT.md §C2).
func (h *NoteHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	n, err := h.feed.UnreadCount(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to count unread"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": n})
}

// MarkSeen — клиент дёргает при открытии Доски, чтобы сбросить точку.
func (h *NoteHandler) MarkSeen(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if err := h.feed.MarkSeen(userID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to mark seen"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
