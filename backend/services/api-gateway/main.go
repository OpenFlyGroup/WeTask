package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/wetask/backend/pkg/common"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(router *http.Request) bool {
		return true
	},
}

func main() {
	// ? Load env variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// ? Init RabbitMQ
	if err := common.InitRabbitMQ(); err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer common.CloseRabbitMQ()

	// ? Init JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "!!! PUT THIS IN ENV FILE !!!"
	}
	common.InitJWT(jwtSecret)

	// ? Set up event exchange
	err := common.RabbitMQChannel.ExchangeDeclare(
		"events", // * name
		"topic",  // * type
		true,     // * durable
		false,    // * auto-deleted
		false,    // * internal
		false,    // * no-wait
		nil,      // * arguments
	)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}

	// ? Declare events queue
	_, err = common.RabbitMQChannel.QueueDeclare(
		"events_queue", // * name
		true,           // * durable
		false,          // * delete when unused
		false,          // * exclusive
		false,          // * no-wait
		nil,            // * arguments
	)
	if err != nil {
		log.Fatal("Failed to declare events queue:", err)
	}

	// ? Bind queue to exchange for all events
	events := []string{
		common.TaskCreated,
		common.TaskUpdated,
		common.TaskDeleted,
		common.BoardUpdated,
		common.TeamMemberAdded,
		common.TeamMemberRemoved,
	}
	for _, event := range events {
		err = common.RabbitMQChannel.QueueBind(
			"events_queue", // * queue name
			event,          // * routing key
			"events",       // * exchange
			false,
			nil,
		)
		if err != nil {
			log.Fatal("Failed to bind queue:", err)
		}
	}

	// ? Start WebSocket hub
	hub := NewHub()
	go hub.Run()

	// ? Set up routes
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := router.Group("/api")
	{
		// ? Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handleRegister)
			auth.POST("/login", handleLogin)
			auth.POST("/refresh", handleRefresh)
		}

		// ? Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware())
		{
			// ? Users
			users := protected.Group("/users")
			{
				users.GET("/me", handleGetMe)
				users.GET("/:id", handleGetUser)
				users.PATCH("/:id", handleUpdateUser)
			}

			// ? Teams
			teams := protected.Group("/teams")
			{
				teams.GET("", handleGetTeams)
				teams.POST("", handleCreateTeam)
				teams.GET("/:id", handleGetTeam)
				teams.POST("/:id/members", handleAddTeamMember)
				teams.DELETE("/:id/members/:userId", handleRemoveTeamMember)
			}

			// ? Boards
			boards := protected.Group("/boards")
			{
				boards.GET("", handleGetBoards)
				boards.POST("", handleCreateBoard)
				boards.GET("/:id", handleGetBoard)
				boards.PUT("/:id", handleUpdateBoard)
				boards.DELETE("/:id", handleDeleteBoard)
			}

			// ? Columns
			columns := protected.Group("/columns")
			{
				columns.POST("", handleCreateColumn)
				columns.GET("/board/:boardId", handleGetColumns)
				columns.PUT("/:id", handleUpdateColumn)
				columns.DELETE("/:id", handleDeleteColumn)
			}

			// ? Tasks
			tasks := protected.Group("/tasks")
			{
				tasks.POST("", handleCreateTask)
				tasks.GET("/:id", handleGetTask)
				tasks.GET("/board/:boardId", handleGetTasksByBoard)
				tasks.PUT("/:id", handleUpdateTask)
				tasks.DELETE("/:id", handleDeleteTask)
				tasks.PUT("/:id/move", handleMoveTask)
				tasks.POST("/:id/comment", handleAddComment)
				tasks.GET("/:id/comments", handleGetComments)
			}
		}
	}

	// ? WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		handleWebSocket(c, hub)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("API Gateway is running on http: // * localhost:%s/api", port)
	log.Printf("WebSocket is available on ws: // * localhost:%s/ws", port)
	router.Run(":" + port)
}
