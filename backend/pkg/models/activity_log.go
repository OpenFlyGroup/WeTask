package models

import (
	"time"
)

// ? stored in MongoDB
type ActivityLog struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	TaskID    *uint     `json:"taskId,omitempty" bson:"taskId,omitempty"`
	BoardID   *uint     `json:"boardId,omitempty" bson:"boardId,omitempty"`
	TeamID    *uint     `json:"teamId,omitempty" bson:"teamId,omitempty"`
	UserID    uint      `json:"userId" bson:"userId"`
	Action    string    `json:"action" bson:"action"`
	Details   string    `json:"details" bson:"details"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
