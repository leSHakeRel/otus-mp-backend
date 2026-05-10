package handlers

import (
	"net/http"

	"movie-night-planner-backend/internal/services"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VoteHandler struct {
	voteService *services.VoteService
}

func NewVoteHandler(voteService *services.VoteService) *VoteHandler {
	return &VoteHandler{voteService: voteService}
}

// CreateVote godoc
// @Summary Vote for film
// @Description Vote for a movie in an evening
// @Tags votes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param eveningId path string true "Evening ID"
// @Param request body services.CreateVoteInput true "Vote data"
// @Success 201 {object} models.Vote
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /api/evenings/{eveningId}/votes [post]
func (h *VoteHandler) CreateVote(c *gin.Context) {
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

	var input services.CreateVoteInput
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
	result, err := h.voteService.Create(eveningID, input, userID.(uuid.UUID))
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
			if appErr.Code == "FILM_NOT_FOUND" {
				c.JSON(http.StatusNotFound, response.ErrorResponse{
					Error: &response.AppError{
						Code:    appErr.Code,
						Message: appErr.Message,
					},
				})
				return
			}
			if appErr.Code == "ALREADY_VOTED" {
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

// GetVotesForEvening godoc
// @Summary Get votes for evening
// @Description Get all votes and vote summaries for an evening
// @Tags votes
// @Produce json
// @Param eveningId path string true "Evening ID"
// @Success 200 {array} response.VoteSummary
// @Failure 404 {object} response.ErrorResponse
// @Router /api/evenings/{eveningId}/votes [get]
func (h *VoteHandler) GetVotesForEvening(c *gin.Context) {
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

	result, err := h.voteService.GetVotesForEvening(eveningID)
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
