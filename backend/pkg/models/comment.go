package models

import (
	"time"
)

// ? stored in MongoDB
type Comment struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	TaskID    uint      `json:"taskId" bson:"taskId"`
	UserID    uint      `json:"userId" bson:"userId"`
	Message   string    `json:"message" bson:"message"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
