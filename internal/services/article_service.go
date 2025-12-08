package services

import (
	"fmt"
	"portfolio-server/internal/database"
	"portfolio-server/internal/errors"
	"portfolio-server/internal/models"

	"gorm.io/gorm"
)

type ArticleService struct {
	db *gorm.DB
}

func NewArticleService() *ArticleService {
	return &ArticleService{
		db: database.GetDB(),
	}
}

type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

type UpdateArticleRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

func (s *ArticleService) CreateArticle(req *CreateArticleRequest, authorID uint) (*models.Article, error) {
	article := models.Article{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: authorID,
	}

	if err := s.db.Create(&article).Error; err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	if err := s.db.Preload("Author").First(&article, article.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load article: %w", err)
	}

	return &article, nil
}

func (s *ArticleService) GetArticleByID(id uint) (*models.ArticleResponse, error) {
	var article models.Article
	if err := s.db.Preload("Author").First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrArticleNotFound()
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	s.db.Model(&article).Update("view_count", gorm.Expr("view_count + ?", 1))
	article.ViewCount++

	var commentCount int64
	s.db.Model(&models.Comment{}).Where("article_id = ?", id).Count(&commentCount)

	response := &models.ArticleResponse{
		ID:           article.ID,
		Title:        article.Title,
		Content:      article.Content,
		AuthorID:     article.AuthorID,
		AuthorName:   article.Author.Username,
		ViewCount:    article.ViewCount,
		CommentCount: int(commentCount),
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}

	return response, nil
}

func (s *ArticleService) GetArticles(lastID *uint, limit int) (*models.ArticleListResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	query := s.db.Model(&models.Article{}).Preload("Author")

	if lastID != nil && *lastID > 0 {
		query = query.Where("id < ?", *lastID)
	}

	query = query.Order("id DESC").Limit(limit + 1)

	var articles []models.Article
	if err := query.Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}

	hasMore := len(articles) > limit
	if hasMore {
		articles = articles[:limit]
	}

	articleIDs := make([]uint, len(articles))
	for i, article := range articles {
		articleIDs[i] = article.ID
	}

	type CommentCountResult struct {
		ArticleID uint
		Count     int64
	}

	var commentCounts []CommentCountResult
	if len(articleIDs) > 0 {
		s.db.Model(&models.Comment{}).
			Select("article_id, count(*) as count").
			Where("article_id IN ?", articleIDs).
			Group("article_id").
				Scan(&commentCounts)
	}

	commentCountMap := make(map[uint]int64)
	for _, cc := range commentCounts {
		commentCountMap[cc.ArticleID] = cc.Count
	}

	responses := make([]models.ArticleResponse, len(articles))
	for i, article := range articles {
		responses[i] = models.ArticleResponse{
			ID:           article.ID,
			Title:        article.Title,
			Content:      article.Content,
			AuthorID:     article.AuthorID,
			AuthorName:   article.Author.Username,
			ViewCount:    article.ViewCount,
			CommentCount: int(commentCountMap[article.ID]),
			CreatedAt:    article.CreatedAt,
			UpdatedAt:    article.UpdatedAt,
		}
	}

	var nextCursor *uint
	if hasMore && len(articles) > 0 {
		lastArticleID := articles[len(articles)-1].ID
		nextCursor = &lastArticleID
	}

	return &models.ArticleListResponse{
		Articles:   responses,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *ArticleService) UpdateArticle(id uint, req *UpdateArticleRequest, userID uint) (*models.Article, error) {
	var article models.Article
	if err := s.db.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrArticleNotFound()
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	if article.AuthorID != userID {
		return nil, errors.ErrPermissionDenied()
	}

	article.Title = req.Title
	article.Content = req.Content

	if err := s.db.Save(&article).Error; err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	// Preload author
	if err := s.db.Preload("Author").First(&article, article.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load article: %w", err)
	}

	return &article, nil
}

func (s *ArticleService) DeleteArticle(id uint, userID uint) error {
	var article models.Article
	if err := s.db.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrArticleNotFound()
		}
		return fmt.Errorf("failed to get article: %w", err)
	}

	// Check if user is the author
	if article.AuthorID != userID {
		return errors.ErrPermissionDenied()
	}

	if err := s.db.Delete(&article).Error; err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	if err := s.db.Where("article_id = ?", id).Delete(&models.Comment{}).Error; err != nil {
		return fmt.Errorf("failed to delete comments: %w", err)
	}

	return nil
}
