package services

import (
	"context"
	"fmt"
	"portfolio-server/internal/database"
	"portfolio-server/internal/errors"
	"portfolio-server/internal/middleware"
	"portfolio-server/internal/models"
	"portfolio-server/internal/utils"
	"time"

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
		return nil, errors.NewAppError(500, "비밀번호 암호화에 실패했습니다", err.Error())
	}

	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.NewAppError(500, "사용자 생성에 실패했습니다", err.Error())
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
		return nil, errors.NewAppError(500, "데이터베이스 오류가 발생했습니다", err.Error())
	}

	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, errors.ErrInvalidCredentials()
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, errors.NewAppError(500, "토큰 생성에 실패했습니다", err.Error())
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

// SendVerificationCodeRequest is the request to send verification code
type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyCodeRequest is the request to verify the code
type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// SendVerificationCode generates and sends a verification code to the user's email
func (s *AuthService) SendVerificationCode(email string) error {
	// 이메일 중복 체크
	var existingUser models.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return errors.ErrEmailAlreadyExists()
	}

	// 6자리 인증 코드 생성
	code, err := utils.GenerateVerificationCode()
	if err != nil {
		return errors.NewAppError(500, "인증 코드 생성에 실패했습니다", err.Error())
	}

	// Redis에 인증 코드 저장 (10분 TTL)
	ctx := context.Background()
	redisKey := fmt.Sprintf("verification:%s", email)
	err = database.GetRedis().Set(ctx, redisKey, code, 10*time.Minute).Err()
	if err != nil {
		return errors.NewAppError(500, "인증 코드 저장에 실패했습니다", err.Error())
	}

	// 이메일 전송
	if err := utils.SendVerificationEmail(email, code); err != nil {
		return errors.NewAppError(500, "이메일 전송에 실패했습니다", err.Error())
	}

	return nil
}

// VerifyCode verifies the verification code
func (s *AuthService) VerifyCode(email, code string) error {
	ctx := context.Background()
	redisKey := fmt.Sprintf("verification:%s", email)

	// Redis에서 저장된 코드 조회
	savedCode, err := database.GetRedis().Get(ctx, redisKey).Result()
	if err != nil {
		return errors.NewAppError(400, "인증 코드가 유효하지 않습니다", "코드를 찾을 수 없거나 만료되었습니다")
	}

	// 코드 비교
	if savedCode != code {
		return errors.NewAppError(400, "인증 코드가 유효하지 않습니다", "코드가 일치하지 않습니다")
	}

	// 인증 성공 시 Redis에서 삭제
	database.GetRedis().Del(ctx, redisKey)

	return nil
}
