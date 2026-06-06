package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aeva-eat/backend/internal/config"
	"github.com/aeva-eat/backend/internal/handler"
	"github.com/aeva-eat/backend/internal/middleware"
	"github.com/aeva-eat/backend/internal/repository"
	"github.com/aeva-eat/backend/internal/service"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to database")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied")

	// Repos
	userRepo := repository.NewUserRepo(db)
	placeRepo := repository.NewPlaceRepo(db)
	catalogRepo := repository.NewCatalogRepo(db)
	reviewRepo := repository.NewReviewRepo(db)
	wishlistRepo := repository.NewWishlistRepo(db)
	inviteRepo := repository.NewInviteRepo(db)
	noteRepo := repository.NewNoteRepo(db)
	feedRepo := repository.NewFeedEventsRepo(db)
	aggRepo := repository.NewAggregateRepo(db)

	// Services
	authService := service.NewAuthService(userRepo, inviteRepo, cfg.JWTSecret)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	placeHandler := handler.NewPlaceHandler(placeRepo, userRepo)
	catalogHandler := handler.NewCatalogHandler(catalogRepo)
	reviewHandler := handler.NewReviewHandler(reviewRepo, wishlistRepo)
	wishlistHandler := handler.NewWishlistHandler(wishlistRepo)
	suggestHandler := handler.NewSuggestHandler(cfg.GeosuggestKey)
	inviteHandler := handler.NewInviteHandler(inviteRepo, userRepo)
	noteHandler := handler.NewNoteHandler(noteRepo, feedRepo)
	aggHandler := handler.NewAggregateHandler(aggRepo, placeRepo)
	shareHandler := handler.NewShareHandler(placeRepo)

	// Router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	// Лимит JSON-тела (multipart-загрузки не трогает — у них свои лимиты).
	r.Use(middleware.BodyLimit(1 << 20)) // 1 MB

	// Rate-limit: брутфорс login/register по IP и слив квоты suggest по userID.
	authLimiter := middleware.NewRateLimiter(20, 10, middleware.IPKey)
	suggestLimiter := middleware.NewRateLimiter(60, 20, middleware.UserKey)

	// Routes
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.With(authLimiter.Handler).Post("/register", authHandler.Register)
		r.With(authLimiter.Handler).Post("/login", authHandler.Login)
		r.With(middleware.JWTAuth(cfg.JWTSecret)).Get("/me", authHandler.Me)
		r.With(middleware.JWTAuth(cfg.JWTSecret)).Put("/password", authHandler.ChangePassword)
		r.With(middleware.JWTAuth(cfg.JWTSecret)).Post("/avatar", authHandler.UploadAvatar)
	})

	// Validate invite (public, no auth needed)
	r.Get("/api/invites/validate/{code}", inviteHandler.ValidateCode)

	// DESIGN-DECISIONS Q3 — публичные share-страницы /p/<id>. Без auth, HTML
	// с OG-метатегами для preview в мессенджерах. nginx проксирует /p/<digits>
	// на бэк (см. nginx/nginx.conf).
	r.Get("/p/{id}", shareHandler.Render)

	// All data routes require auth
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))

		// Places
		r.Route("/api/places", func(r chi.Router) {
			r.Get("/", placeHandler.List)
			r.Get("/cities", placeHandler.ListCities)
			r.Get("/{id}", placeHandler.GetByID)
			r.Get("/{id}/reviews", reviewHandler.ListByPlace)
			r.Post("/", placeHandler.Create)
			r.Put("/{id}", placeHandler.Update)
			r.Delete("/{id}", placeHandler.Delete)
			r.Post("/{id}/restore", placeHandler.Restore) // superuser

			r.Post("/{id}/image", placeHandler.UploadImage)
			r.Post("/{id}/reviews", reviewHandler.Create)
			r.Put("/{id}/reviews/{rid}", reviewHandler.Update)
			r.Delete("/{id}/reviews/{rid}", reviewHandler.Delete)
			r.Post("/{id}/reviews/{rid}/image", reviewHandler.UploadImage)
			r.Post("/{id}/reviews/{rid}/video", reviewHandler.UploadVideo)
			r.Post("/{id}/reviews/{rid}/photos", reviewHandler.UploadPhotos)
			r.Delete("/{id}/reviews/{rid}/photos/{pid}", reviewHandler.DeletePhoto)
		})

		// Random place
		r.Get("/api/random", placeHandler.Random)

		// Notes
		r.Route("/api/notes", func(r chi.Router) {
			r.Get("/", noteHandler.List)
			r.Post("/", noteHandler.Create)
			r.Put("/{id}", noteHandler.Update)
			r.Delete("/{id}", noteHandler.Delete)
			r.Put("/{id}/strike", noteHandler.Strike)
		})

		// Feed events (общая хронология) + unread-индикатор
		r.Get("/api/feed", noteHandler.Feed)
		r.Get("/api/feed/weeks", noteHandler.Weeks)
		r.Get("/api/feed/unread-count", noteHandler.UnreadCount)
		r.Post("/api/feed/seen", noteHandler.MarkSeen)

		// Cities — путеводитель по городам круга
		r.Route("/api/cities", func(r chi.Router) {
			r.Get("/", aggHandler.ListCities)
			r.Get("/{name}", aggHandler.GetCity)
			r.Get("/{name}/places", aggHandler.ListCityPlaces)
			r.Get("/{name}/gems", aggHandler.ListCityGems)
		})

		// Users — публичный профиль и его контент
		r.Route("/api/users", func(r chi.Router) {
			r.Get("/", aggHandler.ListUsers)
			r.Get("/{id}", aggHandler.GetUser)
			r.Get("/{id}/places", aggHandler.ListUserPlaces)
			r.Get("/{id}/gems", aggHandler.ListUserGems)
			r.Get("/{id}/cities", aggHandler.ListUserCities)
			r.Get("/{userId}/reviews", reviewHandler.ListByUser)
		})

		// Gems hub
		r.Get("/api/gems", aggHandler.Gems)

		// Wishlist
		r.Route("/api/wishlist", func(r chi.Router) {
			r.Get("/", wishlistHandler.ListMy)
			r.Get("/all", wishlistHandler.ListAll)
			r.Get("/ids", wishlistHandler.MyIDs)
			r.Post("/{id}", wishlistHandler.Add)
			r.Delete("/{id}", wishlistHandler.Remove)
			r.Get("/custom", wishlistHandler.ListCustom)
			r.Post("/custom", wishlistHandler.AddCustom)
			r.Delete("/custom/{id}", wishlistHandler.DeleteCustom)
		})

		// Catalogs
		r.Get("/api/cuisine-types", catalogHandler.ListCuisineTypes)
		r.Get("/api/categories", catalogHandler.ListCategories)

		// Geosuggest proxy
		r.With(suggestLimiter.Handler).Get("/api/suggest", suggestHandler.Suggest)

		// Invites
		r.Route("/api/invites", func(r chi.Router) {
			r.Get("/", inviteHandler.ListMy)
			r.Post("/", inviteHandler.Create)
			r.Delete("/{id}", inviteHandler.Delete)
			r.Get("/all", inviteHandler.ListAll) // superuser only
		})

		// Admin
		r.Get("/api/admin/users", inviteHandler.ListUsers) // superuser only
	})

	addr := fmt.Sprintf(":%s", cfg.APIPort)
	log.Printf("starting server on %s", addr)
	// ReadHeaderTimeout защищает от Slowloris (медленные заголовки), не ограничивая
	// время заливки тела — иначе оборвутся 20MB-видео на медленном мобильном.
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	files := []string{
		"migrations/001_init.up.sql",
		"migrations/002_place_image.up.sql",
		"migrations/003_rating_float.up.sql",
		"migrations/004_wishlist.up.sql",
		"migrations/006_review_photos_custom_wishlist.up.sql",
		"migrations/007_auth_username_login.up.sql",
		"migrations/008_place_unique.up.sql",
		"migrations/009_invites_roles.up.sql",
		"migrations/010_review_video.up.sql",
		"migrations/011_dict_extensions.up.sql",
		"migrations/012_review_photos.up.sql",
		"migrations/013_notes.up.sql",
		"migrations/014_feed_unread.up.sql",
		"migrations/015_wishlist_struck.up.sql",
		"migrations/016_place_identity.up.sql",
		"migrations/017_place_soft_delete.up.sql",
		// ВНИМАНИЕ: 005_seed_data — это ДЕМО-данные, не схема. Раньше он стоял
		// здесь и прогонялся при каждом старте; т.к. wishlist_custom не имеет
		// уникального индекса, ON CONFLICT не срабатывал и записи плодились
		// (по копии на рестарт). Демо-сид теперь руками: backend/scripts/seed_demo.sh.
	}
	// Леджер применённых миграций: больше не гоняем весь SQL на каждом старте,
	// применяем только новые файлы и фиксируем версию. Существующие up-скрипты
	// остаются идемпотентными (IF NOT EXISTS / ON CONFLICT), поэтому первый старт
	// с леджером на уже мигрированной проде безопасно перенакатит их по разу и
	// запишет — дальше каждый файл применяется максимум once.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`); err != nil {
		return fmt.Errorf("creating schema_migrations: %w", err)
	}

	for _, f := range files {
		var exists bool
		if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, f).Scan(&exists); err != nil {
			return fmt.Errorf("checking migration %s: %w", f, err)
		}
		if exists {
			continue
		}

		migrationSQL, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", f, err)
		}

		// Каждая миграция — в транзакции: запись в леджер коммитится вместе с
		// самой миграцией, поэтому частично-применённого состояния не остаётся.
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", f, err)
		}
		if _, err := tx.Exec(string(migrationSQL)); err != nil {
			tx.Rollback()
			return fmt.Errorf("executing migration %s: %w", f, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, f); err != nil {
			tx.Rollback()
			return fmt.Errorf("recording migration %s: %w", f, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", f, err)
		}
		log.Printf("applied migration %s", f)
	}
	return nil
}
