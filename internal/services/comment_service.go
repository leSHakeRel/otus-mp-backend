package services

import (
	"movie-night-planner-backend/internal/models"
	"movie-night-planner-backend/internal/repositories"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"

	"github.com/google/uuid"
)

type CommentService struct {
	commentRepo *repositories.CommentRepository
	eveningRepo *repositories.EveningRepository
	userRepo    *repositories.UserRepository
}

type CreateCommentInput struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

func NewCommentService(
	commentRepo *repositories.CommentRepository,
	eveningRepo *repositories.EveningRepository,
	userRepo *repositories.UserRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		eveningRepo: eveningRepo,
		userRepo:    userRepo,
	}
}

func (s *CommentService) Create(eveningID uuid.UUID, input CreateCommentInput, userID uuid.UUID) (*models.Comment, error) {
	// Verify evening exists
	_, err := s.eveningRepo.FindByID(eveningID)
	if err != nil {
		return nil, utils.WrapError(err, "EVENING_NOT_FOUND", "Evening not found")
	}

	// Create comment
	comment := &models.Comment{
		EveningID: eveningID,
		UserID:    userID,
		Content:   input.Content,
	}

	err = s.commentRepo.Create(comment)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to create comment")
	}

	return comment, nil
}

func (s *CommentService) GetCommentsForEvening(eveningID uuid.UUID) ([]response.CommentResponse, error) {
	// Verify evening exists
	_, err := s.eveningRepo.FindByID(eveningID)
	if err != nil {
		return nil, utils.WrapError(err, "EVENING_NOT_FOUND", "Evening not found")
	}

	comments, err := s.commentRepo.FindByEveningID(eveningID)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to get comments")
	}

	data := make([]response.CommentResponse, len(comments))
	for i, comment := range comments {
		data[i] = response.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		}
		if comment.User.ID != uuid.Nil {
			data[i].User = response.UserResponse{
				ID:        comment.User.ID,
				Email:     comment.User.Email,
				Username:  comment.User.Username,
				CreatedAt: comment.User.CreatedAt,
			}
		}
	}

	return data, nil
}

func (s *CommentService) Update(commentID uuid.UUID, input CreateCommentInput, userID uuid.UUID) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return nil, utils.WrapError(err, "COMMENT_NOT_FOUND", "Comment not found")
	}

	// Verify ownership
	if comment.UserID != userID {
		return nil, utils.NewAppError("FORBIDDEN", "You can only update your own comments", nil, nil)
	}

	comment.Content = input.Content
	err = s.commentRepo.Update(comment)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to update comment")
	}

	return comment, nil
}

func (s *CommentService) Delete(commentID uuid.UUID, userID uuid.UUID) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return utils.WrapError(err, "COMMENT_NOT_FOUND", "Comment not found")
	}

	// Verify ownership
	if comment.UserID != userID {
		return utils.NewAppError("FORBIDDEN", "You can only delete your own comments", nil, nil)
	}

	err = s.commentRepo.Delete(commentID)
	if err != nil {
		return utils.WrapError(err, "DATABASE_ERROR", "Failed to delete comment")
	}

	return nil
}
