package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wetask/backend/pkg/common"
)

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

// handleCreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task in a column. Title and columnId are required; other fields are optional.
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
