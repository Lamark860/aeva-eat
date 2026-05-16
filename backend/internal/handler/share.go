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
	ID       int
	Name     string
	City     string
	Cuisine  string
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
		ID:      p.ID,
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
	if p.CuisineType != nil {
		d.Cuisine = *p.CuisineType
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
  --ink-soft:     oklch(0.45 0.018 60);
  --terracotta:   oklch(0.58 0.14 30);
  --terra-dark:   oklch(0.45 0.14 30);
  --tape:         oklch(0.92 0.08 70 / 0.7);
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
  position: relative;
}

/* AEVA EAT wordmark — рукописно, верхний правый, кивок к app-каркасу */
.wordmark {
  position: absolute;
  top: 12px;
  right: 16px;
  font-family: var(--hand);
  font-size: 18px;
  letter-spacing: 0.08em;
  color: var(--terracotta);
  z-index: 3;
  user-select: none;
}

/* B5 — cover ~42% высоты (раньше 64vh — слишком много).
   Под плашкой бумаги виден ещё ~58% — там CTA и подпись «камерный дневник еды». */
.cover {
  width: 100%;
  height: 42vh;
  max-height: 320px;
  background: oklch(0.55 0.06 30);
  background-size: cover;
  background-position: center;
}

/* Бумажная плашка — заходит на cover и продолжается под ним. Tape сверху. */
.card-wrap {
  position: relative;
  margin: -48px 18px 18px;
  z-index: 2;
}
.tape {
  position: absolute;
  top: -10px;
  width: 90px;
  height: 22px;
  background: var(--tape);
  box-shadow: 0 1px 2px rgba(40, 30, 20, 0.08);
}
.tape.l { left: 24px; transform: rotate(-9deg); }
.tape.r { right: 24px; transform: rotate(8deg); background: oklch(0.88 0.06 200 / 0.7); }

.card {
  background: var(--paper-card);
  padding: 24px 24px 28px;
  border-radius: 1px;
  text-align: center;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.08),
    0 6px 18px rgba(40, 30, 20, 0.10);
}

.cap-top {
  font-family: var(--hand);
  font-size: 14px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--ink-soft);
  margin-bottom: 12px;
}

.name {
  font-family: var(--serif);
  font-style: italic;
  font-weight: 500;
  font-size: 30px;
  line-height: 1.1;
  margin: 0 0 6px;
}
.meta {
  font-family: var(--hand);
  font-size: 18px;
  color: var(--ink-mute);
  margin-bottom: 12px;
}

.gem-stamp {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin: 4px 0 14px;
}
.gem-stamp .diamond {
  font-size: 18px;
  color: var(--terracotta);
  line-height: 1;
}
.gem-stamp .label {
  font-family: var(--serif);
  font-weight: 600;
  font-size: 11px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--terracotta);
  border: 1.4px solid var(--terracotta);
  padding: 3px 9px 2px;
  border-radius: 2px;
  background: oklch(0.94 0.05 85 / 0.5);
}

/* B5 — CTA как бумажная prompt-кнопка с 1° наклоном, НЕ терракотовый pill.
   Бумажная плашка с пунктирной рамкой, arrow-prefix → */
.cta {
  display: block;
  margin: 14px auto 10px;
  padding: 12px 22px;
  background: var(--paper-card);
  color: var(--ink);
  font-family: var(--serif);
  font-style: italic;
  font-size: 16px;
  text-decoration: none;
  border: 1px dashed var(--terra-dark);
  border-radius: 2px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 3px 8px rgba(40, 30, 20, 0.08);
  transform: rotate(-1deg);
  max-width: 320px;
}
.cta:hover { background: oklch(0.96 0.03 85); }
.cta .arrow { color: var(--terracotta); margin-right: 4px; }

.tagline {
  font-family: var(--hand);
  font-size: 15px;
  color: var(--ink-mute);
  text-align: center;
  margin-top: 12px;
}
</style>
</head>
<body>
<div class="wordmark">AEVA·EAT</div>
<div class="cover" {{if .ImageURL}}style="background-image: url('{{.ImageURL}}')"{{end}}></div>
<div class="card-wrap">
  <span class="tape l"></span>
  <span class="tape r"></span>
  <div class="card">
    <div class="cap-top">из дневника круга</div>
    <h1 class="name">{{.Name}}</h1>
    {{if or .City .Cuisine}}<div class="meta">
      {{if .City}}{{.City}}{{end}}{{if and .City .Cuisine}} · {{end}}{{if .Cuisine}}{{.Cuisine}}{{end}}
    </div>{{end}}
    {{if .IsGem}}<div class="gem-stamp">
      <span class="diamond">◆</span><span class="label">жемчужина</span>
    </div>{{end}}
    <a class="cta" href="/login?next=/places/{{.ID}}">
      <span class="arrow">→</span> войти, чтобы увидеть впечатления
    </a>
    <div class="tagline">камерный дневник еды</div>
  </div>
</div>
</body>
</html>
`
