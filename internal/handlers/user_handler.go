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

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get paginated list of users
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse[response.UserResponse]
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := h.userService.GetAllUsers(page, limit)
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

// GetUser godoc
// @Summary Get user by ID
// @Description Get details of a specific user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Success 200 {object} response.UserResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/users/{userId} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid user ID",
			},
		})
		return
	}

	result, err := h.userService.GetUserByID(userID)
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
