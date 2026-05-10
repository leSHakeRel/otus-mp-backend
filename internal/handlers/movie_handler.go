package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"movie-night-planner-backend/internal/services"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MovieHandler struct {
	movieService *services.MovieService
}

func NewMovieHandler(movieService *services.MovieService) *MovieHandler {
	return &MovieHandler{movieService: movieService}
}

// SearchMovies godoc
// @Summary Search movies
// @Description Search movies in TMDB database
// @Tags movies
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number"
// @Success 200 {object} response.PaginatedResponse[response.EveningFilmResponse]
// @Failure 500 {object} response.ErrorResponse
// @Router /api/movies/search [get]
func (h *MovieHandler) SearchMovies(c *gin.Context) {
	var input services.SearchMoviesInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_INPUT",
				Message: "Invalid input: " + err.Error(),
			},
		})
		return
	}

	result, err := h.movieService.SearchMovies(input)
	if err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: &response.AppError{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetMovieDetails godoc
// @Summary Get movie details
// @Description Get detailed information about a movie from TMDB
// @Tags movies
// @Produce json
// @Param tmdbId path int true "TMDB Movie ID"
// @Success 200 {object} response.EveningFilmResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/movies/{tmdbId} [get]
func (h *MovieHandler) GetMovieDetails(c *gin.Context) {
	tmdbID, err := strconv.Atoi(c.Param("tmdbId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid TMDB ID",
			},
		})
		return
	}

	result, err := h.movieService.GetMovieDetails(tmdbID)
	if err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: &response.AppError{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AddFilmToEvening godoc
// @Summary Add film to evening
// @Description Add a movie to a specific evening
// @Tags movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param eveningId path string true "Evening ID"
// @Param request body services.AddFilmToEveningInput true "Film data"
// @Success 201 {object} models.EveningFilm
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /api/evenings/{eveningId}/movies [post]
func (h *MovieHandler) AddFilmToEvening(c *gin.Context) {
	eveningID, err := uuid.Parse(c.Param("eveningId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid evening ID",
			},
		})
		return
	}

	var input services.AddFilmToEveningInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_INPUT",
				Message: "Invalid input: " + err.Error(),
			},
		})
		return
	}

	userID, _ := c.Get("userID")
	result, err := h.movieService.AddFilmToEvening(eveningID, input.TMDBID, userID.(uuid.UUID))
	if err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			if appErr.Code == "EVENING_NOT_FOUND" {
				c.JSON(http.StatusNotFound, response.ErrorResponse{
					Error: &response.AppError{
						Code:    appErr.Code,
						Message: appErr.Message,
					},
				})
				return
			}
			if appErr.Code == "FILM_EXISTS" {
				c.JSON(http.StatusConflict, response.ErrorResponse{
					Error: &response.AppError{
						Code:    appErr.Code,
						Message: appErr.Message,
					},
				})
				return
			}
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: &response.AppError{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// RemoveFilmFromEvening godoc
// @Summary Remove film from evening
// @Description Remove a movie from a specific evening
// @Tags movies
// @Security BearerAuth
// @Param eveningId path string true "Evening ID"
// @Param tmdbId path int true "TMDB Movie ID"
// @Success 204
// @Failure 404 {object} response.ErrorResponse
// @Router /api/evenings/{eveningId}/movies/{tmdbId} [delete]
func (h *MovieHandler) RemoveFilmFromEvening(c *gin.Context) {
	eveningID, err := uuid.Parse(c.Param("eveningId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid evening ID",
			},
		})
		return
	}

	tmdbID, err := strconv.Atoi(c.Param("tmdbId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid TMDB ID",
			},
		})
		return
	}

	userID, _ := c.Get("userID")
	err = h.movieService.RemoveFilmFromEvening(eveningID, tmdbID, userID.(uuid.UUID))
	if err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			if appErr.Code == "EVENING_NOT_FOUND" || appErr.Code == "FILM_NOT_FOUND" {
				c.JSON(http.StatusNotFound, response.ErrorResponse{
					Error: &response.AppError{
						Code:    appErr.Code,
						Message: appErr.Message,
					},
				})
				return
			}
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: &response.AppError{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetFilmsForEvening godoc
// @Summary Get films for evening
// @Description Get all movies added to a specific evening
// @Tags movies
// @Produce json
// @Param eveningId path string true "Evening ID"
// @Success 200 {array} response.EveningFilmResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/evenings/{eveningId}/movies [get]
func (h *MovieHandler) GetFilmsForEvening(c *gin.Context) {
	eveningID, err := uuid.Parse(c.Param("eveningId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid evening ID",
			},
		})
		return
	}

	result, err := h.movieService.GetFilmsForEvening(eveningID)
	if err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(http.StatusNotFound, response.ErrorResponse{
				Error: &response.AppError{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
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
