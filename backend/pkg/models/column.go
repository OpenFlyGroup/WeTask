package models

import (
	"time"

	"gorm.io/gorm"
)

type Column struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	Order     int            `json:"order" gorm:"not null"`
	BoardID   uint           `json:"boardId" gorm:"not null;index"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Board     Board          `json:"board,omitempty" gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
	Tasks     []Task         `json:"tasks,omitempty" gorm:"foreignKey:ColumnID;constraint:OnDelete:CASCADE"`
}
