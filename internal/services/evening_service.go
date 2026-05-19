package services

import (
	"movie-night-planner-backend/internal/models"
	"movie-night-planner-backend/internal/repositories"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"

	"time"

	"github.com/google/uuid"
)

type EveningService struct {
	eveningRepo     *repositories.EveningRepository
	eveningFilmRepo *repositories.EveningFilmRepository
	voteRepo        *repositories.VoteRepository
	commentRepo     *repositories.CommentRepository
	userRepo        *repositories.UserRepository
}

type CreateEveningInput struct {
	Title       string     `json:"title" validate:"required,min=3,max=255"`
	Description string     `json:"description" validate:"max=1000"`
	ScheduledAt *time.Time `json:"scheduledAt"`
	IsPrivate   bool       `json:"isPrivate"`
}

type UpdateEveningInput struct {
	Title       string     `json:"title,omitempty" validate:"max=255"`
	Description string     `json:"description,omitempty" validate:"max=1000"`
	ScheduledAt *time.Time `json:"scheduledAt"`
	IsPrivate   *bool      `json:"isPrivate"`
}

func NewEveningService(
	eveningRepo *repositories.EveningRepository,
	eveningFilmRepo *repositories.EveningFilmRepository,
	voteRepo *repositories.VoteRepository,
	commentRepo *repositories.CommentRepository,
	userRepo *repositories.UserRepository,
) *EveningService {
	return &EveningService{
		eveningRepo:     eveningRepo,
		eveningFilmRepo: eveningFilmRepo,
		voteRepo:        voteRepo,
		commentRepo:     commentRepo,
		userRepo:        userRepo,
	}
}

func (s *EveningService) Create(input CreateEveningInput, ownerID uuid.UUID) (*models.Evening, error) {
	// Verify owner exists
	_, err := s.userRepo.FindByID(ownerID)
	if err != nil {
		return nil, utils.NewAppError("USER_NOT_FOUND", "Owner not found", nil, nil)
	}

	evening := &models.Evening{
		Title:       input.Title,
		Description: input.Description,
		ScheduledAt: input.ScheduledAt,
		OwnerID:     ownerID,
		IsPrivate:   input.IsPrivate,
	}

	err = s.eveningRepo.Create(evening)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to create evening")
	}

	return evening, nil
}

func (s *EveningService) FindByID(id uuid.UUID) (*response.EveningResponse, error) {
	evening, err := s.eveningRepo.FindByID(id)
	if err != nil {
		return nil, utils.WrapError(err, "EVENING_NOT_FOUND", "Evening not found")
	}

	resp := s.mapEveningToResponse(evening)

	// Load movies
	films, err := s.eveningFilmRepo.FindByEveningID(id)
	if err == nil {
		for _, film := range films {
			movieResp := response.EveningFilmResponse{
				ID:           film.ID,
				TMDBID:       film.TMDBID,
				Title:        film.Title,
				PosterPath:   film.PosterPath,
				BackdropPath: film.BackdropPath,
				VoteAverage:  film.VoteAverage,
				Overview:     film.Overview,
				AddedAt:      film.AddedAt,
			}
			if film.ReleaseDate != nil {
				movieResp.ReleaseDate = film.ReleaseDate.Format("2006-01-02")
			}
			resp.Movies = append(resp.Movies, movieResp)
		}
	}

	// Load votes
	votes, err := s.voteRepo.FindByEveningID(id)
	if err == nil {
		for _, vote := range votes {
			vr := response.VoteResponse{
				ID:            vote.ID,
				EveningFilmID: vote.EveningFilmID,
				UserID:        vote.UserID,
				Value:         vote.Value,
				CreatedAt:     vote.CreatedAt,
			}
			// Load user info if available
			user, userErr := s.userRepo.FindByID(vote.UserID)
			if userErr == nil {
				vr.User = response.UserResponse{
					ID:        user.ID,
					Email:     user.Email,
					Username:  user.Username,
					CreatedAt: user.CreatedAt,
				}
			}
			resp.Votes = append(resp.Votes, vr)
		}
	}

	// Load comments
	comments, err := s.commentRepo.FindByEveningID(id)
	if err == nil {
		for _, comment := range comments {
			cr := response.CommentResponse{
				ID:        comment.ID,
				EveningID: comment.EveningID,
				UserID:    comment.UserID,
				Content:   comment.Content,
				CreatedAt: comment.CreatedAt,
			}
			if comment.User.ID != uuid.Nil {
				cr.User = response.UserResponse{
					ID:        comment.User.ID,
					Email:     comment.User.Email,
					Username:  comment.User.Username,
					CreatedAt: comment.User.CreatedAt,
				}
				cr.Username = comment.User.Username
			}
			resp.Comments = append(resp.Comments, cr)
		}
	}

	return resp, nil
}

