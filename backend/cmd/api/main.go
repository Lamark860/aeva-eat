package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	placeHandler := handler.NewPlaceHandler(placeRepo)
	catalogHandler := handler.NewCatalogHandler(catalogRepo)
	reviewHandler := handler.NewReviewHandler(reviewRepo)
	wishlistHandler := handler.NewWishlistHandler(wishlistRepo)
	suggestHandler := handler.NewSuggestHandler(cfg.GeosuggestKey)

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

	// Routes
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.With(middleware.JWTAuth(cfg.JWTSecret)).Get("/me", authHandler.Me)
	})

	// Places
	r.Route("/api/places", func(r chi.Router) {
		r.Get("/", placeHandler.List)
		r.Get("/cities", placeHandler.ListCities)
		r.Get("/{id}", placeHandler.GetByID)
		r.Get("/{id}/reviews", reviewHandler.ListByPlace)
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg.JWTSecret))
			r.Post("/", placeHandler.Create)
			r.Put("/{id}", placeHandler.Update)
			r.Delete("/{id}", placeHandler.Delete)
			r.Post("/{id}/image", placeHandler.UploadImage)
			r.Post("/{id}/reviews", reviewHandler.Create)
			r.Put("/{id}/reviews/{rid}", reviewHandler.Update)
			r.Delete("/{id}/reviews/{rid}", reviewHandler.Delete)
			r.Post("/{id}/reviews/{rid}/image", reviewHandler.UploadImage)
		})
	})

	// User reviews
	r.Get("/api/users/{userId}/reviews", reviewHandler.ListByUser)

	// Wishlist
	r.Route("/api/wishlist", func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Get("/", wishlistHandler.ListMy)
		r.Get("/ids", wishlistHandler.MyIDs)
		r.Post("/{id}", wishlistHandler.Add)
		r.Delete("/{id}", wishlistHandler.Remove)
		// Custom free-text wishlist
		r.Get("/custom", wishlistHandler.ListCustom)
		r.Post("/custom", wishlistHandler.AddCustom)
		r.Delete("/custom/{id}", wishlistHandler.DeleteCustom)
	})

	// Catalogs
	r.Get("/api/cuisine-types", catalogHandler.ListCuisineTypes)
	r.Get("/api/categories", catalogHandler.ListCategories)

	// Geosuggest proxy
	r.Get("/api/suggest", suggestHandler.Suggest)

	addr := fmt.Sprintf(":%s", cfg.APIPort)
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
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
		"migrations/005_seed_data.up.sql",
	}
	for _, f := range files {
		migrationSQL, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", f, err)
		}
		_, err = db.Exec(string(migrationSQL))
		if err != nil {
			return fmt.Errorf("executing migration %s: %w", f, err)
		}
	}
	return nil
}
