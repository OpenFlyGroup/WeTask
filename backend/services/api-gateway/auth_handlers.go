package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wetask/backend/pkg/common"
)

// RegisterRequest represents user registration request
// @Description User registration request
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"` // User email address
	Password string `json:"password" example:"password123" binding:"required,min=6"`   // User password (min 6 characters)
	Name     string `json:"name" example:"John Doe" binding:"required"`                // User full name
}

// LoginRequest represents user login request
// @Description User login request
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"` // User email address
	Password string `json:"password" example:"password123" binding:"required"`         // User password
}

// RefreshRequest represents token refresh request
// @Description Token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." binding:"required"` // Refresh token
}

// AuthResponse represents authentication response
// @Description Authentication response with tokens
type AuthResponse struct {
	AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`  // JWT access token
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT refresh token
}

func authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			ctx.Abort()
			return
		}

		token := parts[1]
		response, err := common.CallRPC(common.AuthValidate, map[string]interface{}{
			"token": token,
		})

		if err != nil || !response.Success {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		userData, ok := response.Data.(map[string]interface{})
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token data"})
			ctx.Abort()
			return
		}

		userID, ok := userData["id"].(float64)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			ctx.Abort()
			return
		}

		ctx.Set("userId", uint(userID))
		ctx.Next()
	}
}

// handleRegister godoc
// @Summary      Register a new user
// @Description  Register a new user with email, password, and name. Returns JWT tokens for authentication.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Registration request"
// @Success      200      {object}  AuthResponse     "User registered successfully"
// @Failure      400      {object}  ErrorResponse    "Invalid request - validation error"
// @Failure      409      {object}  ErrorResponse    "User already exists"
// @Failure      500      {object}  ErrorResponse    "Internal server error"
// @Router       /auth/register [post]
func handleRegister(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.AuthRegister, map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
		"name":     req.Name,
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

// handleLogin godoc
// @Summary      Login user
// @Description  Authenticate user with email and password. Returns JWT access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest   true  "Login request"
// @Success      200      {object}  AuthResponse   "Login successful - returns access and refresh tokens"
// @Failure      400      {object}  ErrorResponse  "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse  "Invalid credentials"
// @Failure      500      {object}  ErrorResponse  "Internal server error"
// @Router       /auth/login [post]
func handleLogin(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.AuthLogin, map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
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

// handleRefresh godoc
// @Summary      Refresh access token
// @Description  Get a new access token using a valid refresh token. Returns new access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest  true  "Refresh token request"
// @Success      200      {object}  AuthResponse    "Token refreshed successfully - returns new access and refresh tokens"
// @Failure      400      {object}  ErrorResponse   "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse   "Invalid or expired refresh token"
// @Failure      500      {object}  ErrorResponse   "Internal server error"
// @Router       /auth/refresh [post]
func handleRefresh(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.AuthRefresh, map[string]interface{}{
		"refreshToken": req.RefreshToken,
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
