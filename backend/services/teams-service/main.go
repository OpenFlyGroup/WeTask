package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/wetask/backend/pkg/common"
	"github.com/wetask/backend/pkg/models"
)

func attachUserProfile(member *models.TeamMember) {
	if member == nil {
		return
	}
	if rpcResp, err := common.CallRPC(common.UsersGetByID, map[string]interface{}{"id": member.UserID}); err == nil && rpcResp != nil && rpcResp.Success {
		if data, ok := rpcResp.Data.(map[string]interface{}); ok {
			var user models.User
			if id, ok := data["id"].(float64); ok {
				user.ID = uint(id)
			}
			if email, ok := data["email"].(string); ok {
				user.Email = email
			}
			if name, ok := data["name"].(string); ok {
				user.Name = name
			}
			member.User = user
		}
	}
}

func attachUsersToMembers(members []models.TeamMember) {
	for i := range members {
		attachUserProfile(&members[i])
	}
}

func main() {
	// ? Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// ? Initialize database and migrate models
	if err := common.InitPostgreSQL(); err != nil {
		log.Fatal("Failed to initialize PostgreSQL:", err)
	}
	if err := common.MigrateTeamsModels(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// ? Initialize RabbitMQ
	if err := common.InitRabbitMQ(); err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer common.CloseRabbitMQ()

	// ? Declare queues
	queues := []string{
		common.TeamsCreate,
		common.TeamsGetAll,
		common.TeamsGetByID,
		common.TeamsAddMember,
		common.TeamsRemoveMember,
		common.TeamsGetUserTeams,
	}

	for _, queue := range queues {
		_, err := common.RabbitMQChannel.QueueDeclare(
			queue, // * name
			true,  // * durable
			false, // * delete when unused
			false, // * exclusive
			false, // * no-wait
			nil,   // * arguments
		)
		if err != nil {
			log.Fatal("Failed to declare queue:", err)
		}
	}

	// ? Start consuming messages
	for _, queue := range queues {
		messages, err := common.RabbitMQChannel.Consume(
			queue, // * queue
			"",    // * consumer
			false, // * auto-ack
			false, // * exclusive
			false, // * no-local
			false, // * no-wait
			nil,   // * args
		)
		if err != nil {
			log.Fatal("Failed to register consumer:", err)
		}

		go handleMessages(queue, messages)
	}

	log.Println("Teams Service is running...")
	select {} // ? Keep running
}

func handleMessages(queue string, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		var response common.RPCResponse

		switch queue {
		case common.TeamsCreate:
			var req CreateTeamRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleCreate(req)
			}

		case common.TeamsGetAll:
			response = handleGetAll()

		case common.TeamsGetByID:
			var req GetTeamByIDRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetByID(req)
			}

		case common.TeamsAddMember:
			var req AddMemberRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleAddMember(req)
			}

		case common.TeamsRemoveMember:
			var req RemoveMemberRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleRemoveMember(req)
			}

		case common.TeamsGetUserTeams:
			var req GetUserTeamsRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetUserTeams(req)
			}
		}

		// ? Send response
		body, _ := json.Marshal(response)
		delivery.Ack(false)
		common.RabbitMQChannel.Publish(
			"",               // * exchange
			delivery.ReplyTo, // * routing key
			false,            // * mandatory
			false,            // * immediate
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: delivery.CorrelationId,
				Body:          body,
			},
		)
	}
}

