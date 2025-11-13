package models

import (
	"time"

	"gorm.io/gorm"
)

type AuthUser struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	Email             string         `json:"email" gorm:"uniqueIndex;not null"`
	Password          string         `json:"-" gorm:"not null"`
	LastAccessTokenAt time.Time      `json:"lastAccessTokenAt" gorm:"default:null"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}
