package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

type ReviewHandler struct {
	reviewRepo   *repository.ReviewRepo
	wishlistRepo *repository.WishlistRepo
}

func NewReviewHandler(reviewRepo *repository.ReviewRepo, wishlistRepo *repository.WishlistRepo) *ReviewHandler {
	return &ReviewHandler{reviewRepo: reviewRepo, wishlistRepo: wishlistRepo}
}

type createReviewRequest struct {
	FoodRating    float64 `json:"food_rating"`
	ServiceRating float64 `json:"service_rating"`
	VibeRating    float64 `json:"vibe_rating"`
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

	// backend.md §Wishlist: при создании review у соавторов, у которых это
	// место было в wishlist, ставим struck=true. Best-effort — ошибки не
	// фейлят основной запрос; review уже создан.
	if h.wishlistRepo != nil {
		for _, uid := range authorIDs {
			_, _ = h.wishlistRepo.MarkStruck(uid, placeID)
		}
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

func validateRatings(food, service, vibe float64) error {
	for _, v := range []struct {
		name string
		val  float64
	}{
		{"food_rating", food},
		{"service_rating", service},
		{"vibe_rating", vibe},
	} {
		if v.val < 0 || v.val > 10 {
			return fmt.Errorf("%s must be between 0 and 10", v.name)
		}
	}
	return nil
}

const maxPhotosPerReview = 5

// UploadImage — legacy single-photo endpoint. Сохранён для backwards-compat,
// но теперь добавляет фото в review_photos через AddPhoto (то же поведение,
// что и новый /photos endpoint при загрузке одного файла).
func (h *ReviewHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
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
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only upload images for your own reviews"})
		return
	}

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

	url, perr := h.processAndStorePhoto(file, header)
	if perr != nil {
		writeJSON(w, perr.statusCode(), map[string]string{"error": perr.Error()})
		return
	}

	count, err := h.reviewRepo.CountPhotos(reviewID)
	if err != nil {
		removeUploaded(url)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to count photos"})
		return
	}
	if count >= maxPhotosPerReview {
		removeUploaded(url)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("review already has %d photos", maxPhotosPerReview)})
		return
	}

	if _, err := h.reviewRepo.AddPhoto(reviewID, url); err != nil {
		removeUploaded(url)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save photo"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"image_url": url})
}

// UploadPhotos принимает до maxPhotosPerReview файлов в поле "photos" одним
// multipart-запросом и добавляет их в стопку review.photos в порядке прихода.
// Если total после загрузки превысит лимит — отвергает целиком (atomic enough).
func (h *ReviewHandler) UploadPhotos(w http.ResponseWriter, r *http.Request) {
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
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only upload photos for your own reviews"})
		return
	}

	// 5MB на файл × 5 = 25MB лимит формы.
	if err := r.ParseMultipartForm(25 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "form too large (max 25MB total)"})
		return
	}

	files := r.MultipartForm.File["photos"]
	if len(files) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "photos field required"})
		return
	}

	count, err := h.reviewRepo.CountPhotos(reviewID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to count photos"})
		return
	}
	if count+len(files) > maxPhotosPerReview {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("max %d photos per review (have %d, tried to add %d)", maxPhotosPerReview, count, len(files)),
		})
		return
	}

	created := make([]model.ReviewPhoto, 0, len(files))
	urls := make([]string, 0, len(files))

	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			rollbackPhotos(h, urls, created)
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to read file"})
			return
		}
		url, perr := h.processAndStorePhoto(file, header)
		file.Close()
		if perr != nil {
			rollbackPhotos(h, urls, created)
			writeJSON(w, perr.statusCode(), map[string]string{"error": perr.Error()})
			return
		}
		photo, dberr := h.reviewRepo.AddPhoto(reviewID, url)
		if dberr != nil {
			removeUploaded(url)
			rollbackPhotos(h, urls, created)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save photo"})
			return
		}
		urls = append(urls, url)
		created = append(created, *photo)
	}

	writeJSON(w, http.StatusCreated, map[string]any{"photos": created})
}

// DeletePhoto удаляет одно фото из стопки. Author-only через review.
func (h *ReviewHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	photoID, err := strconv.Atoi(chi.URLParam(r, "pid"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid photo id"})
		return
	}

	photo, reviewID, err := h.reviewRepo.GetPhoto(photoID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "photo not found"})
		return
	}

	isAuthor, err := h.reviewRepo.IsAuthor(reviewID, userID)
	if err != nil || !isAuthor {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only delete photos from your own reviews"})
		return
	}

	if err := h.reviewRepo.DeletePhoto(photoID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete photo"})
		return
	}

	removeUploaded(photo.URL)
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// uploadErr — обёртка для проброса HTTP-статуса из processAndStorePhoto наверх.
type uploadErr struct {
	msg  string
	code int
}

func (e *uploadErr) Error() string     { return e.msg }
func (e *uploadErr) statusCode() int   { return e.code }
func badInput(msg string) *uploadErr   { return &uploadErr{msg: msg, code: http.StatusBadRequest} }
func internal(msg string) *uploadErr   { return &uploadErr{msg: msg, code: http.StatusInternalServerError} }

// processAndStorePhoto проверяет тип, прогоняет через imageutil (auto-orient
// + JPEG q=75) и пишет результат в /uploads. Возвращает публичный URL.
func (h *ReviewHandler) processAndStorePhoto(file io.Reader, header *multipart.FileHeader) (string, *uploadErr) {
	ct := header.Header.Get("Content-Type")
	if _, ok := allowedImageTypes[ct]; !ok {
		return "", badInput("only JPEG, PNG and WebP images are allowed")
	}

	if err := os.MkdirAll("uploads", 0o755); err != nil {
		return "", internal("failed to create uploads directory")
	}

	filename := fmt.Sprintf("%s.jpg", uuid.New().String())
	dstPath := filepath.Join("uploads", filename)
	if err := imageutil.Process(file, dstPath); err != nil {
		return "", internal("failed to process image")
	}
	return "/uploads/" + filename, nil
}

func removeUploaded(url string) {
	if url == "" {
		return
	}
	// "/uploads/foo.jpg" → "uploads/foo.jpg"
	rel := strings.TrimPrefix(url, "/")
	_ = os.Remove(rel)
}

// rollbackPhotos чистит уже сохранённые в этом запросе фото при middle-failure.
func rollbackPhotos(h *ReviewHandler, urls []string, created []model.ReviewPhoto) {
	for _, p := range created {
		_ = h.reviewRepo.DeletePhoto(p.ID)
	}
	for _, u := range urls {
		removeUploaded(u)
	}
}

var allowedVideoTypes = map[string]string{
	"video/mp4":  ".mp4",
	"video/webm": ".webm",
}

func (h *ReviewHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
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
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "you can only upload videos for your own reviews"})
		return
	}

	// 20MB max for video
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file too large (max 20MB)"})
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "video field required"})
		return
	}
	defer file.Close()

	ct := header.Header.Get("Content-Type")
	ext, ok := allowedVideoTypes[ct]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "only MP4 and WebM videos are allowed"})
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
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save video"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to write video"})
		return
	}

	videoURL := "/uploads/" + filename
	if err := h.reviewRepo.UpdateVideoURL(reviewID, videoURL); err != nil {
		os.Remove(filepath.Join(uploadsDir, filename))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update review video"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"video_url": videoURL})
}
