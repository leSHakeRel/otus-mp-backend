package handlers

import (
	"net/http"

	"movie-night-planner-backend/internal/services"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateComment godoc
// @Summary Create comment
// @Description Add a comment to an evening
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param eveningId path string true "Evening ID"
// @Param request body services.CreateCommentInput true "Comment data"
// @Success 201 {object} models.Comment
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/evenings/{eveningId}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	eveningID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid evening ID",
			},
		})
		return
	}

	var input services.CreateCommentInput
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
	result, err := h.commentService.Create(eveningID, input, userID.(uuid.UUID))
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

// GetCommentsForEvening godoc
// @Summary Get comments for evening
// @Description Get all comments for an evening
// @Tags comments
// @Produce json
// @Param eveningId path string true "Evening ID"
// @Success 200 {array} response.CommentResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/evenings/{eveningId}/comments [get]
func (h *CommentHandler) GetCommentsForEvening(c *gin.Context) {
	eveningID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: &response.AppError{
				Code:    "INVALID_ID",
				Message: "Invalid evening ID",
			},
		})
		return
	}

	result, err := h.commentService.GetCommentsForEvening(eveningID)
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
