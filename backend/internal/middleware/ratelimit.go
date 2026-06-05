package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RateLimiter — простой in-memory token-bucket по ключу (IP или userID).
// Для invite-only приложения этого достаточно: ключей немного, внешних
// зависимостей не нужно. Защищает login/register от брутфорса и suggest-прокси
// от слива платной квоты Яндекса.
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
	rate    float64 // токенов в секунду
	burst   float64 // максимум накопленных токенов
	keyFn   func(*http.Request) string
}

type tokenBucket struct {
	tokens float64
	last   time.Time
}

// NewRateLimiter: perMinute — установившаяся скорость, burst — мгновенный запас.
func NewRateLimiter(perMinute, burst int, keyFn func(*http.Request) string) *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*tokenBucket),
		rate:    float64(perMinute) / 60.0,
		burst:   float64(burst),
		keyFn:   keyFn,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	b := rl.buckets[key]
	if b == nil {
		b = &tokenBucket{tokens: rl.burst, last: now}
		rl.buckets[key] = b
	}
	b.tokens += now.Sub(b.last).Seconds() * rl.rate
	if b.tokens > rl.burst {
		b.tokens = rl.burst
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// cleanupLoop изредка чистит протухшие бакеты, чтобы map не рос вечно.
func (rl *RateLimiter) cleanupLoop() {
	t := time.NewTicker(10 * time.Minute)
	for range t.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for k, b := range rl.buckets {
			if b.last.Before(cutoff) {
				delete(rl.buckets, k)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(rl.keyFn(r)) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"too many requests, try again later"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// IPKey — ключ по IP (RealIP middleware уже переписывает RemoteAddr на реальный).
func IPKey(r *http.Request) string {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// UserKey — ключ по userID из JWT (фоллбэк на IP, если токена нет).
func UserKey(r *http.Request) string {
	if uid, ok := GetUserID(r); ok {
		return "u:" + strconv.Itoa(uid)
	}
	return IPKey(r)
}

// BodyLimit ограничивает размер JSON-тела (multipart-загрузки не трогает — у них
// свои лимиты в ParseMultipartForm). Защита от раздувания памяти большим payload.
func BodyLimit(max int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				r.Body = http.MaxBytesReader(w, r.Body, max)
			}
			next.ServeHTTP(w, r)
		})
	}
}
