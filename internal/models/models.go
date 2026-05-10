package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Username     string    `gorm:"type:varchar(100);not null" json:"username"`
	CreatedAt    time.Time `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`

	Evenings []Evening `gorm:"foreignKey:OwnerID" json:"evenings,omitempty"`
	Votes    []Vote    `gorm:"foreignKey:UserID" json:"votes,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type Evening struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	ScheduledAt *time.Time `gorm:"type:timestamp with time zone" json:"scheduled_at"`
	OwnerID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"owner_id"`
	IsPrivate   bool       `gorm:"type:boolean;default:false" json:"is_private"`
	CreatedAt   time.Time  `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`

	Owner        User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	EveningFilms []EveningFilm `gorm:"foreignKey:EveningID" json:"evening_films,omitempty"`
	Votes        []Vote        `gorm:"foreignKey:EveningID" json:"votes,omitempty"`
	Comments     []Comment     `gorm:"foreignKey:EveningID" json:"comments,omitempty"`
}

func (Evening) TableName() string {
	return "evenings"
}

func (e *Evening) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

type EveningFilm struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	EveningID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"evening_id"`
	TMDBID       int        `gorm:"type:integer;not null;index" json:"tmdb_id"`
	Title        string     `gorm:"type:varchar(255);not null" json:"title"`
	PosterPath   string     `gorm:"type:varchar(255)" json:"poster_path,omitempty"`
	BackdropPath string     `gorm:"type:varchar(255)" json:"backdrop_path,omitempty"`
	ReleaseDate  *time.Time `gorm:"type:date" json:"release_date,omitempty"`
	VoteAverage  float64    `gorm:"type:decimal(3,1)" json:"vote_average,omitempty"`
	Overview     string     `gorm:"type:text" json:"overview,omitempty"`
	AddedAt      time.Time  `gorm:"type:timestamp with time zone;autoCreateTime" json:"added_at"`

	Evening Evening `gorm:"foreignKey:EveningID" json:"evening,omitempty"`
	Votes   []Vote  `gorm:"foreignKey:EveningFilmID" json:"votes,omitempty"`
}

func (EveningFilm) TableName() string {
	return "evening_films"
}

func (ef *EveningFilm) BeforeCreate(tx *gorm.DB) error {
	if ef.ID == uuid.Nil {
		ef.ID = uuid.New()
	}
	return nil
}

type Vote struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	EveningID     uuid.UUID `gorm:"type:uuid;not null;index" json:"evening_id"`
	EveningFilmID uuid.UUID `gorm:"type:uuid;not null;index" json:"evening_film_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Value         int       `gorm:"type:integer;not null;check:value >= 1 AND value <= 5" json:"value"`
	CreatedAt     time.Time `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at"`

	Evening     Evening     `gorm:"foreignKey:EveningID" json:"evening,omitempty"`
	EveningFilm EveningFilm `gorm:"foreignKey:EveningFilmID" json:"evening_film,omitempty"`
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Vote) TableName() string {
	return "votes"
}

func (v *Vote) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	EveningID uuid.UUID `gorm:"type:uuid;not null;index" json:"evening_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`

	Evening User `gorm:"foreignKey:EveningID" json:"evening,omitempty"`
	User    User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
