package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type EveningResponse struct {
	ID          uuid.UUID             `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	ScheduledAt *time.Time            `json:"scheduledAt,omitempty"`
	CreatedBy   UserResponse          `json:"createdBy"`
	IsPrivate   bool                  `json:"isPrivate"`
	Movies      []EveningFilmResponse `json:"movies"`
	Votes       []VoteResponse        `json:"votes"`
	Comments    []CommentResponse     `json:"comments"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}

type EveningFilmResponse struct {
	ID           uuid.UUID `json:"id"`
	TMDBID       int       `json:"tmdbId"`
	Title        string    `json:"title"`
	PosterPath   string    `json:"posterPath,omitempty"`
	BackdropPath string    `json:"backdropPath,omitempty"`
	ReleaseDate  string    `json:"releaseDate,omitempty"`
	VoteAverage  float64   `json:"voteAverage,omitempty"`
	VoteCount    int       `json:"voteCount"`
	Overview     string    `json:"overview,omitempty"`
	GenreIDs     []int     `json:"genreIds,omitempty"`
	AddedAt      time.Time `json:"addedAt"`
}

type VoteResponse struct {
	ID            uuid.UUID    `json:"id"`
	EveningFilmID uuid.UUID    `json:"eveningFilmId"`
	UserID        uuid.UUID    `json:"userId"`
	User          UserResponse `json:"user"`
	Value         int          `json:"value"`
	CreatedAt     time.Time    `json:"createdAt"`
}

type VoteSummary struct {
	EveningFilmID    uuid.UUID        `json:"eveningFilmId"`
	Title            string           `json:"title"`
	PosterPath       string           `json:"posterPath,omitempty"`
	TotalVotes       int              `json:"totalVotes"`
	AverageScore     float64          `json:"averageScore"`
	VoteDistribution VoteDistribution `json:"voteDistribution"`
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
	EveningID uuid.UUID    `json:"eveningId"`
	UserID    uuid.UUID    `json:"userId"`
	Username  string       `json:"username"`
	User      UserResponse `json:"user"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"createdAt"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
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
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	ExpiresIn    int       `json:"expiresIn"`
}
