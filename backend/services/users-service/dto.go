package main

import (
	"encoding/json"
	"time"

	"github.com/wetask/backend/pkg/common"
	"github.com/wetask/backend/pkg/models"
)

type CreateUserRequest struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GetUserByIDRequest struct {
	ID uint `json:"id"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	ID    uint    `json:"id"`
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

type GetMeRequest struct {
	UserID uint `json:"userId"`
}

type userResponse struct {
	Success    bool         `json:"success"`
	Data       *models.User `json:"data,omitempty"`
	Error      string       `json:"error,omitempty"`
	StatusCode int          `json:"statusCode,omitempty"`
}

type UserData struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toRPC(v any) common.RPCResponse {
	b, _ := json.Marshal(v)
	var resp common.RPCResponse
	_ = json.Unmarshal(b, &resp)
	return resp
}
