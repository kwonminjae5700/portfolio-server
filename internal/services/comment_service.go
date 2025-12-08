package services

import (
	"fmt"
	"portfolio-server/internal/database"
	"portfolio-server/internal/errors"
	"portfolio-server/internal/models"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService() *CommentService {
	return &CommentService{
		db: database.GetDB(),
	}
}

type CreateCommentRequest struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

func (s *CommentService) CreateComment(req *CreateCommentRequest, authorID uint) (*models.Comment, error) {
	var article models.Article
	if err := s.db.First(&article, req.ArticleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrArticleNotFound()
		}
		return nil, fmt.Errorf("failed to verify article: %w", err)
	}

	comment := models.Comment{
		ArticleID: req.ArticleID,
		AuthorID:  authorID,
		Content:   req.Content,
	}

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Preload author
	if err := s.db.Preload("Author").First(&comment, comment.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load comment: %w", err)
	}

	return &comment, nil
}

func (s *CommentService) GetCommentsByArticle(articleID uint, lastID *uint, limit int) (*models.CommentListResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	query := s.db.Model(&models.Comment{}).
		Preload("Author").
		Where("article_id = ?", articleID)

	if lastID != nil && *lastID > 0 {
		query = query.Where("id < ?", *lastID)
	}

	query = query.Order("id DESC").Limit(limit + 1)

	var comments []models.Comment
	if err := query.Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	hasMore := len(comments) > limit
	if hasMore {
		comments = comments[:limit]
	}

	responses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = models.CommentResponse{
			ID:         comment.ID,
			ArticleID:  comment.ArticleID,
			AuthorID:   comment.AuthorID,
			AuthorName: comment.Author.Username,
			Content:    comment.Content,
			CreatedAt:  comment.CreatedAt,
			UpdatedAt:  comment.UpdatedAt,
		}
	}

	var nextCursor *uint
	if hasMore && len(comments) > 0 {
		lastCommentID := comments[len(comments)-1].ID
		nextCursor = &lastCommentID
	}

	return &models.CommentListResponse{
		Comments:   responses,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *CommentService) UpdateComment(id uint, req *UpdateCommentRequest, userID uint) (*models.Comment, error) {
	var comment models.Comment
	if err := s.db.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCommentNotFound()
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	if comment.AuthorID != userID {
		return nil, errors.ErrPermissionDenied()
	}

	comment.Content = req.Content
	if err := s.db.Save(&comment).Error; err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	// Preload author
	if err := s.db.Preload("Author").First(&comment, comment.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load comment: %w", err)
	}

	return &comment, nil
}

func (s *CommentService) DeleteComment(id uint, userID uint) error {
	var comment models.Comment
	if err := s.db.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrCommentNotFound()
		}
		return fmt.Errorf("failed to get comment: %w", err)
	}

	// Check if user is the author
	if comment.AuthorID != userID {
		return errors.ErrPermissionDenied()
	}

	if err := s.db.Delete(&comment).Error; err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}
