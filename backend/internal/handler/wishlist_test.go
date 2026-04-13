package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWishlistHandler_Add_Unauthorized(t *testing.T) {
	h := &WishlistHandler{wishlistRepo: nil}

	req := httptest.NewRequest("POST", "/api/wishlist/1", nil)
	w := httptest.NewRecorder()

	h.Add(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestWishlistHandler_Remove_Unauthorized(t *testing.T) {
	h := &WishlistHandler{wishlistRepo: nil}

	req := httptest.NewRequest("DELETE", "/api/wishlist/1", nil)
	w := httptest.NewRecorder()

	h.Remove(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestWishlistHandler_ListMy_Unauthorized(t *testing.T) {
	h := &WishlistHandler{wishlistRepo: nil}

	req := httptest.NewRequest("GET", "/api/wishlist", nil)
	w := httptest.NewRecorder()

	h.ListMy(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestWishlistHandler_MyIDs_Unauthorized(t *testing.T) {
	h := &WishlistHandler{wishlistRepo: nil}

	req := httptest.NewRequest("GET", "/api/wishlist/ids", nil)
	w := httptest.NewRecorder()

	h.MyIDs(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
