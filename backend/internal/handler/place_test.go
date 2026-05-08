package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/aeva-eat/backend/internal/repository"
	"github.com/go-chi/chi/v5"
)

func TestParsePlaceListFilter_Defaults(t *testing.T) {
	f := parsePlaceListFilter(url.Values{})
	if f.City != "" || f.Search != "" || f.Sort != "" {
		t.Errorf("expected empty string fields, got %+v", f)
	}
	if f.MinRating != 0 {
		t.Errorf("expected min_rating 0, got %v", f.MinRating)
	}
	if f.IsGem != nil {
		t.Errorf("expected IsGem nil, got %v", *f.IsGem)
	}
	if len(f.CuisineTypeIDs) != 0 || len(f.CategoryIDs) != 0 {
		t.Errorf("expected empty id slices, got %+v / %+v", f.CuisineTypeIDs, f.CategoryIDs)
	}
	if f.Limit != 20 {
		t.Errorf("expected default limit 20, got %d", f.Limit)
	}
	if f.Offset != 0 {
		t.Errorf("expected default offset 0, got %d", f.Offset)
	}
}

func TestParsePlaceListFilter_Pagination(t *testing.T) {
	tests := []struct {
		name       string
		limit      string
		page       string
		wantLimit  int
		wantOffset int
	}{
		{"limit=0 (all)", "0", "", 0, 0},
		{"limit=10 page=3", "10", "3", 10, 20},
		{"limit=100 max allowed", "100", "1", 100, 0},
		{"limit=101 over cap → default", "101", "1", 20, 0},
		{"limit negative → default", "-5", "1", 20, 0},
		{"limit non-numeric → default", "abc", "1", 20, 0},
		{"page=0 → default page=1", "20", "0", 20, 0},
		{"page negative → default page=1", "20", "-3", 20, 0},
		{"page non-numeric → default page=1", "20", "abc", 20, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := url.Values{}
			if tt.limit != "" {
				q.Set("limit", tt.limit)
			}
			if tt.page != "" {
				q.Set("page", tt.page)
			}
			f := parsePlaceListFilter(q)
			if f.Limit != tt.wantLimit {
				t.Errorf("limit: want %d, got %d", tt.wantLimit, f.Limit)
			}
			if f.Offset != tt.wantOffset {
				t.Errorf("offset: want %d, got %d", tt.wantOffset, f.Offset)
			}
		})
	}
}

func TestParsePlaceListFilter_MultiIDFilters(t *testing.T) {
	tests := []struct {
		name        string
		cuisine     string
		category    string
		wantCuisine []int
		wantCat     []int
	}{
		{"single ids", "1", "2", []int{1}, []int{2}},
		{"comma list", "1,2,3", "4,5", []int{1, 2, 3}, []int{4, 5}},
		{"with spaces", " 1 , 2 ,3", "", []int{1, 2, 3}, nil},
		{"skip non-numeric", "1,foo,2", "bar,3", []int{1, 2}, []int{3}},
		{"all garbage → empty", "foo,bar", "baz", nil, nil},
		{"empty string", "", "", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := url.Values{}
			if tt.cuisine != "" {
				q.Set("cuisine_type_id", tt.cuisine)
			}
			if tt.category != "" {
				q.Set("category_id", tt.category)
			}
			f := parsePlaceListFilter(q)
			if !reflect.DeepEqual(f.CuisineTypeIDs, tt.wantCuisine) {
				t.Errorf("cuisine: want %v, got %v", tt.wantCuisine, f.CuisineTypeIDs)
			}
			if !reflect.DeepEqual(f.CategoryIDs, tt.wantCat) {
				t.Errorf("category: want %v, got %v", tt.wantCat, f.CategoryIDs)
			}
		})
	}
}

func TestParsePlaceListFilter_IsGem(t *testing.T) {
	tests := []struct {
		val  string
		want *bool
	}{
		{"true", boolPtr(true)},
		{"false", nil},
		{"1", nil}, // only literal "true" is honored
		{"", nil},
		{"TRUE", nil}, // case-sensitive by design
	}
	for _, tt := range tests {
		t.Run("is_gem="+tt.val, func(t *testing.T) {
			q := url.Values{}
			if tt.val != "" {
				q.Set("is_gem", tt.val)
			}
			f := parsePlaceListFilter(q)
			if (f.IsGem == nil) != (tt.want == nil) {
				t.Fatalf("nilness mismatch: got=%v want=%v", f.IsGem, tt.want)
			}
			if f.IsGem != nil && *f.IsGem != *tt.want {
				t.Errorf("value mismatch: got=%v want=%v", *f.IsGem, *tt.want)
			}
		})
	}
}

