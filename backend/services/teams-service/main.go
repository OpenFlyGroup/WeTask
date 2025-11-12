package main

import (
	"encoding/json"
	"log"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"

	"github.com/wetask/backend/pkg/common"
	"github.com/wetask/backend/pkg/models"
)

func main() {
	// ? Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// ? Initialize database
	if err := common.InitPostgreSQL(); err != nil {
		log.Fatal("Failed to initialize PostgreSQL:", err)
	}
	
	// ? Migrate teams models
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
		msgs, err := common.RabbitMQChannel.Consume(
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

		go handleMessages(queue, msgs)
	}

	log.Println("Teams Service is running...")
	select {} // ? Keep running
}

func handleMessages(queue string, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		var response common.RPCResponse

		switch queue {
		case common.TeamsCreate:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleCreate(data)
			}
		case common.TeamsGetAll:
			response = handleGetAll()
		case common.TeamsGetByID:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				id, _ := data["id"].(float64)
				response = handleGetByID(uint(id))
			}
		case common.TeamsAddMember:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleAddMember(data)
			}
		case common.TeamsRemoveMember:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleRemoveMember(data)
			}
		case common.TeamsGetUserTeams:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				userID, _ := data["userId"].(float64)
				response = handleGetUserTeams(uint(userID))
			}
		}

		// ? Send response
		responseBody, _ := json.Marshal(response)
		d.Ack(false)
		common.RabbitMQChannel.Publish(
			"",        // * exchange
			d.ReplyTo, // * routing key
			false,     // * mandatory
			false,     // * immediate
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          responseBody,
			},
		)
	}
}

func handleCreate(data map[string]interface{}) common.RPCResponse {
	name, _ := data["name"].(string)
	userID, _ := data["userId"].(float64)

	if name == "" {
		return common.RPCResponse{Success: false, Error: "Team name required", StatusCode: 400}
	}

	team := models.Team{Name: name}
	if err := common.DB.Create(&team).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to create team", StatusCode: 500}
	}

	// ? Add creator as owner
	member := models.TeamMember{
		TeamID: team.ID,
		UserID: uint(userID),
		Role:   "owner",
	}
	common.DB.Create(&member)

	// ? Load members
	common.DB.Preload("Members.User").First(&team, team.ID)

	return common.RPCResponse{
		Success: true,
		Data:    team,
	}
}

func handleGetAll() common.RPCResponse {
	var teams []models.Team
	common.DB.Preload("Members.User").Find(&teams)

	return common.RPCResponse{
		Success: true,
		Data:    teams,
	}
}

func handleGetByID(id uint) common.RPCResponse {
	var team models.Team
	if err := common.DB.Preload("Members.User").First(&team, id).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Team not found", StatusCode: 404}
	}

	return common.RPCResponse{
		Success: true,
		Data:    team,
	}
}

func handleAddMember(data map[string]interface{}) common.RPCResponse {
	teamID, _ := data["teamId"].(float64)
	userID, _ := data["userId"].(float64)
	role, _ := data["role"].(string)
	if role == "" {
		role = "member"
	}

	// ? Check if member already exists
	var existing models.TeamMember
	if err := common.DB.Where("team_id = ? AND user_id = ?", uint(teamID), uint(userID)).First(&existing).Error; err == nil {
		return common.RPCResponse{Success: false, Error: "Member already exists", StatusCode: 409}
	}

	member := models.TeamMember{
		TeamID: uint(teamID),
		UserID: uint(userID),
		Role:   role,
	}
	if err := common.DB.Create(&member).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to add member", StatusCode: 500}
	}

	common.DB.Preload("User").First(&member, member.ID)

	// ? Publish event
	common.PublishEvent(common.TeamMemberAdded, map[string]interface{}{
		"teamId": teamID,
		"member": member,
	})

	return common.RPCResponse{
		Success: true,
		Data:    member,
	}
}

func handleRemoveMember(data map[string]interface{}) common.RPCResponse {
	teamID, _ := data["teamId"].(float64)
	userID, _ := data["userId"].(float64)

	var member models.TeamMember
	if err := common.DB.Where("team_id = ? AND user_id = ?", uint(teamID), uint(userID)).First(&member).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Member not found", StatusCode: 404}
	}

	common.DB.Delete(&member)

	// ? Publish event
	common.PublishEvent(common.TeamMemberRemoved, map[string]interface{}{
		"teamId": teamID,
		"userId": userID,
	})

	return common.RPCResponse{Success: true}
}

func handleGetUserTeams(userID uint) common.RPCResponse {
	var members []models.TeamMember
	common.DB.Where("user_id = ?", userID).Preload("Team").Preload("User").Find(&members)

	teams := make([]map[string]interface{}, len(members))
	for i, member := range members {
		teams[i] = map[string]interface{}{
			"id":        member.Team.ID,
			"name":      member.Team.Name,
			"role":      member.Role,
			"createdAt": member.Team.CreatedAt,
		}
	}

	return common.RPCResponse{
		Success: true,
		Data:    teams,
	}
}

