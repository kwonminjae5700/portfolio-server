package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ArticleID uint           `gorm:"not null;index" json:"article_id"`
	AuthorID  uint           `gorm:"not null;index" json:"author_id"`
	Content   string         `gorm:"not null;type:text" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Article Article `gorm:"foreignKey:ArticleID" json:"-"`
	Author  User    `gorm:"foreignKey:AuthorID" json:"author"`
}

func (Comment) TableName() string {
	return "comments"
}

type CommentResponse struct {
	ID         uint      `json:"id"`
	ArticleID  uint      `json:"article_id"`
	AuthorID   uint      `json:"author_id"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CommentListResponse struct {
	Comments   []CommentResponse `json:"comments"`
	NextCursor *uint             `json:"next_cursor"`
	HasMore    bool              `json:"has_more"`
}
