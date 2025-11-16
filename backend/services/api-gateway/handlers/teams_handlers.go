package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wetask/backend/pkg/common"
)

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
func HandleGetTeams(ctx *gin.Context) {
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
func HandleCreateTeam(ctx *gin.Context) {
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
func HandleGetTeam(ctx *gin.Context) {
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
func HandleAddTeamMember(ctx *gin.Context) {
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
func HandleRemoveTeamMember(ctx *gin.Context) {
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
