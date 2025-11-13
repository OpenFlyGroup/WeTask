package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wetask/backend/pkg/common"
)

// Request/Response Models for Swagger documentation

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

// TeamResponse represents team information
// @Description Team information response
type TeamResponse struct {
	ID        uint   `json:"id" example:"1"`                           // Team ID
	Name      string `json:"name" example:"Development Team"`          // Team name
	CreatedAt string `json:"createdAt" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt string `json:"updatedAt" example:"2024-01-01T00:00:00Z"` // Last update timestamp
}

// CreateTeamRequest represents team creation request
// @Description Team creation request
type CreateTeamRequest struct {
	Name string `json:"name" example:"Development Team" binding:"required"` // Team name
}

// AddTeamMemberRequest represents add team member request
// @Description Add team member request
type AddTeamMemberRequest struct {
	UserID uint   `json:"userId" example:"2" binding:"required"` // User ID to add
	Role   string `json:"role" example:"member"`                 // Member role (optional)
}

// BoardResponse represents board information
// @Description Board information response
type BoardResponse struct {
	ID        uint   `json:"id" example:"1"`                           // Board ID
	Title     string `json:"title" example:"Project Board"`            // Board title
	TeamID    uint   `json:"teamId" example:"1"`                       // Associated team ID
	CreatedAt string `json:"createdAt" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt string `json:"updatedAt" example:"2024-01-01T00:00:00Z"` // Last update timestamp
}

// CreateBoardRequest represents board creation request
// @Description Board creation request
type CreateBoardRequest struct {
	Title  string `json:"title" example:"Project Board" binding:"required"` // Board title
	TeamID uint   `json:"teamId" example:"1" binding:"required"`            // Team ID
}

// UpdateBoardRequest represents board update request
// @Description Board update request
type UpdateBoardRequest struct {
	Title string `json:"title" example:"Updated Board Title"` // Updated board title
}

// ColumnResponse represents column information
// @Description Column information response
type ColumnResponse struct {
	ID        uint   `json:"id" example:"1"`                           // Column ID
	Title     string `json:"title" example:"To Do"`                    // Column title
	Order     int    `json:"order" example:"1"`                        // Column order
	BoardID   uint   `json:"boardId" example:"1"`                      // Associated board ID
	CreatedAt string `json:"createdAt" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt string `json:"updatedAt" example:"2024-01-01T00:00:00Z"` // Last update timestamp
}

// CreateColumnRequest represents column creation request
// @Description Column creation request
type CreateColumnRequest struct {
	Title   string `json:"title" example:"To Do" binding:"required"` // Column title
	BoardID uint   `json:"boardId" example:"1" binding:"required"`   // Board ID
	Order   int    `json:"order" example:"1"`                        // Column order (optional)
}

// UpdateColumnRequest represents column update request
// @Description Column update request
type UpdateColumnRequest struct {
	Title string `json:"title" example:"In Progress"` // Updated column title (optional)
	Order int    `json:"order" example:"2"`           // Updated column order (optional)
}

// TaskResponse represents task information
// @Description Task information response
type TaskResponse struct {
	ID          uint    `json:"id" example:"1"`                                   // Task ID
	Title       string  `json:"title" example:"Implement feature"`                // Task title
	Description *string `json:"description" example:"Add new feature to the app"` // Task description (optional)
	Priority    *string `json:"priority" example:"high"`                          // Task priority: low, medium, high (optional)
	ColumnID    uint    `json:"columnId" example:"1"`                             // Associated column ID
	AssignedTo  *uint   `json:"assignedTo" example:"2"`                           // Assigned user ID (optional)
	CreatedAt   string  `json:"createdAt" example:"2024-01-01T00:00:00Z"`         // Creation timestamp
	UpdatedAt   string  `json:"updatedAt" example:"2024-01-01T00:00:00Z"`         // Last update timestamp
}

// CreateTaskRequest represents task creation request
// @Description Task creation request
type CreateTaskRequest struct {
	Title       string `json:"title" example:"Implement feature" binding:"required"` // Task title
	Description string `json:"description" example:"Add new feature to the app"`     // Task description (optional)
	ColumnID    uint   `json:"columnId" example:"1" binding:"required"`              // Column ID
	AssignedTo  *uint  `json:"assignedTo" example:"2"`                               // Assigned user ID (optional)
	Priority    string `json:"priority" example:"high"`                              // Task priority: low, medium, high (optional)
}

// UpdateTaskRequest represents task update request
// @Description Task update request
type UpdateTaskRequest struct {
	Title       string `json:"title" example:"Updated task title"`        // Updated task title (optional)
	Description string `json:"description" example:"Updated description"` // Updated task description (optional)
	Priority    string `json:"priority" example:"medium"`                 // Updated priority: low, medium, high (optional)
	AssignedTo  *uint  `json:"assignedTo" example:"3"`                    // Updated assigned user ID (optional)
}

// MoveTaskRequest represents move task request
// @Description Move task request
type MoveTaskRequest struct {
	ColumnID uint `json:"columnId" example:"2" binding:"required"` // Target column ID
}

// CommentResponse represents comment information
// @Description Comment information response
type CommentResponse struct {
	ID        string `json:"id" example:"507f1f77bcf86cd799439011"`    // Comment ID
	TaskID    uint   `json:"taskId" example:"1"`                       // Associated task ID
	UserID    uint   `json:"userId" example:"1"`                       // Comment author user ID
	Message   string `json:"message" example:"This looks good!"`       // Comment message
	CreatedAt string `json:"createdAt" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt string `json:"updatedAt" example:"2024-01-01T00:00:00Z"` // Last update timestamp
}

