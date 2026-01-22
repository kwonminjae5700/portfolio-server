package handlers

import (
	"net/http"
	"portfolio-server/internal/middleware"
	"portfolio-server/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService *services.ArticleService
}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{
		articleService: services.NewArticleService(),
	}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	// @Summary 글 작성
	// @Description 새로운 게시글을 작성합니다
	// @Tags articles
	// @Accept json
	// @Produce json
	// @Security Bearer
	// @Param request body services.CreateArticleRequest true "글 작성 요청"
	// @Success 201 {object} models.Article
	// @Failure 400 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Router /articles [post]
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req services.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "요청 형식이 올바르지 않습니다",
			"detail":  err.Error(),
		})
		return
	}

	article, err := h.articleService.CreateArticle(&req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	// @Summary 글 상세 조회
	// @Description 특정 게시글의 상세 정보를 조회합니다
	// @Tags articles
	// @Accept json
	// @Produce json
	// @Param id path uint true "글 ID"
	// @Success 200 {object} models.ArticleResponse
	// @Failure 400 {object} map[string]interface{}
	// @Failure 404 {object} map[string]interface{}
	// @Router /articles/{id} [get]
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid article ID",
			"detail":  err.Error(),
		})
		return
	}

	article, err := h.articleService.GetArticleByID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	// @Summary 글 목록 조회
	// @Description 커서 기반 무한 스크롤로 게시글 목록을 조회합니다
	// @Tags articles
	// @Accept json
	// @Produce json
	// @Param last_id query uint false "마지막 글 ID"
	// @Param limit query int false "조회할 개수 (기본값: 20)"
	// @Success 200 {object} models.ArticleListResponse
	// @Failure 400 {object} map[string]interface{}
	// @Router /articles [get]
	lastIDStr := c.Query("last_id")
	limitStr := c.DefaultQuery("limit", "20")

	var lastID *uint
	if lastIDStr != "" {
		id, err := strconv.ParseUint(lastIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "last_id 파라미터가 올바르지 않습니다",
				"detail":  err.Error(),
			})
			return
		}
		lastIDUint := uint(id)
		lastID = &lastIDUint
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "limit 파라미터가 올바르지 않습니다",
			"detail":  err.Error(),
		})
		return
	}

	articles, err := h.articleService.GetArticles(lastID, limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	// @Summary 글 수정
	// @Description 자신의 게시글을 수정합니다
	// @Tags articles
	// @Accept json
	// @Produce json
	// @Security Bearer
	// @Param id path uint true "글 ID"
	// @Param request body services.UpdateArticleRequest true "글 수정 요청"
	// @Success 200 {object} models.Article
	// @Failure 400 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Failure 403 {object} map[string]interface{}
	// @Failure 404 {object} map[string]interface{}
	// @Router /articles/{id} [put]
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "게시글 ID가 올바르지 않습니다",
			"detail":  err.Error(),
		})
		return
	}

	var req services.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "요청 형식이 올바르지 않습니다",
			"detail":  err.Error(),
		})
		return
	}

	article, err := h.articleService.UpdateArticle(uint(id), &req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	// @Summary 글 삭제
	// @Description 자신의 게시글을 삭제합니다
	// @Tags articles
	// @Accept json
	// @Produce json
	// @Security Bearer
	// @Param id path uint true "글 ID"
	// @Success 204
	// @Failure 400 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Failure 403 {object} map[string]interface{}
	// @Failure 404 {object} map[string]interface{}
	// @Router /articles/{id} [delete]
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid article ID",
			"detail":  err.Error(),
		})
		return
	}

	if err := h.articleService.DeleteArticle(uint(id), userID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ArticleHandler) GetTopArticles(c *gin.Context) {
	// @Summary 조회수 TOP 5 글
	// @Description 조회수가 가장 높은 5개의 글을 반환합니다
	// @Tags articles
	// @Accept json
	// @Produce json
	// @Success 200 {object} []models.TopArticleInfo
	// @Failure 500 {object} map[string]interface{}
	// @Router /articles/top/views [get]
	articles, err := h.articleService.GetTopArticlesByViewCount()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, articles)
}
