package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Members   []TeamMember   `json:"members,omitempty" gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE"`
	Boards    []Board        `json:"boards,omitempty" gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE"`
}
