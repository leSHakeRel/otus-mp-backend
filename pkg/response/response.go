package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type EveningResponse struct {
	ID          uuid.UUID    `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	ScheduledAt *time.Time   `json:"scheduled_at,omitempty"`
	Owner       UserResponse `json:"owner"`
	IsPrivate   bool         `json:"is_private"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type EveningFilmResponse struct {
	ID           uuid.UUID  `json:"id"`
	TMDBID       int        `json:"tmdb_id"`
	Title        string     `json:"title"`
	PosterPath   string     `json:"poster_path,omitempty"`
	BackdropPath string     `json:"backdrop_path,omitempty"`
	ReleaseDate  *time.Time `json:"release_date,omitempty"`
	VoteAverage  float64    `json:"vote_average,omitempty"`
	Overview     string     `json:"overview,omitempty"`
	AddedAt      time.Time  `json:"added_at"`
}

type VoteResponse struct {
	ID            uuid.UUID    `json:"id"`
	EveningFilmID uuid.UUID    `json:"evening_film_id"`
	User          UserResponse `json:"user"`
	Value         int          `json:"value"`
	CreatedAt     time.Time    `json:"created_at"`
}

type VoteSummary struct {
	EveningFilmID    uuid.UUID        `json:"evening_film_id"`
	Title            string           `json:"title"`
	PosterPath       string           `json:"poster_path,omitempty"`
	TotalVotes       int              `json:"total_votes"`
	AverageScore     float64          `json:"average_score"`
	VoteDistribution VoteDistribution `json:"vote_distribution"`
}

type VoteDistribution struct {
	One   int `json:"1"`
	Two   int `json:"2"`
	Three int `json:"3"`
	Four  int `json:"4"`
	Five  int `json:"5"`
}

type CommentResponse struct {
	ID        uuid.UUID    `json:"id"`
	User      UserResponse `json:"user"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"created_at"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type ErrorResponse struct {
	Error *AppError `json:"error"`
}

type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type AuthResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresIn    int       `json:"expires_in"`
}
