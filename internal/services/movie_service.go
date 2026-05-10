package services

import (
	"encoding/json"
	"io"
	"movie-night-planner-backend/internal/models"
	"movie-night-planner-backend/internal/repositories"
	"movie-night-planner-backend/internal/tmdb"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type MovieService struct {
	eveningFilmRepo *repositories.EveningFilmRepository
	eveningRepo     *repositories.EveningRepository
	tmdbClient      *tmdb.Client
}

type SearchMoviesInput struct {
	Query string `form:"q" validate:"required,min=2"`
	Page  int    `form:"page"`
}

type AddFilmToEveningInput struct {
	TMDBID int `json:"tmdb_id" validate:"required"`
}

func NewMovieService(
	eveningFilmRepo *repositories.EveningFilmRepository,
	eveningRepo *repositories.EveningRepository,
	tmdbClient *tmdb.Client,
) *MovieService {
	return &MovieService{
		eveningFilmRepo: eveningFilmRepo,
		eveningRepo:     eveningRepo,
		tmdbClient:      tmdbClient,
	}
}

func (s *MovieService) SearchMovies(input SearchMoviesInput) (*response.PaginatedResponse[response.EveningFilmResponse], error) {
	if input.Page < 1 {
		input.Page = 1
	}

	result, err := s.tmdbClient.SearchMovies(input.Query, input.Page)
	if err != nil {
		return nil, utils.WrapError(err, "TMDB_API_ERROR", "Failed to search movies from TMDB")
	}

	data := make([]response.EveningFilmResponse, len(result.Results))
	for i, movie := range result.Results {
		data[i] = response.EveningFilmResponse{
			ID:           uuid.Nil,
			TMDBID:       movie.ID,
			Title:        movie.Title,
			PosterPath:   movie.PosterPath,
			BackdropPath: movie.BackdropPath,
			ReleaseDate:  movie.ReleaseDate,
			VoteAverage:  movie.VoteAverage,
			Overview:     movie.Overview,
			AddedAt:      time.Now(),
		}
	}

	return &response.PaginatedResponse[response.EveningFilmResponse]{
		Data: data,
		Pagination: response.Pagination{
			Page:       input.Page,
			Limit:      20,
			Total:      int64(result.TotalResults),
			TotalPages: result.TotalPages,
		},
	}, nil
}

func (s *MovieService) GetMovieDetails(tmdbID int) (*response.EveningFilmResponse, error) {
	movie, err := s.tmdbClient.GetMovieDetails(tmdbID)
	if err != nil {
		return nil, utils.WrapError(err, "TMDB_API_ERROR", "Failed to get movie details from TMDB")
	}

	return &response.EveningFilmResponse{
		ID:           uuid.Nil,
		TMDBID:       movie.ID,
		Title:        movie.Title,
		PosterPath:   movie.PosterPath,
		BackdropPath: movie.BackdropPath,
		ReleaseDate:  movie.ReleaseDate,
		VoteAverage:  movie.VoteAverage,
		Overview:     movie.Overview,
		AddedAt:      time.Now(),
	}, nil
}

func (s *MovieService) AddFilmToEvening(eveningID uuid.UUID, tmdbID int, userID uuid.UUID) (*models.EveningFilm, error) {
	// Verify evening exists
	_, err := s.eveningRepo.FindByID(eveningID)
	if err != nil {
		return nil, utils.WrapError(err, "EVENING_NOT_FOUND", "Evening not found")
	}

	// Check if film already exists in evening
	_, err = s.eveningFilmRepo.FindByEveningIDAndTMDBID(eveningID, tmdbID)
	if err == nil {
		return nil, utils.NewAppError("FILM_EXISTS", "Film already added to this evening", nil, nil)
	}

	// Get movie details from TMDB
	movie, err := s.tmdbClient.GetMovieDetails(tmdbID)
	if err != nil {
		return nil, utils.WrapError(err, "TMDB_API_ERROR", "Failed to get movie details from TMDB")
	}

	// Create evening film
	eveningFilm := &models.EveningFilm{
		EveningID:    eveningID,
		TMDBID:       movie.ID,
		Title:        movie.Title,
		PosterPath:   movie.PosterPath,
		BackdropPath: movie.BackdropPath,
		ReleaseDate:  movie.ReleaseDate,
		VoteAverage:  movie.VoteAverage,
		Overview:     movie.Overview,
	}

	err = s.eveningFilmRepo.Create(eveningFilm)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to add film to evening")
	}

	return eveningFilm, nil
}

func (s *MovieService) RemoveFilmFromEvening(eveningID uuid.UUID, tmdbID int, userID uuid.UUID) error {
	// Verify evening exists
	_, err := s.eveningRepo.FindByID(eveningID)
	if err != nil {
		return utils.WrapError(err, "EVENING_NOT_FOUND", "Evening not found")
	}

	// Find and delete film
	eveningFilm, err := s.eveningFilmRepo.FindByEveningIDAndTMDBID(eveningID, tmdbID)
	if err != nil {
		return utils.NewAppError("FILM_NOT_FOUND", "Film not found in this evening", nil, nil)
	}

	err = s.eveningFilmRepo.Delete(eveningFilm.ID)
	if err != nil {
		return utils.WrapError(err, "DATABASE_ERROR", "Failed to remove film from evening")
	}

	return nil
}

func (s *MovieService) GetFilmsForEvening(eveningID uuid.UUID) ([]response.EveningFilmResponse, error) {
	eveningFilms, err := s.eveningFilmRepo.FindByEveningID(eveningID)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to get films for evening")
	}

	data := make([]response.EveningFilmResponse, len(eveningFilms))
	for i, film := range eveningFilms {
		data[i] = response.EveningFilmResponse{
			ID:           film.ID,
			TMDBID:       film.TMDBID,
			Title:        film.Title,
			PosterPath:   film.PosterPath,
			BackdropPath: film.BackdropPath,
			ReleaseDate:  film.ReleaseDate,
			VoteAverage:  film.VoteAverage,
			Overview:     film.Overview,
			AddedAt:      film.AddedAt,
		}
	}

	return data, nil
}

// Helper function to read JSON from HTTP response
func readJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
