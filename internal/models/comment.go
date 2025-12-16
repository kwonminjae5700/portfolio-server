package models

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	AuthorID  uint      `json:"author_id" gorm:"not null"`
	Author    User      `json:"author" gorm:"foreignKey:AuthorID"`
	ArticleID uint      `json:"article_id" gorm:"not null"`
	Article   Article   `json:"-" gorm:"foreignKey:ArticleID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentResponse struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	AuthorID   uint      `json:"author_id"`
	AuthorName string    `json:"author_name"`
	ArticleID  uint      `json:"article_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CommentListResponse struct {
	Comments   []CommentResponse `json:"comments"`
	NextCursor *uint             `json:"next_cursor,omitempty"`
	HasMore    bool              `json:"has_more"`
}