func (s *EveningService) FindAll(page, limit int, isPrivate *bool, createdBy *uuid.UUID) (*response.PaginatedResponse[response.EveningResponse], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	evenings, total, err := s.eveningRepo.FindAll(page, limit, isPrivate, createdBy)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to fetch evenings")
	}

	data := make([]response.EveningResponse, len(evenings))
	for i, evening := range evenings {
		data[i] = *s.mapEveningToResponse(&evening)
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &response.PaginatedResponse[response.EveningResponse]{
		Data: data,
		Pagination: response.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *EveningService) FindByOwner(ownerID uuid.UUID, page, limit int) (*response.PaginatedResponse[response.EveningResponse], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	evenings, total, err := s.eveningRepo.FindByOwnerID(ownerID, page, limit)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to fetch evenings")
	}

	data := make([]response.EveningResponse, len(evenings))
	for i, evening := range evenings {
		data[i] = *s.mapEveningToResponse(&evening)
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &response.PaginatedResponse[response.EveningResponse]{
		Data: data,
		Pagination: response.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *EveningService) Update(id uuid.UUID, input UpdateEveningInput, ownerID uuid.UUID) (*models.Evening, error) {
	evening, err := s.eveningRepo.FindByID(id)
	if err != nil {
		return nil, utils.WrapError(err, "EVENING_NOT_FOUND", "Evening not found")
	}

	// Verify ownership
	if evening.OwnerID != ownerID {
		return nil, utils.NewAppError("FORBIDDEN", "You can only update your own evenings", nil, nil)
	}

	if input.Title != "" {
		evening.Title = input.Title
	}
	if input.Description != "" {
		evening.Description = input.Description
	}
	if input.ScheduledAt != nil {
		evening.ScheduledAt = input.ScheduledAt
	}
	if input.IsPrivate != nil {
		evening.IsPrivate = *input.IsPrivate
	}

	err = s.eveningRepo.Update(evening)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to update evening")
	}

	return evening, nil
}

func (s *EveningService) Delete(id uuid.UUID, ownerID uuid.UUID) error {
	evening, err := s.eveningRepo.FindByID(id)
	if err != nil {
		return utils.WrapError(err, "EVENING_NOT_FOUND", "Evening not found")
	}

	// Verify ownership
	if evening.OwnerID != ownerID {
		return utils.NewAppError("FORBIDDEN", "You can only delete your own evenings", nil, nil)
	}

	// Delete related films and votes
	_ = s.eveningFilmRepo.DeleteByEveningID(id)
	_ = s.eveningRepo.Delete(id)

	return nil
}

func (s *EveningService) mapEveningToResponse(evening *models.Evening) *response.EveningResponse {
	resp := &response.EveningResponse{
		ID:          evening.ID,
		Title:       evening.Title,
		Description: evening.Description,
		ScheduledAt: evening.ScheduledAt,
		IsPrivate:   evening.IsPrivate,
		CreatedAt:   evening.CreatedAt,
		UpdatedAt:   evening.UpdatedAt,
		Movies:      make([]response.EveningFilmResponse, 0),
		Votes:       make([]response.VoteResponse, 0),
		Comments:    make([]response.CommentResponse, 0),
	}

	if evening.Owner.ID != uuid.Nil {
		resp.CreatedBy = response.UserResponse{
			ID:        evening.Owner.ID,
			Email:     evening.Owner.Email,
			Username:  evening.Owner.Username,
			CreatedAt: evening.Owner.CreatedAt,
		}
	}

	return resp
}
