package model

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	DisplayName  *string   `json:"display_name,omitempty"`
	PasswordHash string    `json:"-"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Invite struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	CreatedBy   int        `json:"created_by"`
	CreatorName string     `json:"creator_name,omitempty"`
	UsedBy      *int       `json:"used_by,omitempty"`
	UsedAt      *time.Time `json:"used_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Place struct {
	ID                 int           `json:"id"`
	Name               string        `json:"name"`
	Address            *string       `json:"address,omitempty"`
	City               *string       `json:"city,omitempty"`
	Lat                *float64      `json:"lat,omitempty"`
	Lng                *float64      `json:"lng,omitempty"`
	CuisineTypeID      *int          `json:"cuisine_type_id,omitempty"`
	CuisineType        *string       `json:"cuisine_type,omitempty"`
	Website            *string       `json:"website,omitempty"`
	CreatedBy          *int          `json:"created_by,omitempty"`
	Categories         []string      `json:"categories,omitempty"`
	AvgFood            *float64      `json:"avg_food,omitempty"`
	AvgService         *float64      `json:"avg_service,omitempty"`
	AvgVibe            *float64      `json:"avg_vibe,omitempty"`
	ImageURL           *string       `json:"image_url,omitempty"`
	IsGemPlace         bool          `json:"is_gem_place"`
	HasVideo           bool          `json:"has_video"`
	// VideoURL — последний review.video_url для этого места (если есть).
	// Фронт нужен, чтобы рендерить кружочек с реальным video-poster, а не
	// тёмной заглушкой; обновляется автоматически как review добавляются.
	VideoURL *string `json:"video_url,omitempty"`
	ReviewCount        int           `json:"review_count"`
	Reviewers          []Reviewer    `json:"reviewers,omitempty"`
	FeedPhotos         []ReviewPhoto `json:"feed_photos,omitempty"`
	// DESIGN-DECISIONS Q2 — populated only on GET /api/places/:id (not on list).
	GemStatus      *GemStatus       `json:"gem_status,omitempty"`
	Attendance     []Attendance     `json:"attendance,omitempty"`
	RatingsPerUser []RatingsPerUser `json:"ratings_per_user,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// GemStatus — кто первый отметил место жемчужиной + остальные «подписавшиеся».
// Q2: рукописная подпись «отметила Аня · 12 марта (+ Серёжа, Миша)» под штампом.
type GemStatus struct {
	MarkedBy      []Reviewer `json:"marked_by"`
	FirstMarkedAt time.Time  `json:"first_marked_at"`
}

// Attendance — кто и сколько раз был. ×N рендерится рукописным Caveat.
type Attendance struct {
	UserID     int     `json:"user_id"`
	Username   string  `json:"username"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	VisitCount int     `json:"visit_count"`
}

// RatingsPerUser — усреднённые оценки этого пользователя по всем визитам в это место.
// Только для бэка: сортировка `rating_user:N` + потенциальная мини-сравнялка
// «Аня: 9/8/9 · ты: 7/7/8» в шапке (в next).
type RatingsPerUser struct {
	UserID  int      `json:"user_id"`
	Food    *float64 `json:"food,omitempty"`
	Service *float64 `json:"service,omitempty"`
	Vibe    *float64 `json:"vibe,omitempty"`
}

type Reviewer struct {
	ID        int     `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type Review struct {
	ID            int           `json:"id"`
	PlaceID       int           `json:"place_id"`
	PlaceName     string        `json:"place_name,omitempty"`
	FoodRating    float64       `json:"food_rating"`
	ServiceRating float64       `json:"service_rating"`
	VibeRating    float64       `json:"vibe_rating"`
	IsGem         bool          `json:"is_gem"`
	Comment       *string       `json:"comment,omitempty"`
	ImageURL      *string       `json:"image_url,omitempty"`
	VideoURL      *string       `json:"video_url,omitempty"`
	VisitedAt     *string       `json:"visited_at,omitempty"`
	Authors       []User        `json:"authors"`
	Photos        []ReviewPhoto `json:"photos"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type ReviewPhoto struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Position int    `json:"position"`
}

