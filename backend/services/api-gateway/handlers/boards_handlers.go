package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wetask/backend/pkg/common"
)

// handleGetBoards godoc
// @Summary      Get all boards
// @Description  Get a list of all boards for teams where the user is a member
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   BoardResponse  "List of boards"
// @Failure      401  {object}  ErrorResponse  "Unauthorized - invalid or missing token"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Router       /boards [get]
func HandleGetBoards(ctx *gin.Context) {
	response, err := common.CallRPC(common.BoardsGetAll, nil)
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
// @Description  Create a new board in the specified team. User must be a team member.
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateBoardRequest  true  "Board creation request"
// @Success      200      {object}  BoardResponse       "Board created successfully"
// @Failure      400      {object}  ErrorResponse       "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse       "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse       "Team not found or user not a team member"
// @Failure      500      {object}  ErrorResponse       "Internal server error"
// @Router       /boards [post]
func HandleCreateBoard(ctx *gin.Context) {
	userIDVal, _ := ctx.Get("userId")
	userID := userIDVal.(uint)

	var req CreateBoardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.BoardsCreate, map[string]interface{}{
		"name":   req.Name,
		"teamId": req.TeamID,
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

// handleGetBoard godoc
// @Summary      Get board by ID
// @Description  Get board information including columns and metadata
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int           true  "Board ID"  example(1)
// @Success      200  {object}  BoardResponse "Board information"
// @Failure      400  {object}  ErrorResponse "Invalid board ID format"
// @Failure      401  {object}  ErrorResponse "Unauthorized - invalid or missing token"
// @Failure      404  {object}  ErrorResponse "Board not found"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /boards/{id} [get]
func HandleGetBoard(ctx *gin.Context) {
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
// @Description  Update board name and metadata
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                 true  "Board ID"  example(1)
// @Param        request  body      UpdateBoardRequest  true  "Board update request"
// @Success      200      {object}  BoardResponse       "Board updated successfully"
// @Failure      400      {object}  ErrorResponse       "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse       "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse       "Board not found"
// @Failure      500      {object}  ErrorResponse       "Internal server error"
// @Router       /boards/{id} [put]
func HandleUpdateBoard(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	var req UpdateBoardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.BoardsUpdate, map[string]interface{}{
		"id":   uint(id),
		"name": req.Name,
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
// @Description  Delete a board and all its columns and tasks
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int            true  "Board ID"  example(1)
// @Success      200  {object}  SuccessResponse "Board deleted successfully"
// @Failure      400  {object}  ErrorResponse   "Invalid board ID format"
// @Failure      401  {object}  ErrorResponse   "Unauthorized - invalid or missing token"
// @Failure      404  {object}  ErrorResponse   "Board not found"
// @Failure      500  {object}  ErrorResponse   "Internal server error"
// @Router       /boards/{id} [delete]
func HandleDeleteBoard(ctx *gin.Context) {
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
// @Description  Create a new column in the specified board
// @Tags         boards
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
func HandleCreateColumn(ctx *gin.Context) {
	var req CreateColumnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.ColumnsCreate, map[string]interface{}{
		"name":     req.Name,
		"boardId":  req.BoardID,
		"position": req.Position,
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
// @Summary      Get columns for board
// @Description  Get all columns in the specified board
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        boardId  query     int              true  "Board ID"  example(1)
// @Success      200      {array}   ColumnResponse   "List of columns"
// @Failure      400      {object}  ErrorResponse    "Invalid board ID format"
// @Failure      401      {object}  ErrorResponse    "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse    "Board not found"
// @Failure      500      {object}  ErrorResponse    "Internal server error"
// @Router       /columns [get]
func HandleGetColumns(ctx *gin.Context) {
	boardIDStr := ctx.Query("boardId")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
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
// @Description  Update column name and position
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                 true  "Column ID"  example(1)
// @Param        request  body      UpdateColumnRequest  true  "Column update request"
// @Success      200      {object}  ColumnResponse       "Column updated successfully"
// @Failure      400      {object}  ErrorResponse        "Invalid request - validation error"
// @Failure      401      {object}  ErrorResponse        "Unauthorized - invalid or missing token"
// @Failure      404      {object}  ErrorResponse        "Column not found"
// @Failure      500      {object}  ErrorResponse        "Internal server error"
// @Router       /columns/{id} [put]
func HandleUpdateColumn(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column ID"})
		return
	}

	var req UpdateColumnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := common.CallRPC(common.ColumnsUpdate, map[string]interface{}{
		"id":       uint(id),
		"name":     req.Name,
		"position": req.Position,
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

// handleDeleteColumn godoc
// @Summary      Delete column
// @Description  Delete a column and all its tasks
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int            true  "Column ID"  example(1)
// @Success      200  {object}  SuccessResponse "Column deleted successfully"
// @Failure      400  {object}  ErrorResponse   "Invalid column ID format"
// @Failure      401  {object}  ErrorResponse   "Unauthorized - invalid or missing token"
// @Failure      404  {object}  ErrorResponse   "Column not found"
// @Failure      500  {object}  ErrorResponse   "Internal server error"
// @Router       /columns/{id} [delete]
func HandleDeleteColumn(ctx *gin.Context) {
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
