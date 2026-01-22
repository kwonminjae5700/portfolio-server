package models

import (
	"time"

	"gorm.io/gorm"
)

type VerificationCode struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"index;not null" json:"email"`
	Code      string         `gorm:"not null" json:"-"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	IsUsed    bool           `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (VerificationCode) TableName() string {
	return "verification_codes"
}