func TestParsePlaceListFilter_MinRating(t *testing.T) {
	tests := []struct {
		val  string
		want float64
	}{
		{"7.5", 7.5},
		{"0", 0},
		{"abc", 0}, // invalid → unset (0)
		{"", 0},
	}
	for _, tt := range tests {
		t.Run("min_rating="+tt.val, func(t *testing.T) {
			q := url.Values{}
			if tt.val != "" {
				q.Set("min_rating", tt.val)
			}
			f := parsePlaceListFilter(q)
			if f.MinRating != tt.want {
				t.Errorf("want %v, got %v", tt.want, f.MinRating)
			}
		})
	}
}

func TestParsePlaceListFilter_StringFields(t *testing.T) {
	q := url.Values{}
	q.Set("city", "Москва")
	q.Set("search", "пиццерия")
	q.Set("sort", "rating")
	f := parsePlaceListFilter(q)
	if f.City != "Москва" || f.Search != "пиццерия" || f.Sort != "rating" {
		t.Errorf("unexpected: %+v", f)
	}
}

// --- handler-level tests (no DB) ---

func TestPlaceHandler_Create_Unauthorized(t *testing.T) {
	h := &PlaceHandler{placeRepo: nil}
	req := httptest.NewRequest("POST", "/api/places", bytes.NewBufferString(`{"name":"X"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestPlaceHandler_Create_BadInput(t *testing.T) {
	h := &PlaceHandler{placeRepo: nil}
	tests := []struct {
		name string
		body string
		code int
	}{
		{"invalid json", `{bad`, http.StatusBadRequest},
		{"missing name", `{"city":"Москва"}`, http.StatusBadRequest},
		{"empty name", `{"name":""}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/places", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req = withUserID(req, 1)
			w := httptest.NewRecorder()
			h.Create(w, req)
			if w.Code != tt.code {
				t.Errorf("want %d, got %d", tt.code, w.Code)
			}
		})
	}
}

func TestPlaceHandler_Update_Unauthorized(t *testing.T) {
	h := &PlaceHandler{placeRepo: nil}
	req := httptest.NewRequest("PUT", "/api/places/1", bytes.NewBufferString(`{"name":"X"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Update(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestPlaceHandler_Update_InvalidID(t *testing.T) {
	h := &PlaceHandler{placeRepo: nil}
	req := httptest.NewRequest("PUT", "/api/places/abc", bytes.NewBufferString(`{"name":"X"}`))
	req.Header.Set("Content-Type", "application/json")
	req = withUserID(req, 1)
	req = withChiURLParam(req, "id", "abc")
	w := httptest.NewRecorder()
	h.Update(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPlaceHandler_Delete_Unauthorized(t *testing.T) {
	h := &PlaceHandler{placeRepo: nil}
	req := httptest.NewRequest("DELETE", "/api/places/1", nil)
	w := httptest.NewRecorder()
	h.Delete(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestPlaceHandler_Delete_InvalidID(t *testing.T) {
	h := &PlaceHandler{placeRepo: nil}
	req := httptest.NewRequest("DELETE", "/api/places/abc", nil)
	req = withUserID(req, 1)
	req = withChiURLParam(req, "id", "abc")
	w := httptest.NewRecorder()
	h.Delete(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPlaceHandler_GetByID_InvalidID(t *testing.T) {
	h := &PlaceHandler{placeRepo: nil}
	req := httptest.NewRequest("GET", "/api/places/abc", nil)
	req = withChiURLParam(req, "id", "abc")
	w := httptest.NewRecorder()
	h.GetByID(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPlaceHandler_UploadImage_Unauthorized(t *testing.T) {
	h := &PlaceHandler{placeRepo: nil}
	req := httptest.NewRequest("POST", "/api/places/1/image", nil)
	w := httptest.NewRecorder()
	h.UploadImage(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// Sanity: confirm the filter struct passes through repository's expected shape.
func TestParsePlaceListFilter_TypeShape(t *testing.T) {
	q := url.Values{}
	q.Set("cuisine_type_id", "1,2")
	q.Set("category_id", "3")
	q.Set("city", "Spb")
	q.Set("search", "bar")
	q.Set("sort", "name")
	q.Set("min_rating", "5")
	q.Set("is_gem", "true")
	q.Set("limit", "5")
	q.Set("page", "2")
	got := parsePlaceListFilter(q)
	want := repository.PlaceFilter{
		City:           "Spb",
		Search:         "bar",
		Sort:           "name",
		CuisineTypeIDs: []int{1, 2},
		CategoryIDs:    []int{3},
		MinRating:      5,
		IsGem:          boolPtr(true),
		Limit:          5,
		Offset:         5,
	}
	// IsGem comparison: dereference if both non-nil
	if got.IsGem == nil || *got.IsGem != *want.IsGem {
		t.Errorf("IsGem mismatch: got=%v want=%v", got.IsGem, want.IsGem)
	}
	got.IsGem, want.IsGem = nil, nil
	if !reflect.DeepEqual(got, want) {
		t.Errorf("filter mismatch:\n got: %+v\nwant: %+v", got, want)
	}
}

// --- helpers ---

func boolPtr(b bool) *bool { return &b }

func withChiURLParam(r *http.Request, key, val string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
