package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"not null;type:varchar(200)" json:"title"`
	Content   string         `gorm:"not null;type:text" json:"content"`
	AuthorID  uint           `gorm:"not null;index" json:"author_id"`
	ViewCount int            `gorm:"default:0" json:"view_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Author     User       `gorm:"foreignKey:AuthorID" json:"author"`
	Categories []Category `gorm:"many2many:article_categories;" json:"-"`
}

func (Article) TableName() string {
	return "articles"
}

type CategoryInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ArticleResponse struct {
	ID         uint           `json:"id"`
	Title      string         `json:"title"`
	Content    string         `json:"content"`
	AuthorID   uint           `json:"author_id"`
	AuthorName string         `json:"author_name"`
	ViewCount  int            `json:"view_count"`
	Categories []CategoryInfo `json:"categories"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type ArticleListResponse struct {
	Articles   []ArticleResponse `json:"articles"`
	NextCursor *uint             `json:"next_cursor"`
	HasMore    bool              `json:"has_more"`
}
