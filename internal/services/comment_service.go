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
	Content string `json:"content" binding:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

func (s *CommentService) CreateComment(articleID uint, req *CreateCommentRequest, authorID uint) (*models.CommentResponse, error) {
	// Check if article exists
	var article models.Article
	if err := s.db.First(&article, articleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrArticleNotFound()
		}
		return nil, fmt.Errorf("게시글 조회 실패: %w", err)
	}

	comment := models.Comment{
		Content:   req.Content,
		AuthorID:  authorID,
		ArticleID: articleID,
	}

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, fmt.Errorf("댓글 생성 실패: %w", err)
	}

	// Load author
	if err := s.db.Preload("Author").First(&comment, comment.ID).Error; err != nil {
		return nil, fmt.Errorf("댓글 로드 실패: %w", err)
	}

	return &models.CommentResponse{
		ID:         comment.ID,
		Content:    comment.Content,
		AuthorID:   comment.AuthorID,
		AuthorName: comment.Author.Username,
		ArticleID:  comment.ArticleID,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
	}, nil
}

func (s *CommentService) GetCommentsByArticleID(articleID uint, lastID *uint, limit int) (*models.CommentListResponse, error) {
	// Check if article exists
	var article models.Article
	if err := s.db.First(&article, articleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrArticleNotFound()
		}
		return nil, fmt.Errorf("게시글 조회 실패: %w", err)
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	query := s.db.Model(&models.Comment{}).
		Where("article_id = ?", articleID).
		Preload("Author")

	if lastID != nil && *lastID > 0 {
		query = query.Where("id > ?", *lastID)
	}

	query = query.Order("id ASC").Limit(limit + 1)

	var comments []models.Comment
	if err := query.Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("댓글 목록 조회 실패: %w", err)
	}

	hasMore := len(comments) > limit
	if hasMore {
		comments = comments[:limit]
	}

	responses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = models.CommentResponse{
			ID:         comment.ID,
			Content:    comment.Content,
			AuthorID:   comment.AuthorID,
			AuthorName: comment.Author.Username,
			ArticleID:  comment.ArticleID,
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

func (s *CommentService) UpdateComment(commentID uint, req *UpdateCommentRequest, userID uint) (*models.CommentResponse, error) {
	var comment models.Comment
	if err := s.db.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCommentNotFound()
		}
		return nil, fmt.Errorf("댓글 조회 실패: %w", err)
	}

	if comment.AuthorID != userID {
		return nil, errors.ErrPermissionDenied()
	}

	comment.Content = req.Content

	if err := s.db.Save(&comment).Error; err != nil {
		return nil, fmt.Errorf("댓글 수정 실패: %w", err)
	}

	// Load author
	if err := s.db.Preload("Author").First(&comment, comment.ID).Error; err != nil {
		return nil, fmt.Errorf("댓글 로드 실패: %w", err)
	}

	return &models.CommentResponse{
		ID:         comment.ID,
		Content:    comment.Content,
		AuthorID:   comment.AuthorID,
		AuthorName: comment.Author.Username,
		ArticleID:  comment.ArticleID,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
	}, nil
}

func (s *CommentService) DeleteComment(commentID uint, userID uint) error {
	var comment models.Comment
	if err := s.db.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrCommentNotFound()
		}
		return fmt.Errorf("댓글 조회 실패: %w", err)
	}

	if comment.AuthorID != userID {
		return errors.ErrPermissionDenied()
	}

	if err := s.db.Delete(&comment).Error; err != nil {
		return fmt.Errorf("댓글 삭제 실패: %w", err)
	}

	return nil
}
