package handler

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

// suggestClient — отдельный клиент с таймаутом: дефолтный http.Get без таймаута
// при зависшем Яндексе копит горутины-обработчики без ограничения.
var suggestClient = &http.Client{Timeout: 5 * time.Second}

type SuggestHandler struct {
	apiKey string
}

func NewSuggestHandler(apiKey string) *SuggestHandler {
	return &SuggestHandler{apiKey: apiKey}
}

func (h *SuggestHandler) Suggest(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("text")
	if q == "" {
		http.Error(w, `{"error":"text required"}`, http.StatusBadRequest)
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
	if ll := r.URL.Query().Get("ll"); ll != "" {
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
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
