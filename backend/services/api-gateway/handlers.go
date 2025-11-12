package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wetask/backend/pkg/common"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		token := parts[1]
		response, err := common.CallRPC(common.AuthValidate, map[string]interface{}{
			"token": token,
		})

		if err != nil || !response.Success {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userData, ok := response.Data.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token data"})
			c.Abort()
			return
		}

		userID, ok := userData["id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		c.Set("userId", uint(userID))
		c.Next()
	}
}

func handleRegister(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.AuthRegister, map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
		"name":     req.Name,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.AuthLogin, map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleRefresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.AuthRefresh, map[string]interface{}{
		"refreshToken": req.RefreshToken,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleGetMe(c *gin.Context) {
	userIDVal, _ := c.Get("userId")
	userID := userIDVal.(uint)

	response, err := common.CallRPC(common.UsersGetMe, map[string]interface{}{
		"userId": userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleGetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	response, err := common.CallRPC(common.UsersGetByID, map[string]interface{}{
		"id": uint(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleUpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleGetTeams(c *gin.Context) {
	response, err := common.CallRPC(common.TeamsGetAll, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleCreateTeam(c *gin.Context) {
	userIDVal, _ := c.Get("userId")
	userID := userIDVal.(uint)

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.TeamsCreate, map[string]interface{}{
		"name":   req.Name,
		"userId": userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleGetTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	response, err := common.CallRPC(common.TeamsGetByID, map[string]interface{}{
		"id": uint(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleAddTeamMember(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req struct {
		UserID uint   `json:"userId" binding:"required"`
		Role   string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.TeamsAddMember, map[string]interface{}{
		"teamId": uint(id),
		"userId": req.UserID,
		"role":   req.Role,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleRemoveTeamMember(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	response, err := common.CallRPC(common.TeamsRemoveMember, map[string]interface{}{
		"teamId": uint(teamID),
		"userId": uint(userID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func handleGetBoards(c *gin.Context) {
	userIDVal, _ := c.Get("userId")
	userID := userIDVal.(uint)

	response, err := common.CallRPC(common.BoardsGetAll, map[string]interface{}{
		"userId": userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleCreateBoard(c *gin.Context) {
	var req struct {
		Title  string `json:"title" binding:"required"`
		TeamID uint   `json:"teamId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.BoardsCreate, map[string]interface{}{
		"title":  req.Title,
		"teamId": req.TeamID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleGetBoard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	response, err := common.CallRPC(common.BoardsGetByID, map[string]interface{}{
		"id": uint(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleUpdateBoard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	var req struct {
		Title string `json:"title"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.BoardsUpdate, map[string]interface{}{
		"id":    uint(id),
		"title": req.Title,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleDeleteBoard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	response, err := common.CallRPC(common.BoardsDelete, map[string]interface{}{
		"id": uint(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func handleCreateColumn(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		BoardID uint   `json:"boardId" binding:"required"`
		Order   int    `json:"order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.ColumnsCreate, map[string]interface{}{
		"title":   req.Title,
		"boardId": req.BoardID,
		"order":   req.Order,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleGetColumns(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("boardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	response, err := common.CallRPC(common.ColumnsGetByBoard, map[string]interface{}{
		"boardId": uint(boardID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleUpdateColumn(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column ID"})
		return
	}

	var req struct {
		Title string `json:"title"`
		Order int    `json:"order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleDeleteColumn(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column ID"})
		return
	}

	response, err := common.CallRPC(common.ColumnsDelete, map[string]interface{}{
		"id": uint(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func handleCreateTask(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ColumnID    uint   `json:"columnId" binding:"required"`
		AssignedTo  *uint  `json:"assignedTo"`
		Priority    string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleGetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	response, err := common.CallRPC(common.TasksGetByID, map[string]interface{}{
		"id": uint(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleGetTasksByBoard(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("boardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	response, err := common.CallRPC(common.TasksGetByBoard, map[string]interface{}{
		"boardId": uint(boardID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleUpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		AssignedTo  *uint  `json:"assignedTo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleDeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	response, err := common.CallRPC(common.TasksDelete, map[string]interface{}{
		"id": uint(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func handleMoveTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		ColumnID uint `json:"columnId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.TasksMove, map[string]interface{}{
		"id":       uint(id),
		"columnId": req.ColumnID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleAddComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userIDVal, _ := c.Get("userId")
	userID := userIDVal.(uint)

	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.TasksAddComment, map[string]interface{}{
		"taskId":  uint(id),
		"userId":  userID,
		"message": req.Message,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}

func handleGetComments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	response, err := common.CallRPC(common.TasksGetComments, map[string]interface{}{
		"taskId": uint(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(response.StatusCode, gin.H{"error": response.Error})
		return
	}

	c.JSON(http.StatusOK, response.Data)
}
