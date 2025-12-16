package handlers

import (
	"net/http"
	"portfolio-server/internal/middleware"
	"portfolio-server/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	// @Summary 사용자 회원가입
	// @Description 새로운 사용자를 등록합니다
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param request body services.RegisterRequest true "회원가입 요청"
	// @Success 201 {object} services.AuthResponse
	// @Failure 400 {object} map[string]interface{}
	// @Failure 409 {object} map[string]interface{}
	// @Router /auth/register [post]
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request body",
			"detail":  err.Error(),
		})
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	// @Summary 사용자 로그인
	// @Description 사용자 인증정보로 로그인합니다
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param request body services.LoginRequest true "로그인 요청"
	// @Success 200 {object} services.AuthResponse
	// @Failure 400 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Router /auth/login [post]
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request body",
			"detail":  err.Error(),
		})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	// @Summary 프로필 조회
	// @Description 현재 로그인한 사용자의 프로필을 조회합니다
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Security Bearer
	// @Success 200 {object} models.User
	// @Failure 401 {object} map[string]interface{}
	// @Router /auth/profile [get]
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}
