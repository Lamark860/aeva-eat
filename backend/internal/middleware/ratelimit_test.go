package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func constKey(string) func(*http.Request) string {
	return func(*http.Request) string { return "k" }
}

func TestRateLimiter_BurstThenBlock(t *testing.T) {
	// perMinute=0 → пополнения нет, доступен ровно burst запросов.
	rl := NewRateLimiter(0, 3, constKey(""))
	for i := 0; i < 3; i++ {
		if !rl.allow("k") {
			t.Fatalf("request %d within burst should be allowed", i+1)
		}
	}
	if rl.allow("k") {
		t.Fatal("request beyond burst should be blocked")
	}
}

func TestRateLimiter_SeparateKeys(t *testing.T) {
	rl := NewRateLimiter(0, 1, constKey(""))
	if !rl.allow("a") || !rl.allow("b") {
		t.Fatal("different keys must have independent buckets")
	}
	if rl.allow("a") {
		t.Fatal("key a exhausted its burst")
	}
}

func TestRateLimiter_Handler429(t *testing.T) {
	rl := NewRateLimiter(0, 1, func(*http.Request) string { return "same" })
	h := rl.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, httptest.NewRequest("POST", "/login", nil))
	if rec1.Code != http.StatusOK {
		t.Fatalf("first request: want 200, got %d", rec1.Code)
	}

	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, httptest.NewRequest("POST", "/login", nil))
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request: want 429, got %d", rec2.Code)
	}
}