type CuisineType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// WishlistEntry — запись общего wishlist круга для GET /api/wishlist/all.
// Содержит автора, место и флаг "зачёркнуто" (визит уже состоялся).
type WishlistEntry struct {
	UserID    int        `json:"user_id"`
	Username  string     `json:"username"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	Place     Place      `json:"place"`
	IsStruck  bool       `json:"is_struck"`
	StruckAt  *time.Time `json:"struck_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type WishlistCustom struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Note      *string   `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Note — записка от руки на доске (backend.md §notes). Может быть привязана
// к месту/городу, имеет цвет бумаги/тейпа и флаг "зачёркнуто".
type Note struct {
	ID         int       `json:"id"`
	AuthorID   int       `json:"author_id"`
	Author     *User     `json:"author,omitempty"`
	Text       string    `json:"text"`
	PlaceID    *int      `json:"place_id,omitempty"`
	PlaceName  *string   `json:"place_name,omitempty"`
	City       *string   `json:"city,omitempty"`
	PaperColor *string   `json:"paper_color,omitempty"`
	TapeColor  *string   `json:"tape_color,omitempty"`
	IsStruck   bool      `json:"is_struck"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FeedWeek — агрегат недели для свернутых полосок CollapsedStrip
// (backend.md §Лента/Доска: GET /api/feed/weeks). WeekStart — понедельник
// 00:00 UTC, ISO-week формата YYYY-Www в Key.
type FeedWeek struct {
	Key       string    `json:"key"`        // YYYY-Www, например 2026-W19
	WeekStart time.Time `json:"week_start"`
	Count     int       `json:"count"`
	GemCount  int       `json:"gem_count"`
}

// FeedEvent — строка из VIEW feed_events. Объединяет review_added /
// note_added в одну хронологию. Поля review_id и note_id взаимоисключающие.
// Attendees — все участники события (для review_added это co-authors;
// для note_added — единственный автор). Нужно для Q4-группировки на фронте
// по (place_id, дата, набор-участников).
type FeedEvent struct {
	Kind       string    `json:"kind"`
	EventID    int       `json:"event_id"`
	OccurredAt time.Time `json:"occurred_at"`
	PlaceID    *int      `json:"place_id,omitempty"`
	AuthorID   *int      `json:"author_id,omitempty"`
	ReviewID   *int      `json:"review_id,omitempty"`
	NoteID     *int      `json:"note_id,omitempty"`
	Attendees  []int     `json:"attendees,omitempty"`
}

// CityAggregate — строка для /api/cities: имя города + счётчики мест,
// жемчужин и уникальных авторов отзывов в этом городе. Используется на
// странице города и в полке "По городам" в Найти.
type CityAggregate struct {
	City             string `json:"city"`
	Count            int    `json:"count"`
	GemCount         int    `json:"gem_count"`
	ContributorCount int    `json:"contributor_count"`
}

// UserProfile — публичный профиль пользователя для /api/users/:id.
// Содержит то, что не стыдно показать гостю круга: статистики посещений,
// жемчужин и городов. Email/role/password остаются на User.
type UserProfile struct {
	ID                   int     `json:"id"`
	Username             string  `json:"username"`
	DisplayName          *string `json:"display_name,omitempty"`
	AvatarURL            *string `json:"avatar_url,omitempty"`
	PlaceCount           int     `json:"place_count"`
	GemCount             int     `json:"gem_count"`
	CityCount            int     `json:"city_count"`
	ReviewCount          int     `json:"review_count"`
	// FavoriteCuisine — самая частая кухня по визитам пользователя. nil, если
	// у мест нет cuisine_type_id или у пользователя ещё нет визитов.
	// Фронт рендерит как «любит грузинскую — N раз» (рукописно, с глаголом).
	FavoriteCuisine      *string `json:"favorite_cuisine,omitempty"`
	FavoriteCuisineCount int     `json:"favorite_cuisine_count,omitempty"`
}

// UserCity — город визитов конкретного пользователя.
type UserCity struct {
	City  string `json:"city"`
	Count int    `json:"count"`
}

// GemsHub — ответ /api/gems: список мест-жемчужин + агрегаты по городам
// и авторам.
type GemsHub struct {
	Places []Place         `json:"places"`
	Total  int             `json:"total"`
	ByCity []CityAggregate `json:"by_city"`
	ByUser []UserGemCount  `json:"by_user"`
}

type UserGemCount struct {
	UserID    int     `json:"user_id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Count     int     `json:"count"`
}
