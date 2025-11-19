package handlers

// ? General responses

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

// ? Auth requests and responses

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

// ? Board requests and responses

// BoardResponse represents board information
// @Description Board information response
type BoardResponse struct {
	ID        uint   `json:"id" example:"1"`                           // Board ID
	Name      string `json:"name" example:"Q1 Project"`                // Board name
	TeamID    uint   `json:"teamId" example:"1"`                       // Team ID
	CreatedAt string `json:"createdAt" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt string `json:"updatedAt" example:"2024-01-01T00:00:00Z"` // Last update timestamp
}

// ColumnResponse represents column information
// @Description Column information response
type ColumnResponse struct {
	ID        uint   `json:"id" example:"1"`                           // Column ID
	Name      string `json:"name" example:"To Do"`                     // Column name
	BoardID   uint   `json:"boardId" example:"1"`                      // Board ID
	Position  int    `json:"position" example:"0"`                     // Column position in board
	CreatedAt string `json:"createdAt" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt string `json:"updatedAt" example:"2024-01-01T00:00:00Z"` // Last update timestamp
}

// CreateBoardRequest represents board creation request
// @Description Board creation request
type CreateBoardRequest struct {
	Title  string `json:"title" example:"Q1 Project" binding:"required"` // Board title
	TeamID uint   `json:"teamId" example:"1" binding:"required"`         // Team ID
}

// UpdateBoardRequest represents board update request
// @Description Board update request
type UpdateBoardRequest struct {
	Name string `json:"name" example:"Q2 Project" binding:"required"` // Board name
}

// CreateColumnRequest represents column creation request
// @Description Column creation request
type CreateColumnRequest struct {
	Title    string `json:"title" example:"To Do" binding:"required"` // Column title
	BoardID  uint   `json:"boardId" example:"1" binding:"required"`   // Board ID
	Position int    `json:"position" example:"0"`                     // Column position
}

// UpdateColumnRequest represents column update request
// @Description Column update request
type UpdateColumnRequest struct {
	Title    string `json:"title" example:"In Progress"` // Column title
	Position int    `json:"position" example:"1"`        // Column position
}

// GetColumnsRequest represents a request to fetch columns for a board
// @Description Get columns by board request (supports URI `/board/:boardId` or query `?boardId=...`)
type GetColumnsRequest struct {
	BoardID uint `json:"boardId" form:"boardId" uri:"boardId" example:"1" binding:"required"` // Board ID
}

// ? Task requests and responses

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

// ? Team requests and responses

// TeamResponse represents team information
// @Description Team information response
type TeamResponse struct {
	ID        uint   `json:"id" example:"1"`                           // Team ID
	Name      string `json:"name" example:"Development Team"`          // Team name
	CreatedAt string `json:"createdAt" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt string `json:"updatedAt" example:"2024-01-01T00:00:00Z"` // Last update timestamp
}

// MemberResponse represents a team member
// @Description Team member information
type MemberResponse struct {
	ID        uint         `json:"id" example:"1"`                           // Member ID
	TeamID    uint         `json:"teamId" example:"1"`                       // Team ID
	UserID    uint         `json:"userId" example:"2"`                       // User ID
	Role      string       `json:"role" example:"owner"`                     // Role: owner, admin, member
	User      UserResponse `json:"user,omitempty"`                           // Nested user info
	CreatedAt string       `json:"createdAt" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt string       `json:"updatedAt" example:"2024-01-01T00:00:00Z"` // Last update timestamp
}

// Extended TeamResponse with members and boards
// @Description Team information with members and boards
type TeamFullResponse struct {
	ID        uint             `json:"id" example:"1"`                           // Team ID
	Name      string           `json:"name" example:"Development Team"`          // Team name
	Members   []MemberResponse `json:"members"`                                  // Team members
	Boards    []BoardResponse  `json:"boards"`                                   // Team boards
	CreatedAt string           `json:"createdAt" example:"2024-01-01T00:00:00Z"` // Creation timestamp
	UpdatedAt string           `json:"updatedAt" example:"2024-01-01T00:00:00Z"` // Last update timestamp
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

// ? Team requests and responses

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
