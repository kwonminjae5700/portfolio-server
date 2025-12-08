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
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req services.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request body",
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
	lastIDStr := c.Query("last_id")
	limitStr := c.DefaultQuery("limit", "20")

	var lastID *uint
	if lastIDStr != "" {
		id, err := strconv.ParseUint(lastIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Invalid last_id parameter",
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
			"message": "Invalid limit parameter",
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

	var req services.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request body",
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
