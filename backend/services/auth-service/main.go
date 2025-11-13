package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"golang.org/x/crypto/bcrypt"

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

	// ? Migrate auth models
	if err := common.MigrateAuthModels(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// ? Initialize RabbitMQ
	if err := common.InitRabbitMQ(); err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer common.CloseRabbitMQ()

	// ? Initialize JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "SECRET_KEY"
	}
	common.InitJWT(jwtSecret)

	// ? Declare queues
	queues := []string{
		common.AuthRegister,
		common.AuthLogin,
		common.AuthRefresh,
		common.AuthValidate,
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

	log.Println("Auth Service is running...")
	select {} // ? Keep running
}

func handleMessages(queue string, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		var response common.RPCResponse

		switch queue {
		case common.AuthRegister:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleRegister(data)
			}
		case common.AuthLogin:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleLogin(data)
			}
		case common.AuthRefresh:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				refreshToken, _ := data["refreshToken"].(string)
				response = handleRefresh(refreshToken)
			}
		case common.AuthValidate:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				token, _ := data["token"].(string)
				response = handleValidate(token)
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

func handleRegister(data map[string]any) common.RPCResponse {
	email, _ := data["email"].(string)
	password, _ := data["password"].(string)
	name, _ := data["name"].(string)

	if email == "" || password == "" || name == "" {
		return common.RPCResponse{Success: false, Error: "Missing required fields", StatusCode: 400}
	}

	// ? Check if user exists
	var existingUser models.User
	if err := common.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return common.RPCResponse{Success: false, Error: "User already exists", StatusCode: 409}
	}

	// ? Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to hash password", StatusCode: 500}
	}

	// ? Create user
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	if err := common.DB.Create(&user).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to create user", StatusCode: 500}
	}

	// ? Also create user in users service (sync)
	_, syncErr := common.CallRPC(common.UsersCreate, map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
	if syncErr != nil {
		// ? Log error but don't fail registration
		log.Printf("Warning: Failed to sync user to users service: %v", syncErr)
	}

	// ? Generate tokens
	tokens, err := generateTokens(user.ID)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to generate tokens", StatusCode: 500}
	}

	return common.RPCResponse{
		Success: true,
		Data: map[string]any{
			"user": map[string]any{
				"id":        user.ID,
				"email":     user.Email,
				"name":      user.Name,
				"createdAt": user.CreatedAt,
			},
			"accessToken":  tokens["accessToken"],
			"refreshToken": tokens["refreshToken"],
		},
	}
}

func handleLogin(data map[string]any) common.RPCResponse {
	email, _ := data["email"].(string)
	password, _ := data["password"].(string)

	if email == "" || password == "" {
		return common.RPCResponse{Success: false, Error: "Missing email or password", StatusCode: 400}
	}

	// ? Find user
	var user models.User
	if err := common.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Invalid credentials", StatusCode: 401}
	}

	// ? Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return common.RPCResponse{Success: false, Error: "Invalid credentials", StatusCode: 401}
	}

	// ? Ensure user exists in users service (sync if missing)
	_, syncErr := common.CallRPC(common.UsersCreate, map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
	if syncErr != nil {
		// ? Log error but don't fail login
		log.Printf("Warning: Failed to sync user to users service: %v", syncErr)
	}

	// ? Generate tokens
	tokens, err := generateTokens(user.ID)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to generate tokens", StatusCode: 500}
	}

	return common.RPCResponse{
		Success: true,
		Data: map[string]any{
			"user": map[string]any{
				"id":        user.ID,
				"email":     user.Email,
				"name":      user.Name,
				"createdAt": user.CreatedAt,
			},
			"accessToken":  tokens["accessToken"],
			"refreshToken": tokens["refreshToken"],
		},
	}
}

func handleRefresh(refreshToken string) common.RPCResponse {
	if refreshToken == "" {
		return common.RPCResponse{Success: false, Error: "Refresh token required", StatusCode: 400}
	}

	// ? Find refresh token
	var token models.RefreshToken
	if err := common.DB.Where("token = ? AND expires_at > ?", refreshToken, time.Now()).First(&token).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Invalid refresh token", StatusCode: 401}
	}

	// ? Generate new tokens
	tokens, err := generateTokens(token.UserID)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to generate tokens", StatusCode: 500}
	}

	// ? Delete old refresh token
	common.DB.Delete(&token)

	return common.RPCResponse{
		Success: true,
		Data: map[string]any{
			"accessToken":  tokens["accessToken"],
			"refreshToken": tokens["refreshToken"],
		},
	}
}

func handleValidate(token string) common.RPCResponse {
	if token == "" {
		return common.RPCResponse{Success: false, Error: "Token required", StatusCode: 400}
	}

	claims, err := common.ValidateToken(token)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Invalid token", StatusCode: 401}
	}

	// ? Find user
	var user models.User
	if err := common.DB.First(&user, claims.UserID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 401}
	}

	return common.RPCResponse{
		Success: true,
		Data: map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	}
}

func generateTokens(userID uint) (map[string]string, error) {
	accessToken, err := common.GenerateToken(userID, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	// ? Generate refresh token
	refreshTokenBytes := make([]byte, 64)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		return nil, err
	}
	refreshToken := fmt.Sprintf("%x", refreshTokenBytes)

	// ? Save refresh token
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	rt := models.RefreshToken{
		Token:     refreshToken,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	common.DB.Create(&rt)

	return map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}, nil
}
