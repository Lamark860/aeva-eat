package model

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Place struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Address       *string   `json:"address,omitempty"`
	City          *string   `json:"city,omitempty"`
	Lat           *float64  `json:"lat,omitempty"`
	Lng           *float64  `json:"lng,omitempty"`
	CuisineTypeID *int      `json:"cuisine_type_id,omitempty"`
	CuisineType   *string   `json:"cuisine_type,omitempty"`
	Website       *string   `json:"website,omitempty"`
	CreatedBy     *int      `json:"created_by,omitempty"`
	Categories    []string  `json:"categories,omitempty"`
	AvgFood       *float64  `json:"avg_food,omitempty"`
	AvgService    *float64  `json:"avg_service,omitempty"`
	AvgVibe       *float64  `json:"avg_vibe,omitempty"`
	ImageURL      *string   `json:"image_url,omitempty"`
	ReviewCount   int       `json:"review_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Review struct {
	ID            int       `json:"id"`
	PlaceID       int       `json:"place_id"`
	FoodRating    int       `json:"food_rating"`
	ServiceRating int       `json:"service_rating"`
	VibeRating    int       `json:"vibe_rating"`
	IsGem         bool      `json:"is_gem"`
	Comment       *string   `json:"comment,omitempty"`
	VisitedAt     *string   `json:"visited_at,omitempty"`
	Authors       []User    `json:"authors"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CuisineType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
