package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	UserID    uint           `json:"userId" gorm:"not null;index"`
	ExpiresAt time.Time      `json:"expiresAt" gorm:"not null"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
