package models

import (
	"time"

	"gorm.io/gorm"
)

type Board struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	TeamID    uint           `json:"teamId" gorm:"not null;index"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Team      Team           `json:"team,omitempty" gorm:"-"`
	Columns   []Column       `json:"columns,omitempty" gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
}
