package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wetask/backend/pkg/common"
)

// UserResponse represents user information
// @Description User information response
type UserResponse struct {
	ID        uint   `json:"id" example:"1"`                           // User ID
	Email     string `json:"email" example:"user@example.com"`         // User email
	Name      string `json:"name" example:"John Doe"`                  // User name
	CreatedAt string `json:"createdAt" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt string `json:"updatedAt" example:"2024-01-01T00:00:00Z"` // Last update timestamp
}

// UpdateUserRequest represents user update request
// @Description User update request
type UpdateUserRequest struct {
	Name  string `json:"name" example:"John Doe Updated"`      // Updated user name (optional)
	Email string `json:"email" example:"newemail@example.com"` // Updated user email (optional)
}

// handleGetMe godoc
// @Summary      Get current user
// @Description  Get the authenticated user's information. User ID is extracted from JWT token.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  UserResponse   "User information"
// @Failure      401  {object}  ErrorResponse  "Unauthorized - invalid or missing token"
// @Failure      404  {object}  ErrorResponse  "User not found"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Router       /users/me [get]
func handleGetMe(ctx *gin.Context) {
	userIDVal, _ := ctx.Get("userId")
	userID := userIDVal.(uint)

	response, err := common.CallRPC(common.UsersGetMe, map[string]interface{}{
		"userId": userID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		ctx.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	ctx.JSON(http.StatusOK, response.Data)
}

// handleGetUser godoc
// @Summary      Get user by ID
// @Description  Get user information by user ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int           true  "User ID"  example(1)
// @Success      200  {object}  UserResponse  "User information"
// @Failure      400  {object}  ErrorResponse "Invalid user ID format"
// @Failure      401  {object}  ErrorResponse "Unauthorized - invalid or missing token"
// @Failure      404  {object}  ErrorResponse "User not found"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /users/{id} [get]
func handleGetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	response, err := common.CallRPC(common.UsersGetByID, map[string]interface{}{
		"id": uint(id),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		ctx.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	ctx.JSON(http.StatusOK, response.Data)
}

// handleUpdateUser godoc
// @Summary      Update user
// @Description  Update user information (name and/or email). Both fields are optional.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                true  "User ID"  example(1)
// @Param        request  body      UpdateUserRequest  true  "Update request - all fields optional"
// @Success      200      {object}  UserResponse       "User updated successfully"
// @Failure      400      {object}  ErrorResponse      "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse      "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse      "User not found"
// @Failure      500      {object}  ErrorResponse      "Internal server error"
// @Router       /users/{id} [patch]
func handleUpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{"id": uint(id)}
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.Email != "" {
		data["email"] = req.Email
	}

	response, err := common.CallRPC(common.UsersUpdate, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		ctx.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	ctx.JSON(http.StatusOK, response.Data)
}
