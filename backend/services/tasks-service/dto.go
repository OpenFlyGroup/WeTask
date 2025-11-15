package main

import (
	"encoding/json"

	"github.com/wetask/backend/pkg/common"
	"github.com/wetask/backend/pkg/models"
)

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	ColumnID    uint    `json:"columnId"`
	AssignedTo  *uint   `json:"assignedTo,omitempty"`
	Priority    *string `json:"priority,omitempty"`
}

type GetTaskByIDRequest struct {
	ID uint `json:"id"`
}

type GetTasksByBoardRequest struct {
	BoardID uint `json:"boardId"`
}

type UpdateTaskRequest struct {
	ID          uint    `json:"id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Priority    *string `json:"priority,omitempty"`
	AssignedTo  *uint   `json:"assignedTo,omitempty"`
}

type DeleteTaskRequest struct {
	ID uint `json:"id"`
}

type MoveTaskRequest struct {
	ID       uint `json:"id"`
	ColumnID uint `json:"columnId"`
}

type AddCommentRequest struct {
	TaskID  uint   `json:"taskId"`
	UserID  uint   `json:"userId"`
	Message string `json:"message"`
}

type GetCommentsRequest struct {
	TaskID uint `json:"taskId"`
}

type taskResponse struct {
	Success    bool         `json:"success"`
	Data       *models.Task `json:"data,omitempty"`
	Error      string       `json:"error,omitempty"`
	StatusCode int          `json:"statusCode,omitempty"`
}

type tasksListResponse struct {
	Success    bool          `json:"success"`
	Data       []models.Task `json:"data,omitempty"`
	Error      string        `json:"error,omitempty"`
	StatusCode int           `json:"statusCode,omitempty"`
}

type commentResponse struct {
	Success    bool            `json:"success"`
	Data       *models.Comment `json:"data,omitempty"`
	Error      string          `json:"error,omitempty"`
	StatusCode int             `json:"statusCode,omitempty"`
}

type commentsListResponse struct {
	Success    bool             `json:"success"`
	Data       []models.Comment `json:"data,omitempty"`
	Error      string           `json:"error,omitempty"`
	StatusCode int              `json:"statusCode,omitempty"`
}

type successResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}

func toRPC(v any) common.RPCResponse {
	b, _ := json.Marshal(v)
	var resp common.RPCResponse
	_ = json.Unmarshal(b, &resp)
	return resp
}
