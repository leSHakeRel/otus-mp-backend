package services

import (
	"movie-night-planner-backend/internal/models"
	"movie-night-planner-backend/internal/repositories"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo   *repositories.UserRepository
	jwtService *utils.JWTService
}

type AuthInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Username string `json:"username" validate:"required,min=3,max=100"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	ExpiresIn    int       `json:"expiresIn"`
}

func NewAuthService(userRepo *repositories.UserRepository, jwtService *utils.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(input AuthInput) (*AuthResponse, error) {
	// Check if email already exists
	_, err := s.userRepo.FindByEmail(input.Email)
	if err == nil {
		return nil, utils.NewAppError("EMAIL_EXISTS", "Email already registered", nil, nil)
	}

	// Check if username already exists
	_, err = s.userRepo.FindByUsername(input.Username)
	if err == nil {
		return nil, utils.NewAppError("USERNAME_EXISTS", "Username already taken", nil, nil)
	}

	// Hash password
	passwordHash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, utils.WrapError(err, "PASSWORD_HASH_ERROR", "Failed to hash password")
	}

	// Create user
	user := &models.User{
		Email:        input.Email,
		PasswordHash: passwordHash,
		Username:     input.Username,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to create user")
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, utils.WrapError(err, "TOKEN_GENERATION_ERROR", "Failed to generate access token")
	}

	return &AuthResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		AccessToken: accessToken,
		ExpiresIn:   86400, // 24 hours
	}, nil
}

func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, utils.NewAppError("INVALID_CREDENTIALS", "Invalid email or password", nil, nil)
	}

	// Verify password
	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		return nil, utils.NewAppError("INVALID_CREDENTIALS", "Invalid email or password", nil, nil)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, utils.WrapError(err, "TOKEN_GENERATION_ERROR", "Failed to generate access token")
	}

	return &AuthResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		AccessToken: accessToken,
		ExpiresIn:   86400, // 24 hours
	}, nil
}

func (s *AuthService) GetCurrentUser(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, utils.WrapError(err, "USER_NOT_FOUND", "User not found")
	}
	return user, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	return s.jwtService.ValidateToken(tokenString)
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*response.UserResponse, error) {
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
