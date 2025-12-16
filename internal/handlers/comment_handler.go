package handlers

import (
	"net/http"
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
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid article ID",
			"detail":  err.Error(),
		})
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

	userID := c.GetUint("user_id")
	comment, err := h.commentService.CreateComment(uint(articleID), &req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid article ID",
			"detail":  err.Error(),
		})
		return
	}

	var lastID *uint
	if lastIDStr := c.Query("last_id"); lastIDStr != "" {
		if id, err := strconv.ParseUint(lastIDStr, 10, 32); err == nil {
			uid := uint(id)
			lastID = &uid
		}
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	comments, err := h.commentService.GetCommentsByArticleID(uint(articleID), lastID, limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
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

	userID := c.GetUint("user_id")
	comment, err := h.commentService.UpdateComment(uint(commentID), &req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid comment ID",
			"detail":  err.Error(),
		})
		return
	}

	userID := c.GetUint("user_id")
	if err := h.commentService.DeleteComment(uint(commentID), userID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
