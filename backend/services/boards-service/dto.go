package main

import (
	"encoding/json"

	"github.com/wetask/backend/pkg/common"
	"github.com/wetask/backend/pkg/models"
)

// ? Request DTOs

type CreateBoardRequest struct {
	Title  string `json:"title"`
	TeamID uint   `json:"teamId"`
}

type GetAllBoardsRequest struct {
	UserID uint `json:"userId"`
}

type GetBoardByIDRequest struct {
	ID uint `json:"id"`
}

type UpdateBoardRequest struct {
	ID    uint    `json:"id"`
	Title *string `json:"title,omitempty"`
}

type DeleteBoardRequest struct {
	ID uint `json:"id"`
}

type GetBoardsByTeamRequest struct {
	TeamID uint `json:"teamId"`
}

type CreateColumnRequest struct {
	Title   string `json:"title"`
	BoardID uint   `json:"boardId"`
	Order   int    `json:"order"`
}

type GetColumnsByBoardRequest struct {
	BoardID uint `json:"boardId"`
}

type UpdateColumnRequest struct {
	ID    uint    `json:"id"`
	Title *string `json:"title,omitempty"`
	Order *int    `json:"order,omitempty"`
}

type DeleteColumnRequest struct {
	ID uint `json:"id"`
}

// ? Response DTOs

type boardResponse struct {
	Success    bool          `json:"success"`
	Data       *models.Board `json:"data,omitempty"`
	Error      string        `json:"error,omitempty"`
	StatusCode int           `json:"statusCode,omitempty"`
}

type boardsListResponse struct {
	Success    bool           `json:"success"`
	Data       []models.Board `json:"data,omitempty"`
	Error      string         `json:"error,omitempty"`
	StatusCode int            `json:"statusCode,omitempty"`
}

type columnResponse struct {
	Success    bool           `json:"success"`
	Data       *models.Column `json:"data,omitempty"`
	Error      string         `json:"error,omitempty"`
	StatusCode int            `json:"statusCode,omitempty"`
}

type columnsListResponse struct {
	Success    bool            `json:"success"`
	Data       []models.Column `json:"data,omitempty"`
	Error      string          `json:"error,omitempty"`
	StatusCode int             `json:"statusCode,omitempty"`
}

type successResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}

func toRPC(resp any) common.RPCResponse {
	b, _ := json.Marshal(resp)
	var generic common.RPCResponse
	_ = json.Unmarshal(b, &generic)
	return generic
}