// AddCommentRequest represents add comment request
// @Description Add comment request
type AddCommentRequest struct {
	Message string `json:"message" example:"This looks good!" binding:"required"` // Comment message
}

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

// handleGetTeams godoc
// @Summary      Get all teams
// @Description  Get a list of all teams accessible to the authenticated user
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   TeamResponse  "List of teams"
// @Failure      401  {object}  ErrorResponse "Unauthorized - invalid or missing token"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /teams [get]
func handleGetTeams(ctx *gin.Context) {
	response, err := common.CallRPC(common.TeamsGetAll, nil)
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

// handleCreateTeam godoc
// @Summary      Create a new team
// @Description  Create a new team with the authenticated user as owner. The creator is automatically added as a member.
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateTeamRequest  true  "Team creation request"
// @Success      200      {object}  TeamResponse       "Team created successfully"
// @Failure      400      {object}  ErrorResponse      "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse      "Unauthorized - invalid or missing token"
// @Failure      500      {object}  ErrorResponse      "Internal server error"
// @Router       /teams [post]
func handleCreateTeam(ctx *gin.Context) {
	userIDVal, _ := ctx.Get("userId")
	userID := userIDVal.(uint)

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.TeamsCreate, map[string]interface{}{
		"name":   req.Name,
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

// handleGetTeam godoc
// @Summary      Get team by ID
// @Description  Get team information by team ID including members and boards
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int           true  "Team ID"  example(1)
// @Success      200  {object}  TeamResponse  "Team information"
// @Failure      400  {object}  ErrorResponse "Invalid team ID format"
// @Failure      401  {object}  ErrorResponse "Unauthorized - invalid or missing token"
// @Failure      404  {object}  ErrorResponse "Team not found"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /teams/{id} [get]
func handleGetTeam(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	response, err := common.CallRPC(common.TeamsGetByID, map[string]interface{}{
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

// handleAddTeamMember godoc
// @Summary      Add member to team
// @Description  Add a user to a team with an optional role. The user must exist in the system.
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                   true  "Team ID"  example(1)
// @Param        request  body      AddTeamMemberRequest  true  "Add member request"
// @Success      200      {object}  TeamResponse          "Member added successfully"
// @Failure      400      {object}  ErrorResponse         "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse         "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse         "Team or user not found"
// @Failure      409      {object}  ErrorResponse         "User is already a member of the team"
// @Failure      500      {object}  ErrorResponse         "Internal server error"
// @Router       /teams/{id}/members [post]
func handleAddTeamMember(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req struct {
		UserID uint   `json:"userId" binding:"required"`
		Role   string `json:"role"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.TeamsAddMember, map[string]interface{}{
		"teamId": uint(id),
		"userId": req.UserID,
		"role":   req.Role,
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

// handleRemoveTeamMember godoc
// @Summary      Remove member from team
// @Description  Remove a user from a team. The user will lose access to all team boards.
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int            true  "Team ID"   example(1)
// @Param        userId  path      int            true  "User ID"   example(2)
// @Success      200     {object}  SuccessResponse "Member removed successfully"
// @Failure      400     {object}  ErrorResponse   "Invalid team or user ID format"
// @Failure      401     {object}  ErrorResponse   "Unauthorized - invalid or missing token"
// @Failure      404     {object}  ErrorResponse   "Team or member not found"
// @Failure      500     {object}  ErrorResponse   "Internal server error"
// @Router       /teams/{id}/members/{userId} [delete]
func handleRemoveTeamMember(ctx *gin.Context) {
	teamID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	userID, err := strconv.ParseUint(ctx.Param("userId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	response, err := common.CallRPC(common.TeamsRemoveMember, map[string]interface{}{
		"teamId": uint(teamID),
		"userId": uint(userID),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		ctx.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// handleGetBoards godoc
// @Summary      Get all boards
// @Description  Get all boards accessible to the authenticated user (boards from teams where user is a member)
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   BoardResponse  "List of boards"
// @Failure      401  {object}  ErrorResponse  "Unauthorized - invalid or missing token"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Router       /boards [get]
func handleGetBoards(ctx *gin.Context) {
	userIDVal, _ := ctx.Get("userId")
	userID := userIDVal.(uint)

	response, err := common.CallRPC(common.BoardsGetAll, map[string]interface{}{
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

// handleCreateBoard godoc
// @Summary      Create a new board
// @Description  Create a new board for a team. User must be a member of the team.
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateBoardRequest  true  "Board creation request"
// @Success      200      {object}  BoardResponse       "Board created successfully"
// @Failure      400      {object}  ErrorResponse       "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse       "Unauthorized - invalid or missing token"
// @Failure      403      {object}  ErrorResponse       "Forbidden - user is not a member of the team"
// @Failure      404      {object}  ErrorResponse       "Team not found"
// @Failure      500      {object}  ErrorResponse       "Internal server error"
// @Router       /boards [post]
func handleCreateBoard(ctx *gin.Context) {
	var req struct {
		Title  string `json:"title" binding:"required"`
		TeamID uint   `json:"teamId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.BoardsCreate, map[string]interface{}{
		"title":  req.Title,
		"teamId": req.TeamID,
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

// handleGetBoard godoc
// @Summary      Get board by ID
// @Description  Get board information by board ID including columns and tasks
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int            true  "Board ID"  example(1)
// @Success      200  {object}  BoardResponse  "Board information"
// @Failure      400  {object}  ErrorResponse  "Invalid board ID format"
// @Failure      401  {object}  ErrorResponse  "Unauthorized - invalid or missing token"
// @Failure      403  {object}  ErrorResponse  "Forbidden - user doesn't have access to this board"
// @Failure      404  {object}  ErrorResponse  "Board not found"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Router       /boards/{id} [get]
func handleGetBoard(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	response, err := common.CallRPC(common.BoardsGetByID, map[string]interface{}{
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

// handleUpdateBoard godoc
// @Summary      Update board
// @Description  Update board information (title). User must have access to the board.
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                 true  "Board ID"  example(1)
// @Param        request  body      UpdateBoardRequest  true  "Update request"
// @Success      200      {object}  BoardResponse       "Board updated successfully"
// @Failure      400      {object}  ErrorResponse       "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse       "Unauthorized - invalid or missing token"
// @Failure      403      {object}  ErrorResponse       "Forbidden - user doesn't have access to this board"
// @Failure      404      {object}  ErrorResponse       "Board not found"
// @Failure      500      {object}  ErrorResponse       "Internal server error"
// @Router       /boards/{id} [put]
func handleUpdateBoard(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	var req struct {
		Title string `json:"title"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.BoardsUpdate, map[string]interface{}{
		"id":    uint(id),
		"title": req.Title,
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

// handleDeleteBoard godoc
// @Summary      Delete board
// @Description  Delete a board by ID. This will also delete all columns and tasks in the board.
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int             true  "Board ID"  example(1)
// @Success      200  {object}  SuccessResponse "Board deleted successfully"
// @Failure      400  {object}  ErrorResponse   "Invalid board ID format"
// @Failure      401  {object}  ErrorResponse   "Unauthorized - invalid or missing token"
// @Failure      403  {object}  ErrorResponse   "Forbidden - user doesn't have permission to delete this board"
// @Failure      404  {object}  ErrorResponse   "Board not found"
// @Failure      500  {object}  ErrorResponse   "Internal server error"
// @Router       /boards/{id} [delete]
func handleDeleteBoard(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	response, err := common.CallRPC(common.BoardsDelete, map[string]interface{}{
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

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// handleCreateColumn godoc
// @Summary      Create a new column
// @Description  Create a new column in a board. Order determines the position of the column.
// @Tags         columns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateColumnRequest  true  "Column creation request"
// @Success      200      {object}  ColumnResponse       "Column created successfully"
// @Failure      400      {object}  ErrorResponse        "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse        "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse        "Board not found"
// @Failure      500      {object}  ErrorResponse        "Internal server error"
// @Router       /columns [post]
func handleCreateColumn(ctx *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		BoardID uint   `json:"boardId" binding:"required"`
		Order   int    `json:"order"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.ColumnsCreate, map[string]interface{}{
		"title":   req.Title,
		"boardId": req.BoardID,
		"order":   req.Order,
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

// handleGetColumns godoc
// @Summary      Get columns by board
// @Description  Get all columns for a specific board, ordered by their order field
// @Tags         columns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        boardId  path      int             true  "Board ID"  example(1)
// @Success      200      {array}   ColumnResponse  "List of columns"
// @Failure      400      {object}  ErrorResponse   "Invalid board ID format"
// @Failure      401      {object}  ErrorResponse   "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse   "Board not found"
// @Failure      500      {object}  ErrorResponse   "Internal server error"
// @Router       /columns/board/{boardId} [get]
func handleGetColumns(ctx *gin.Context) {
	boardID, err := strconv.ParseUint(ctx.Param("boardId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	response, err := common.CallRPC(common.ColumnsGetByBoard, map[string]interface{}{
		"boardId": uint(boardID),
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

// handleUpdateColumn godoc
// @Summary      Update column
// @Description  Update column information (title and/or order). Both fields are optional.
// @Tags         columns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                 true  "Column ID"  example(1)
// @Param        request  body      UpdateColumnRequest true  "Update request - all fields optional"
// @Success      200      {object}  ColumnResponse      "Column updated successfully"
// @Failure      400      {object}  ErrorResponse       "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse       "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse       "Column not found"
// @Failure      500      {object}  ErrorResponse       "Internal server error"
// @Router       /columns/{id} [put]
func handleUpdateColumn(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column ID"})
		return
	}

	var req struct {
		Title string `json:"title"`
		Order int    `json:"order"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{"id": uint(id)}
	if req.Title != "" {
		data["title"] = req.Title
	}
	if req.Order > 0 {
		data["order"] = req.Order
	}

	response, err := common.CallRPC(common.ColumnsUpdate, data)
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

// handleDeleteColumn godoc
// @Summary      Delete column
// @Description  Delete a column by ID. This will also delete all tasks in the column.
// @Tags         columns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int             true  "Column ID"  example(1)
// @Success      200  {object}  SuccessResponse "Column deleted successfully"
// @Failure      400  {object}  ErrorResponse   "Invalid column ID format"
// @Failure      401  {object}  ErrorResponse   "Unauthorized - invalid or missing token"
// @Failure      404  {object}  ErrorResponse   "Column not found"
// @Failure      500  {object}  ErrorResponse   "Internal server error"
// @Router       /columns/{id} [delete]
func handleDeleteColumn(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column ID"})
		return
	}

	response, err := common.CallRPC(common.ColumnsDelete, map[string]interface{}{
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

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// handleCreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task in a column. Priority can be: low, medium, high (default: medium)
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateTaskRequest  true  "Task creation request"
// @Success      200      {object}  TaskResponse       "Task created successfully"
// @Failure      400      {object}  ErrorResponse      "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse      "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse      "Column not found"
// @Failure      500      {object}  ErrorResponse      "Internal server error"
// @Router       /tasks [post]
func handleCreateTask(ctx *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ColumnID    uint   `json:"columnId" binding:"required"`
		AssignedTo  *uint  `json:"assignedTo"`
		Priority    string `json:"priority"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{
		"title":    req.Title,
		"columnId": req.ColumnID,
	}
	if req.Description != "" {
		data["description"] = req.Description
	}
	if req.AssignedTo != nil {
		data["assignedTo"] = *req.AssignedTo
	}
	if req.Priority != "" {
		data["priority"] = req.Priority
	}

	response, err := common.CallRPC(common.TasksCreate, data)
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

// handleGetTask godoc
// @Summary      Get task by ID
// @Description  Get task information by task ID including assigned user and column details
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int            true  "Task ID"  example(1)
// @Success      200  {object}  TaskResponse   "Task information"
// @Failure      400  {object}  ErrorResponse  "Invalid task ID format"
// @Failure      401  {object}  ErrorResponse  "Unauthorized - invalid or missing token"
// @Failure      404  {object}  ErrorResponse  "Task not found"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Router       /tasks/{id} [get]
func handleGetTask(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	response, err := common.CallRPC(common.TasksGetByID, map[string]interface{}{
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

// handleGetTasksByBoard godoc
// @Summary      Get tasks by board
// @Description  Get all tasks for a specific board, organized by columns
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        boardId  path      int             true  "Board ID"  example(1)
// @Success      200      {array}   TaskResponse    "List of tasks"
// @Failure      400      {object}  ErrorResponse   "Invalid board ID format"
// @Failure      401      {object}  ErrorResponse   "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse   "Board not found"
// @Failure      500      {object}  ErrorResponse   "Internal server error"
// @Router       /tasks/board/{boardId} [get]
func handleGetTasksByBoard(ctx *gin.Context) {
	boardID, err := strconv.ParseUint(ctx.Param("boardId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	response, err := common.CallRPC(common.TasksGetByBoard, map[string]interface{}{
		"boardId": uint(boardID),
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

// handleUpdateTask godoc
// @Summary      Update task
// @Description  Update task information (title, description, priority, assignedTo). All fields are optional. Priority: low, medium, high
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                true  "Task ID"  example(1)
// @Param        request  body      UpdateTaskRequest  true  "Update request - all fields optional"
// @Success      200      {object}  TaskResponse       "Task updated successfully"
// @Failure      400      {object}  ErrorResponse      "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse      "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse      "Task not found"
// @Failure      500      {object}  ErrorResponse      "Internal server error"
// @Router       /tasks/{id} [put]
func handleUpdateTask(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		AssignedTo  *uint  `json:"assignedTo"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{"id": uint(id)}
	if req.Title != "" {
		data["title"] = req.Title
	}
	if req.Description != "" {
		data["description"] = req.Description
	}
	if req.Priority != "" {
		data["priority"] = req.Priority
	}
	if req.AssignedTo != nil {
		data["assignedTo"] = *req.AssignedTo
	}

	response, err := common.CallRPC(common.TasksUpdate, data)
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

// handleDeleteTask godoc
// @Summary      Delete task
// @Description  Delete a task by ID. This will also delete all comments associated with the task.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int             true  "Task ID"  example(1)
// @Success      200  {object}  SuccessResponse "Task deleted successfully"
// @Failure      400  {object}  ErrorResponse   "Invalid task ID format"
// @Failure      401  {object}  ErrorResponse   "Unauthorized - invalid or missing token"
// @Failure      404  {object}  ErrorResponse   "Task not found"
// @Failure      500  {object}  ErrorResponse   "Internal server error"
// @Router       /tasks/{id} [delete]
func handleDeleteTask(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	response, err := common.CallRPC(common.TasksDelete, map[string]interface{}{
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

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// handleMoveTask godoc
// @Summary      Move task to another column
// @Description  Move a task from one column to another within the same board
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int              true  "Task ID"  example(1)
// @Param        request  body      MoveTaskRequest  true  "Move request"
// @Success      200      {object}  TaskResponse     "Task moved successfully"
// @Failure      400      {object}  ErrorResponse    "Invalid request - validation error or columns not in same board"
// @Failure      401      {object}  ErrorResponse    "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse    "Task or column not found"
// @Failure      500      {object}  ErrorResponse    "Internal server error"
// @Router       /tasks/{id}/move [put]
func handleMoveTask(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		ColumnID uint `json:"columnId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.TasksMove, map[string]interface{}{
		"id":       uint(id),
		"columnId": req.ColumnID,
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

// handleAddComment godoc
// @Summary      Add comment to task
// @Description  Add a comment to a task. The comment author is automatically set to the authenticated user.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                true  "Task ID"  example(1)
// @Param        request  body      AddCommentRequest  true  "Comment request"
// @Success      200      {object}  CommentResponse    "Comment added successfully"
// @Failure      400      {object}  ErrorResponse      "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse      "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse      "Task not found"
// @Failure      500      {object}  ErrorResponse      "Internal server error"
// @Router       /tasks/{id}/comment [post]
func handleAddComment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userIDVal, _ := ctx.Get("userId")
	userID := userIDVal.(uint)

	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.TasksAddComment, map[string]interface{}{
		"taskId":  uint(id),
		"userId":  userID,
		"message": req.Message,
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

// handleGetComments godoc
// @Summary      Get task comments
// @Description  Get all comments for a specific task, ordered by creation date (oldest first)
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int               true  "Task ID"  example(1)
// @Success      200  {array}   CommentResponse   "List of comments"
// @Failure      400  {object}  ErrorResponse     "Invalid task ID format"
// @Failure      401  {object}  ErrorResponse     "Unauthorized - invalid or missing token"
// @Failure      404  {object}  ErrorResponse     "Task not found"
// @Failure      500  {object}  ErrorResponse     "Internal server error"
// @Router       /tasks/{id}/comments [get]
func handleGetComments(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	response, err := common.CallRPC(common.TasksGetComments, map[string]interface{}{
		"taskId": uint(id),
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
