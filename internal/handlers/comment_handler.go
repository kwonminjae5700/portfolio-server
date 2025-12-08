package handlers

import (
	"net/http"
	"portfolio-server/internal/middleware"
	"portfolio-server/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		commentService: services.NewCommentService(),
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req services.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request body",
			"detail":  err.Error(),
		})
		return
	}

	comment, err := h.commentService.CreateComment(&req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetCommentsByArticle(c *gin.Context) {
	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid article ID",
			"detail":  err.Error(),
		})
		return
	}

	lastIDStr := c.Query("last_id")
	limitStr := c.DefaultQuery("limit", "50")

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

	comments, err := h.commentService.GetCommentsByArticle(uint(articleID), lastID, limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
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
			"message": "Invalid comment ID",
			"detail":  err.Error(),
		})
		return
	}

	var req services.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid request body",
			"detail":  err.Error(),
		})
		return
	}

	comment, err := h.commentService.UpdateComment(uint(id), &req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
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
			"message": "Invalid comment ID",
			"detail":  err.Error(),
		})
		return
	}

	if err := h.commentService.DeleteComment(uint(id), userID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
