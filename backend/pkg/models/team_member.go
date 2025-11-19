package models

import (
	"gorm.io/gorm"
)

type TeamMember struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	TeamID    uint           `json:"teamId" gorm:"not null;index;uniqueIndex:idx_team_user"`
	UserID    uint           `json:"userId" gorm:"not null;index;uniqueIndex:idx_team_user"`
	Role      string         `json:"role" gorm:"default:'member'"` //? owner, admin, member
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Team      Team           `json:"-" gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE"`
	User      User           `json:"user,omitempty" gorm:"-"`
}
