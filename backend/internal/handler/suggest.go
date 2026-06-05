package handler

import (
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// suggestClient — отдельный клиент с таймаутом: дефолтный http.Get без таймаута
// при зависшем Яндексе копит горутины-обработчики без ограничения.
var suggestClient = &http.Client{Timeout: 5 * time.Second}

const (
	suggestCacheTTL = 60 * time.Second
	suggestCacheMax = 500
)

type SuggestHandler struct {
	apiKey string
	mu     sync.Mutex
	cache  map[string]suggestCacheEntry
}

type suggestCacheEntry struct {
	body []byte
	exp  time.Time
}

func NewSuggestHandler(apiKey string) *SuggestHandler {
	return &SuggestHandler{apiKey: apiKey, cache: make(map[string]suggestCacheEntry)}
}

func (h *SuggestHandler) cacheGet(key string) ([]byte, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	e, ok := h.cache[key]
	if !ok || time.Now().After(e.exp) {
		return nil, false
	}
	return e.body, true
}

func (h *SuggestHandler) cacheSet(key string, body []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := time.Now()
	if len(h.cache) >= suggestCacheMax {
		for k, e := range h.cache {
			if now.After(e.exp) {
				delete(h.cache, k)
			}
		}
		if len(h.cache) >= suggestCacheMax {
			h.cache = make(map[string]suggestCacheEntry) // крайний случай — сброс
		}
	}
	h.cache[key] = suggestCacheEntry{body: body, exp: now.Add(suggestCacheTTL)}
}

func (h *SuggestHandler) Suggest(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("text")
	if q == "" {
		http.Error(w, `{"error":"text required"}`, http.StatusBadRequest)
		return
	}
	ll := r.URL.Query().Get("ll")

	// Кэш по нормализованному запросу: автокомплит часто шлёт одни и те же
	// строки (повторный ввод, дебаунс) — экономит платную квоту Яндекса.
	cacheKey := q + "|" + ll
	if body, ok := h.cacheGet(cacheKey); ok {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
		return
	}

	params := url.Values{
		"apikey":  {h.apiKey},
		"text":    {q},
		"lang":    {"ru_RU"},
		"types":   {"biz,geo"},
		"results": {"7"},
		"attrs":   {"uri"},
	}
	if ll != "" {
		params.Set("ll", ll)
		params.Set("spn", "0.5,0.5")
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet,
		"https://suggest-maps.yandex.ru/v1/suggest?"+params.Encode(), nil)
	if err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadGateway)
		return
	}
	resp, err := suggestClient.Do(req)
	if err != nil {
		http.Error(w, `{"error":"upstream error"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		h.cacheSet(cacheKey, body)
		_, _ = w.Write(body)
		return
	}
	// Ошибки апстрима не кэшируем — отдаём как есть.
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
