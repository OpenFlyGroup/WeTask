package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wetask/backend/pkg/common"
)

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
func HandleCreateTask(ctx *gin.Context) {
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

	data := map[string]any{
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
func HandleGetTask(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	response, err := common.CallRPC(common.TasksGetByID, map[string]any{
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
func HandleGetTasksByBoard(ctx *gin.Context) {
	boardID, err := strconv.ParseUint(ctx.Param("boardId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	response, err := common.CallRPC(common.TasksGetByBoard, map[string]any{
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
func HandleUpdateTask(ctx *gin.Context) {
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

	data := map[string]any{"id": uint(id)}
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
func HandleDeleteTask(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	response, err := common.CallRPC(common.TasksDelete, map[string]any{
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
func HandleMoveTask(ctx *gin.Context) {
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

	response, err := common.CallRPC(common.TasksMove, map[string]any{
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
func HandleAddComment(ctx *gin.Context) {
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

	response, err := common.CallRPC(common.TasksAddComment, map[string]any{
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
func HandleGetComments(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	response, err := common.CallRPC(common.TasksGetComments, map[string]any{
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
