package services

import (
	"fmt"
	"portfolio-server/internal/database"
	"portfolio-server/internal/models"

	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		db: database.GetDB(),
	}
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

func (s *CategoryService) CreateCategory(req *CreateCategoryRequest) (*models.Category, error) {
	category := models.Category{
		Name: req.Name,
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := s.db.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

func (s *CategoryService) UpdateCategory(id uint, req *UpdateCategoryRequest) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	category.Name = req.Name
	if err := s.db.Save(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &category, nil
}

func (s *CategoryService) DeleteCategory(id uint) error {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("category not found")
		}
		return fmt.Errorf("failed to get category: %w", err)
	}

	// Delete category associations first
	if err := s.db.Exec("DELETE FROM article_categories WHERE category_id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete category associations: %w", err)
	}

	if err := s.db.Delete(&category).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
