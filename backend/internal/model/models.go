package model

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	DisplayName  *string   `json:"display_name,omitempty"`
	PasswordHash string    `json:"-"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Place struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Address       *string    `json:"address,omitempty"`
	City          *string    `json:"city,omitempty"`
	Lat           *float64   `json:"lat,omitempty"`
	Lng           *float64   `json:"lng,omitempty"`
	CuisineTypeID *int       `json:"cuisine_type_id,omitempty"`
	CuisineType   *string    `json:"cuisine_type,omitempty"`
	Website       *string    `json:"website,omitempty"`
	CreatedBy     *int       `json:"created_by,omitempty"`
	Categories    []string   `json:"categories,omitempty"`
	AvgFood       *float64   `json:"avg_food,omitempty"`
	AvgService    *float64   `json:"avg_service,omitempty"`
	AvgVibe       *float64   `json:"avg_vibe,omitempty"`
	ImageURL      *string    `json:"image_url,omitempty"`
	IsGemPlace    bool       `json:"is_gem_place"`
	ReviewCount   int        `json:"review_count"`
	Reviewers     []Reviewer `json:"reviewers,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Reviewer struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Review struct {
	ID            int       `json:"id"`
	PlaceID       int       `json:"place_id"`
	PlaceName     string    `json:"place_name,omitempty"`
	FoodRating    float64   `json:"food_rating"`
	ServiceRating float64   `json:"service_rating"`
	VibeRating    float64   `json:"vibe_rating"`
	IsGem         bool      `json:"is_gem"`
	Comment       *string   `json:"comment,omitempty"`
	ImageURL      *string   `json:"image_url,omitempty"`
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

type WishlistCustom struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Note      *string   `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
