package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description *string        `json:"description"`
	Status      string         `json:"status" gorm:"not null"`
	Priority    *string        `json:"priority" gorm:"default:'medium'"` //? low, medium, high
	ColumnID    uint           `json:"columnId" gorm:"not null;index"`
	AssignedTo  *uint          `json:"assignedTo" gorm:"index"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Column      Column         `json:"column,omitempty" gorm:"foreignKey:ColumnID;constraint:OnDelete:CASCADE"`
	User        *User          `json:"user,omitempty" gorm:"-"`
}
