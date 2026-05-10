package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/aeva-eat/backend/internal/repository"
	"github.com/go-chi/chi/v5"
)

// ShareHandler рендерит публичные страницы /p/<place_id> — превью места без
// авторизации, с OG-метатегами для красивого preview в мессенджерах.
// DESIGN-DECISIONS Q3: cover full-bleed → бумажная плашка → серифа имя →
// каракули город → штамп жемчужины (если is_gem). БЕЗ рейтингов / БЕЗ имён авторов.
type ShareHandler struct {
	placeRepo *repository.PlaceRepo
	tpl       *template.Template
}

func NewShareHandler(placeRepo *repository.PlaceRepo) *ShareHandler {
	tpl := template.Must(template.New("share").Parse(shareTemplate))
	return &ShareHandler{placeRepo: placeRepo, tpl: tpl}
}

type shareData struct {
	Name     string
	City     string
	ImageURL string
	IsGem    bool
	OGTitle  string
	OGImage  string
}

// Render — GET /p/<id>. Без auth. 404 если места нет.
func (h *ShareHandler) Render(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	p, err := h.placeRepo.GetByID(id)
	if err != nil || p == nil {
		http.NotFound(w, r)
		return
	}

	d := shareData{
		Name:    p.Name,
		IsGem:   p.IsGemPlace,
		OGTitle: p.Name,
	}
	if p.City != nil {
		d.City = *p.City
		if d.City != "" {
			d.OGTitle = p.Name + " · " + d.City
		}
	}
	if p.ImageURL != nil {
		d.ImageURL = *p.ImageURL
		d.OGImage = absURL(r, *p.ImageURL)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Превью-страницы можно кэшировать минут на 5 — если место поправили,
	// мессенджер всё равно держит свой preview-кэш дольше.
	w.Header().Set("Cache-Control", "public, max-age=300")
	if err := h.tpl.Execute(w, d); err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// absURL — превращает относительный путь /uploads/... в абсолютный URL для
// og:image, иначе мессенджеры не смогут вытащить картинку.
func absURL(r *http.Request, path string) string {
	if len(path) >= 4 && (path[:4] == "http") {
		return path
	}
	scheme := "https"
	if xfp := r.Header.Get("X-Forwarded-Proto"); xfp != "" {
		scheme = xfp
	} else if r.TLS == nil {
		scheme = "http"
	}
	return scheme + "://" + r.Host + path
}

const shareTemplate = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
<title>{{.OGTitle}}</title>

<meta property="og:type" content="article">
<meta property="og:title" content="{{.OGTitle}}">
<meta property="og:description" content="камерный дневник еды">
{{if .OGImage}}<meta property="og:image" content="{{.OGImage}}">{{end}}
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.OGTitle}}">
<meta name="twitter:description" content="камерный дневник еды">
{{if .OGImage}}<meta name="twitter:image" content="{{.OGImage}}">{{end}}

<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Lora:ital,wght@0,500;1,500&family=Caveat:wght@500&display=swap" rel="stylesheet">

<style>
:root {
  color-scheme: light;
  --paper:        oklch(0.965 0.018 82);
  --paper-card:   #fdfcf7;
  --ink:          oklch(0.22 0.02 60);
  --ink-mute:     oklch(0.6 0.015 60);
  --terracotta:   oklch(0.58 0.14 30);
  --serif:        'Lora', Georgia, serif;
  --hand:         'Caveat', 'Marker Felt', cursive;
}
* { box-sizing: border-box; }
body {
  margin: 0;
  background: var(--paper);
  color: var(--ink);
  font-family: var(--serif);
  min-height: 100vh;
}
.cover {
  width: 100%;
  height: 64vh;
  max-height: 480px;
  background: oklch(0.55 0.06 30);
  background-size: cover;
  background-position: center;
}
.card {
  background: var(--paper-card);
  margin: -32px 16px 0;
  padding: 22px 22px 28px;
  border-radius: 1px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.08),
    0 6px 18px rgba(40, 30, 20, 0.10);
  text-align: center;
  position: relative;
  z-index: 2;
}
.name {
  font-family: var(--serif);
  font-style: italic;
  font-weight: 500;
  font-size: 28px;
  line-height: 1.15;
  margin: 0 0 6px;
}
.city {
  font-family: var(--hand);
  font-size: 20px;
  color: var(--ink-mute);
  margin-bottom: 10px;
}
.gem {
  display: inline-block;
  font-family: var(--serif);
  font-weight: 600;
  font-size: 10px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--terracotta);
  border: 1.4px solid var(--terracotta);
  padding: 3px 8px 2px;
  border-radius: 2px;
  background: oklch(0.94 0.05 85 / 0.5);
  margin-bottom: 14px;
}
.cta {
  display: inline-block;
  margin-top: 8px;
  padding: 12px 22px;
  background: var(--terracotta);
  color: #fff;
  font-family: var(--serif);
  font-style: italic;
  font-size: 16px;
  text-decoration: none;
  border-radius: 999px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.1),
    0 4px 10px rgba(40, 30, 20, 0.15);
}
</style>
</head>
<body>
<div class="cover" {{if .ImageURL}}style="background-image: url('{{.ImageURL}}')"{{end}}></div>
<div class="card">
  <h1 class="name">{{.Name}}</h1>
  {{if .City}}<div class="city">{{.City}}</div>{{end}}
  {{if .IsGem}}<div class="gem">жемчужина</div>{{end}}
  <a class="cta" href="/login">войти, чтобы увидеть наши впечатления</a>
</div>
</body>
</html>
`
