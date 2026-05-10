package handlers

import (
	"net/http"

	"movie-night-planner-backend/internal/services"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.AuthInput true "Registration data"
// @Success 201 {object} response.AuthResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input services.AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_INPUT",
				Message: "Invalid input: " + err.Error(),
			},
		})
		return
	}

	result, err := h.authService.Register(input)
	if err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(http.StatusConflict, response.ErrorResponse{
				Error: &response.AppError{
					Code:    appErr.Code,
					Message: appErr.Message,
					Details: appErr.Details,
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

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.LoginInput true "Login credentials"
// @Success 200 {object} response.AuthResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_INPUT",
				Message: "Invalid input: " + err.Error(),
			},
		})
		return
	}

	result, err := h.authService.Login(input)
	if err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{
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

// GetMe godoc
// @Summary Get current user
// @Description Get current authenticated user info
// @Tags auth
// @Produce json
// @Success 200 {object} response.UserResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")

	result, err := h.authService.GetUserByID(userID.(uuid.UUID))
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
