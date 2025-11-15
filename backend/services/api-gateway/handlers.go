package main

// ErrorResponse represents error response
// @Description Error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"` // Error message
}

// SuccessResponse represents success response
// @Description Success response
type SuccessResponse struct {
	Success bool `json:"success" example:"true"` // Success status
}
