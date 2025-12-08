package services

import (
	"fmt"
	"portfolio-server/internal/database"
	"portfolio-server/internal/errors"
	"portfolio-server/internal/middleware"
	"portfolio-server/internal/models"
	"portfolio-server/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{
		db: database.GetDB(),
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string           `json:"token"`
	User  models.User      `json:"user"`
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.ErrEmailAlreadyExists()
	}

	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.ErrUsernameAlreadyExists()
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to hash password", err.Error())
	}

	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.NewAppError(500, "Failed to create user", err.Error())
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to generate token", err.Error())
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrInvalidCredentials()
		}
		return nil, errors.NewAppError(500, "Database error", err.Error())
	}

	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, errors.ErrInvalidCredentials()
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to generate token", err.Error())
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound()
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
