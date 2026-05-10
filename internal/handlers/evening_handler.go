package handlers

import (
	"net/http"
	"strconv"

	"movie-night-planner-backend/internal/services"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EveningHandler struct {
	eveningService *services.EveningService
}

func NewEveningHandler(eveningService *services.EveningService) *EveningHandler {
	return &EveningHandler{eveningService: eveningService}
}

// CreateEvening godoc
// @Summary Create new evening
// @Description Create a new movie night event
// @Tags evenings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateEveningInput true "Evening data"
// @Success 201 {object} models.Evening
// @Failure 400 {object} response.ErrorResponse
// @Router /api/evenings [post]
func (h *EveningHandler) CreateEvening(c *gin.Context) {
	var input services.CreateEveningInput
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
	result, err := h.eveningService.Create(input, userID.(uuid.UUID))
	if err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			if appErr.Code == "USER_NOT_FOUND" {
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

	c.JSON(http.StatusCreated, result)
}

// GetEvening godoc
// @Summary Get evening by ID
// @Description Get details of a specific movie night event
// @Tags evenings
// @Produce json
// @Param id path string true "Evening ID"
// @Success 200 {object} response.EveningResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/evenings/{id} [get]
func (h *EveningHandler) GetEvening(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid evening ID",
			},
		})
		return
	}

	result, err := h.eveningService.FindByID(id)
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

// GetAllEvenings godoc
// @Summary Get all evenings
// @Description Get paginated list of movie night events
// @Tags evenings
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param is_private query bool false "Filter by privacy"
// @Success 200 {object} response.PaginatedResponse[response.EveningResponse]
// @Router /api/evenings [get]
func (h *EveningHandler) GetAllEvenings(c *gin.Context) {
	page, _ := c.GetQuery("page")
	limit, _ := c.GetQuery("limit")
	isPrivate, _ := c.GetQuery("is_private")

	var isPrivateBool *bool
	if isPrivate != "" {
		val := isPrivate == "true"
		isPrivateBool = &val
	}

	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	l, _ := strconv.Atoi(limit)
	if l < 1 {
		l = 10
	}

	result, err := h.eveningService.FindAll(p, l, isPrivateBool)
	if err != nil {
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

// UpdateEvening godoc
// @Summary Update evening
// @Description Update a movie night event
// @Tags evenings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Evening ID"
// @Param request body services.UpdateEveningInput true "Evening data"
// @Success 200 {object} models.Evening
// @Failure 404 {object} response.ErrorResponse
// @Router /api/evenings/{id} [put]
func (h *EveningHandler) UpdateEvening(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid evening ID",
			},
		})
		return
	}

	var input services.UpdateEveningInput
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
	result, err := h.eveningService.Update(id, input, userID.(uuid.UUID))
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
			if appErr.Code == "FORBIDDEN" {
				c.JSON(http.StatusForbidden, response.ErrorResponse{
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

	c.JSON(http.StatusOK, result)
}

// DeleteEvening godoc
// @Summary Delete evening
// @Description Delete a movie night event
// @Tags evenings
// @Security BearerAuth
// @Param id path string true "Evening ID"
// @Success 204
// @Failure 404 {object} response.ErrorResponse
// @Router /api/evenings/{id} [delete]
func (h *EveningHandler) DeleteEvening(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid evening ID",
			},
		})
		return
	}

	userID, _ := c.Get("userID")
	err = h.eveningService.Delete(id, userID.(uuid.UUID))
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
			if appErr.Code == "FORBIDDEN" {
				c.JSON(http.StatusForbidden, response.ErrorResponse{
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