func handleCreate(req CreateTeamRequest) common.RPCResponse {
	if req.Name == "" {
		return common.RPCResponse{Success: false, Error: "Team name required", StatusCode: 400}
	}
	userReq := map[string]uint{"id": req.UserID}
	resp, rpcErr := common.CallRPC(common.UsersGetByID, userReq)
	if rpcErr != nil {
		return common.RPCResponse{Success: false, Error: "Failed to verify user", StatusCode: 500}
	}
	if !resp.Success {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 404}
	}

	var remoteUserID uint
	var remoteEmail, remoteName string
	if data, ok := resp.Data.(map[string]interface{}); ok {
		if idf, ok := data["id"].(float64); ok {
			remoteUserID = uint(idf)
		}
		if e, ok := data["email"].(string); ok {
			remoteEmail = e
		}
		if n, ok := data["name"].(string); ok {
			remoteName = n
		}
	}

	var createdTeam models.Team
	err := common.DB.Transaction(func(tx *gorm.DB) error {
		if remoteUserID != 0 {
			u := models.User{ID: remoteUserID, Email: remoteEmail, Name: remoteName}
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&u).Error; err != nil {
				return err
			}
		}

		team := models.Team{Name: req.Name}
		if err := tx.Create(&team).Error; err != nil {
			return err
		}

		member := models.TeamMember{
			TeamID: team.ID,
			UserID: req.UserID,
			Role:   "owner",
		}
		if err := tx.Create(&member).Error; err != nil {
			return err
		}

		if err := tx.Preload("Members").First(&team, team.ID).Error; err != nil {
			return err
		}
		createdTeam = team
		return nil
	})

	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to create team", StatusCode: 500}
	}

	// Attach user profiles to members before returning
	attachUsersToMembers(createdTeam.Members)

	fmt.Println("Created team :", createdTeam)
	return toRPC(teamResponse{Success: true, Data: &createdTeam})
}

func handleGetAll() common.RPCResponse {
	var teams []models.Team
	common.DB.Preload("Members").Find(&teams)
	for i := range teams {
		if teams[i].Members == nil {
			teams[i].Members = make([]models.TeamMember, 0)
		}
		// Attach user profiles for members
		attachUsersToMembers(teams[i].Members)
	}
	return toRPC(teamsListResponse{Success: true, Data: teams})
}

func handleGetByID(req GetTeamByIDRequest) common.RPCResponse {
	var team models.Team
	if err := common.DB.Preload("Members").First(&team, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Team not found", StatusCode: 404}
	}
	if team.Members == nil {
		fmt.Println("No members found, initializing empty slice")
		team.Members = make([]models.TeamMember, 0)
	}
	// Attach user profiles for members
	attachUsersToMembers(team.Members)
	return toRPC(teamResponse{Success: true, Data: &team})
}

func handleAddMember(req AddMemberRequest) common.RPCResponse {
	role := req.Role
	if role == "" {
		role = "member"
	}

	var existing models.TeamMember
	if err := common.DB.Where("team_id = ? AND user_id = ?", req.TeamID, req.UserID).First(&existing).Error; err == nil {
		return common.RPCResponse{Success: false, Error: "Member already exists", StatusCode: 409}
	}

	member := models.TeamMember{
		TeamID: req.TeamID,
		UserID: req.UserID,
		Role:   role,
	}
	if err := common.DB.Create(&member).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to add member", StatusCode: 500}
	}

	// Enrich member with user profile
	attachUserProfile(&member)

	common.PublishEvent(common.TeamMemberAdded, map[string]interface{}{
		"teamId": req.TeamID,
		"member": member,
	})

	return toRPC(memberResponse{Success: true, Data: &member})
}

func handleRemoveMember(req RemoveMemberRequest) common.RPCResponse {
	var member models.TeamMember
	if err := common.DB.Where("team_id = ? AND user_id = ?", req.TeamID, req.UserID).First(&member).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Member not found", StatusCode: 404}
	}

	common.DB.Delete(&member)

	common.PublishEvent(common.TeamMemberRemoved, map[string]interface{}{
		"teamId": req.TeamID,
		"userId": req.UserID,
	})

	return toRPC(successResponse{Success: true})
}

func handleGetUserTeams(req GetUserTeamsRequest) common.RPCResponse {
	var members []models.TeamMember
	common.DB.Where("user_id = ?", req.UserID).
		Preload("Team").
		Find(&members)

	teams := make([]UserTeamSummary, len(members))
	for i, member := range members {
		teams[i] = UserTeamSummary{
			ID:        member.Team.ID,
			Name:      member.Team.Name,
			Role:      member.Role,
			CreatedAt: member.Team.CreatedAt,
		}
	}

	return toRPC(userTeamResponse{Success: true, Data: teams})
}

