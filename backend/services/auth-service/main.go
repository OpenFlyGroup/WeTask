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
			var req RegisterRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleRegister(req)
			}
		case common.AuthLogin:
			var req LoginRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleLogin(req)
			}
		case common.AuthRefresh:
			var req RefreshRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleRefresh(req)
			}
		case common.AuthValidate:
			var req ValidateRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleValidate(req)
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

func handleRegister(req RegisterRequest) common.RPCResponse {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return common.RPCResponse{Success: false, Error: "Missing required fields", StatusCode: 400}
	}

	// ? Check if user exists (auth-side)
	var existingUser models.AuthUser
	if err := common.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return common.RPCResponse{Success: false, Error: "User already exists", StatusCode: 409}
	}

	// ? Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to hash password", StatusCode: 500}
	}

	// ? Create minimal auth user (do not store rich profile data here)
	user := models.AuthUser{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := common.DB.Create(&user).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to create user", StatusCode: 500}
	}

	// ? Also create user in users service (sync the rich profile)
	if resp, syncErr := common.CallRPC(common.UsersCreate, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  req.Name,
	}); syncErr != nil || resp == nil || !resp.Success {
		common.DB.Delete(&user)
		errMsg := "Failed to create user profile"
		if syncErr != nil {
			log.Printf("Error: failed to sync user to users service: %v", syncErr)
			errMsg = syncErr.Error()
		} else if resp != nil && !resp.Success {
			log.Printf("Error: users service returned failure: %v", resp.Error)
			errMsg = resp.Error
		}
		return common.RPCResponse{Success: false, Error: errMsg, StatusCode: 500}
	}

	// ? Generate tokens
	tokens, err := generateTokens(user.ID)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to generate tokens", StatusCode: 500}
	}

	// ? Try to fetch rich profile from users service for response
	var userResp *UserResponse
	if resp, err := common.CallRPC(common.UsersGetByID, map[string]interface{}{"id": user.ID}); err == nil && resp != nil && resp.Success {
		if data, ok := resp.Data.(map[string]interface{}); ok {
			userResp = &UserResponse{
				ID:        uint(data["id"].(float64)),
				Email:     data["email"].(string),
				Name:      data["name"].(string),
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}
		}
	}

	if userResp == nil {
		userResp = &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}
	}

	return common.RPCResponse{
		Success: true,
		Data: AuthResponse{
			User:         userResp,
			AccessToken:  tokens["accessToken"],
			RefreshToken: tokens["refreshToken"],
		},
	}
}

func handleLogin(req LoginRequest) common.RPCResponse {
	if req.Email == "" || req.Password == "" {
		return common.RPCResponse{Success: false, Error: "Missing email or password", StatusCode: 400}
	}

	// ? Find auth user
	var user models.AuthUser
	if err := common.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Invalid credentials", StatusCode: 401}
	}

	// ? Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return common.RPCResponse{Success: false, Error: "Invalid credentials", StatusCode: 401}
	}

	// ? Ensure user exists in users service (sync if missing). Auth doesn't hold rich profile.
	_, syncErr := common.CallRPC(common.UsersCreate, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
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

	// ? Try to fetch rich profile from users service for response
	var userResp *UserResponse
	if resp, err := common.CallRPC(common.UsersGetByID, map[string]interface{}{"id": user.ID}); err == nil && resp != nil && resp.Success {
		if data, ok := resp.Data.(map[string]interface{}); ok {
			userResp = &UserResponse{
				ID:        uint(data["id"].(float64)),
				Email:     data["email"].(string),
				Name:      data["name"].(string),
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}
		}
	}

	if userResp == nil {
		userResp = &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}
	}

	return common.RPCResponse{
		Success: true,
		Data: AuthResponse{
			User:         userResp,
			AccessToken:  tokens["accessToken"],
			RefreshToken: tokens["refreshToken"],
		},
	}
}

func handleRefresh(req RefreshRequest) common.RPCResponse {
	if req.RefreshToken == "" {
		return common.RPCResponse{Success: false, Error: "Refresh token required", StatusCode: 400}
	}

	// ? Find refresh token
	var token models.RefreshToken
	if err := common.DB.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&token).Error; err != nil {
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
		Data: RefreshResponse{
			AccessToken:  tokens["accessToken"],
			RefreshToken: tokens["refreshToken"],
		},
	}
}

func handleValidate(req ValidateRequest) common.RPCResponse {
	if req.Token == "" {
		return common.RPCResponse{Success: false, Error: "Token required", StatusCode: 400}
	}

	claims, err := common.ValidateToken(req.Token)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Invalid token", StatusCode: 401}
	}
	// ? Ensure auth user exists
	var authUser models.AuthUser
	if err := common.DB.First(&authUser, claims.UserID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "User not found", StatusCode: 401}
	}

	// ? Ensure token's issued-at matches user's LastAccessTokenAt
	if claims.IssuedAt == nil {
		return common.RPCResponse{Success: false, Error: "Invalid token (missing iat)", StatusCode: 401}
	}
	if claims.IssuedAt.Time.Unix() != authUser.LastAccessTokenAt.Unix() {
		return common.RPCResponse{Success: false, Error: "Token has been invalidated", StatusCode: 401}
	}

	// ? Fetch rich profile from users service (if available)
	if resp, err := common.CallRPC(common.UsersGetByID, map[string]interface{}{"id": authUser.ID}); err == nil && resp != nil && resp.Success {
		if data, ok := resp.Data.(map[string]interface{}); ok {
			return common.RPCResponse{
				Success: true,
				Data: ValidateResponse{
					ID:    uint(data["id"].(float64)),
					Email: data["email"].(string),
					Name:  data["name"].(string),
				},
			}
		}
	}

	return common.RPCResponse{
		Success: true,
		Data: ValidateResponse{
			ID:    authUser.ID,
			Email: authUser.Email,
		},
	}
}

func generateTokens(userID uint) (map[string]string, error) {
	now := time.Now()
	accessToken, err := common.GenerateToken(userID, 15*time.Minute, now)
	if err != nil {
		return nil, err
	}

	// ? Update user's token generation timestamp
	if err := common.DB.Model(&models.AuthUser{}).Where("id = ?", userID).Update("last_access_token_at", now).Error; err != nil {
		log.Printf("Warning: Failed to update token timestamp: %v", err)
	}

	// ? Delete any old refresh tokens for this user (invalidate previous sessions)
	if err := common.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error; err != nil {
		log.Printf("Warning: Failed to invalidate old refresh tokens: %v", err)
	}

	// ? Generate refresh token
	refreshTokenBytes := make([]byte, 64)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		return nil, err
	}
	refreshToken := fmt.Sprintf("%x", refreshTokenBytes)

	// ? Save new refresh token
	expiresAt := now.Add(7 * 24 * time.Hour)
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
