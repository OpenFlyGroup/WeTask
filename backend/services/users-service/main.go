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

	// ? Initialize database and migrate models
	if err := common.InitPostgreSQL(); err != nil {
		log.Fatal("Failed to initialize PostgreSQL:", err)
	}
	if err := common.MigrateUsersModels(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// ? Initialize RabbitMQ
	if err := common.InitRabbitMQ(); err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer common.CloseRabbitMQ()

	// ? Declare queues
	queues := []string{
		common.UsersCreate,
		common.UsersGetByID,
		common.UsersGetByEmail,
		common.UsersUpdate,
		common.UsersGetMe,
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

	log.Println("Users Service is running...")
	select {} // ? Keep running
}

// ? Handles incoming messages for user-related queues
func handleMessages(queue string, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		var response common.RPCResponse

		switch queue {
		case common.UsersCreate:
			var req CreateUserRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleCreate(req)
			}

		case common.UsersGetByID:
			var req GetUserByIDRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetByID(req)
			}

		case common.UsersGetByEmail:
			var req GetUserByEmailRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetByEmail(req)
			}

		case common.UsersUpdate:
			var req UpdateUserRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleUpdate(req)
			}

		case common.UsersGetMe:
			var req GetMeRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetMe(req)
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

func handleCreate(req CreateUserRequest) common.RPCResponse {
	if req.Email == "" || req.Name == "" {
		return common.RPCResponse{Success: false, Error: "Missing required fields", StatusCode: 400}
	}

	var existingUser models.User
	err := common.DB.First(&existingUser, req.ID).Error
	if err == nil {
		// ? User exists — update if needed
		updated := false
		if existingUser.Email != req.Email {
			existingUser.Email = req.Email
			updated = true
		}
		if existingUser.Name != req.Name {
			existingUser.Name = req.Name
			updated = true
		}
		if updated {
			if saveErr := common.DB.Save(&existingUser).Error; saveErr != nil {
				return common.RPCResponse{Success: false, Error: "Failed to update user", StatusCode: 500}
			}
		}
	} else {
		// ? Create new user
		user := models.User{
			ID:    req.ID,
			Email: req.Email,
			Name:  req.Name,
		}
		if createErr := common.DB.Create(&user).Error; createErr != nil {
			return common.RPCResponse{Success: false, Error: "Failed to create user", StatusCode: 500}
		}
		existingUser = user
	}

	return toRPC(userResponse{
		Success: true,
		Data:    &existingUser,
	})
}

func handleGetByID(req GetUserByIDRequest) common.RPCResponse {
	var user models.User
	if err := common.DB.First(&user, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 404}
	}
	return toRPC(userResponse{Success: true, Data: &user})
}

func handleGetByEmail(req GetUserByEmailRequest) common.RPCResponse {
	var user models.User
	if err := common.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 404}
	}
	return toRPC(userResponse{Success: true, Data: &user})
}

func handleUpdate(req UpdateUserRequest) common.RPCResponse {
	var user models.User
	if err := common.DB.First(&user, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 404}
	}

	if req.Name != nil && *req.Name != "" {
		user.Name = *req.Name
	}
	if req.Email != nil && *req.Email != "" {
		// ? Check if email already exists (excluding current user)
		var existing models.User
		if err := common.DB.Where("email = ? AND id != ?", *req.Email, req.ID).First(&existing).Error; err == nil {
			return common.RPCResponse{Success: false, Error: "Email already exists", StatusCode: 409}
		}
		user.Email = *req.Email
	}

	if err := common.DB.Save(&user).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to update user", StatusCode: 500}
	}

	return toRPC(userResponse{Success: true, Data: &user})
}

func handleGetMe(req GetMeRequest) common.RPCResponse {
	var user models.User
	if err := common.DB.First(&user, req.UserID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 404}
	}
	return toRPC(userResponse{Success: true, Data: &user})
}
