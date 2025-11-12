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

	// ? Initialize database (users service has its own user data)
	if err := common.InitPostgreSQL(); err != nil {
		log.Fatal("Failed to initialize PostgreSQL:", err)
	}
	
	// ? Migrate users models (user profile data)
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
		msgs, err := common.RabbitMQChannel.Consume(
			queue, // ? queue
			"",    // ? consumer
			false, // ? auto-ack
			false, // ? exclusive
			false, // ? no-local
			false, // ? no-wait
			nil,   // ? args
		)
		if err != nil {
			log.Fatal("Failed to register consumer:", err)
		}

		go handleMessages(queue, msgs)
	}

	log.Println("Users Service is running...")
	select {} // ? Keep running
}

func handleMessages(queue string, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		var response common.RPCResponse

		switch queue {
		case common.UsersGetByID:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				id, _ := data["id"].(float64)
				response = handleGetByID(uint(id))
			}
		case common.UsersGetByEmail:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				email, _ := data["email"].(string)
				response = handleGetByEmail(email)
			}
		case common.UsersUpdate:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleUpdate(data)
			}
		case common.UsersGetMe:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				userID, _ := data["userId"].(float64)
				response = handleGetMe(uint(userID))
			}
		}

		// ? Send response
		responseBody, _ := json.Marshal(response)
		d.Ack(false)
		common.RabbitMQChannel.Publish(
			"",        // ? exchange
			d.ReplyTo, // ? routing key
			false,     // ? mandatory
			false,     // ? immediate
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          responseBody,
			},
		)
	}
}

func handleGetByID(id uint) common.RPCResponse {
	var user models.User
	if err := common.DB.First(&user, id).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 404}
	}

	return common.RPCResponse{
		Success: true,
		Data: map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	}
}

func handleGetByEmail(email string) common.RPCResponse {
	var user models.User
	if err := common.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 404}
	}

	return common.RPCResponse{
		Success: true,
		Data: map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	}
}

func handleUpdate(data map[string]interface{}) common.RPCResponse {
	id, ok := data["id"].(float64)
	if !ok {
		return common.RPCResponse{Success: false, Error: "User ID required", StatusCode: 400}
	}

	var user models.User
	if err := common.DB.First(&user, uint(id)).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 404}
	}

	// ? Update fields
	if name, ok := data["name"].(string); ok && name != "" {
		user.Name = name
	}
	if email, ok := data["email"].(string); ok && email != "" {
		// ? Check if email already exists
		var existingUser models.User
		if err := common.DB.Where("email = ? AND id != ?", email, id).First(&existingUser).Error; err == nil {
			return common.RPCResponse{Success: false, Error: "Email already exists", StatusCode: 409}
		}
		user.Email = email
	}

	if err := common.DB.Save(&user).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to update user", StatusCode: 500}
	}

	return common.RPCResponse{
		Success: true,
		Data: map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	}
}

func handleGetMe(userID uint) common.RPCResponse {
	var user models.User
	if err := common.DB.First(&user, userID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 404}
	}

	return common.RPCResponse{
		Success: true,
		Data: map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	}
}

