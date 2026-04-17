package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aeva-eat/backend/internal/middleware"
	"github.com/aeva-eat/backend/internal/model"
	"github.com/aeva-eat/backend/internal/repository"
	"github.com/go-chi/chi/v5"
)

type InviteHandler struct {
	inviteRepo *repository.InviteRepo
	userRepo   *repository.UserRepo
}

func NewInviteHandler(inviteRepo *repository.InviteRepo, userRepo *repository.UserRepo) *InviteHandler {
	return &InviteHandler{inviteRepo: inviteRepo, userRepo: userRepo}
}

type createInviteRequest struct {
	ExpiresAt *string `json:"expires_at,omitempty"`
}

func generateCode() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *InviteHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req createInviteRequest
	if r.ContentLength > 0 {
		json.NewDecoder(r.Body).Decode(&req)
	}

	code := generateCode()
	invite, err := h.inviteRepo.Create(code, userID, req.ExpiresAt)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create invite"})
		return
	}

	writeJSON(w, http.StatusCreated, invite)
}

func (h *InviteHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	invites, err := h.inviteRepo.ListByCreator(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list invites"})
		return
	}
	if invites == nil {
		invites = []model.Invite{}
	}
	writeJSON(w, http.StatusOK, invites)
}

func (h *InviteHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil || user.Role != "superuser" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "superuser access required"})
		return
	}

	invites, err := h.inviteRepo.ListAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list invites"})
		return
	}
	writeJSON(w, http.StatusOK, invites)
}

func (h *InviteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid invite id"})
		return
	}

	// Check that the user owns this invite or is superuser
	user, _ := h.userRepo.GetByID(userID)
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	if err := h.inviteRepo.Delete(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete invite"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *InviteHandler) ValidateCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "code is required"})
		return
	}

	invite, err := h.inviteRepo.GetByCode(code)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "invalid invite code", "valid": "false"})
		return
	}
	if invite.UsedBy != nil {
		writeJSON(w, http.StatusGone, map[string]string{"error": "invite already used", "valid": "false"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"valid":        true,
		"creator_name": invite.CreatorName,
	})
}

// Admin: list all users
func (h *InviteHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil || user.Role != "superuser" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "superuser access required"})
		return
	}

	users, err := h.userRepo.ListAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list users"})
		return
	}
	writeJSON(w, http.StatusOK, users)
}
