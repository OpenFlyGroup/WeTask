package main

import (
	"encoding/json"
	"time"

	"github.com/wetask/backend/pkg/common"
)

// ? Request DTOs
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type ValidateRequest struct {
	Token string `json:"token"`
}

// ? Response DTOs
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
}

type ValidateResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// ? RPC wrapper types to unify with other services
type authResponseWrapper struct {
	Success    bool         `json:"success"`
	Data       AuthResponse `json:"data,omitempty"`
	Error      string       `json:"error,omitempty"`
	StatusCode int          `json:"statusCode,omitempty"`
}

type refreshResponseWrapper struct {
	Success    bool            `json:"success"`
	Data       RefreshResponse `json:"data,omitempty"`
	Error      string          `json:"error,omitempty"`
	StatusCode int             `json:"statusCode,omitempty"`
}

type validateResponseWrapper struct {
	Success    bool             `json:"success"`
	Data       ValidateResponse `json:"data,omitempty"`
	Error      string           `json:"error,omitempty"`
	StatusCode int              `json:"statusCode,omitempty"`
}

func toRPC(v any) common.RPCResponse {
	b, _ := json.Marshal(v)
	var resp common.RPCResponse
	_ = json.Unmarshal(b, &resp)
	return resp
}
