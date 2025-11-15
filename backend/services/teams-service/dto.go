package main

import (
	"encoding/json"
	"time"

	"github.com/wetask/backend/pkg/common"
	"github.com/wetask/backend/pkg/models"
)

type CreateTeamRequest struct {
	Name   string `json:"name"`
	UserID uint   `json:"userId"`
}

type GetTeamByIDRequest struct {
	ID uint `json:"id"`
}

type AddMemberRequest struct {
	TeamID uint   `json:"teamId"`
	UserID uint   `json:"userId"`
	Role   string `json:"role,omitempty"`
}

type RemoveMemberRequest struct {
	TeamID uint `json:"teamId"`
	UserID uint `json:"userId"`
}

type GetUserTeamsRequest struct {
	UserID uint `json:"userId"`
}

type teamResponse struct {
	Success    bool         `json:"success"`
	Data       *models.Team `json:"data,omitempty"`
	Error      string       `json:"error,omitempty"`
	StatusCode int          `json:"statusCode,omitempty"`
}

type teamsListResponse struct {
	Success    bool          `json:"success"`
	Data       []models.Team `json:"data,omitempty"`
	Error      string        `json:"error,omitempty"`
	StatusCode int           `json:"statusCode,omitempty"`
}

type memberResponse struct {
	Success    bool               `json:"success"`
	Data       *models.TeamMember `json:"data,omitempty"`
	Error      string             `json:"error,omitempty"`
	StatusCode int                `json:"statusCode,omitempty"`
}

type userTeamResponse struct {
	Success    bool              `json:"success"`
	Data       []UserTeamSummary `json:"data,omitempty"`
	Error      string            `json:"error,omitempty"`
	StatusCode int               `json:"statusCode,omitempty"`
}

type successResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}

type UserTeamSummary struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

func toRPC(v any) common.RPCResponse {
	b, _ := json.Marshal(v)
	var resp common.RPCResponse
	_ = json.Unmarshal(b, &resp)
	return resp
}
