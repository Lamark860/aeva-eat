package handler

import (
	"io"
	"net/http"
	"net/url"
)

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

	resp, err := http.Get("https://suggest-maps.yandex.ru/v1/suggest?" + params.Encode())
	if err != nil {
		http.Error(w, `{"error":"upstream error"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
