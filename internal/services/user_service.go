package services

import (
	"math"

	"movie-night-planner-backend/internal/repositories"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetAllUsers возвращает пагинированный список пользователей
func (s *UserService) GetAllUsers(page, limit int) (*response.PaginatedResponse[response.UserResponse], error) {
	users, total, err := s.userRepo.FindAll(page, limit)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to fetch users")
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	data := make([]response.UserResponse, len(users))
	for i, u := range users {
		data[i] = response.UserResponse{
			ID:        u.ID,
			Email:     u.Email,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
		}
	}

	return &response.PaginatedResponse[response.UserResponse]{
		Data: data,
		Pagination: response.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetUserByID возвращает пользователя по ID
func (s *UserService) GetUserByID(userID uuid.UUID) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, utils.WrapError(err, "USER_NOT_FOUND", "User not found")
	}

	return &response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}
