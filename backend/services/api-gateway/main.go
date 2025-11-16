package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/wetask/backend/pkg/common"
	_ "github.com/wetask/backend/services/api-gateway/docs"
	"github.com/wetask/backend/services/api-gateway/handlers"
)

// @title           WeTask API Gateway
// @version         1.0
// @description     API Gateway for WeTask - A collaborative task management system
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@wetask.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

var upgrader = websocket.Upgrader{
	CheckOrigin: func(router *http.Request) bool {
		return true
	},
}

func main() {
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
			auth.POST("/register", handlers.HandleRegister)
			auth.POST("/login", handlers.HandleLogin)
			auth.POST("/refresh", handlers.HandleRefresh)
		}

		// ? Protected routes
		protected := api.Group("")
		protected.Use(handlers.AuthMiddleware())
		{
			// ? Users
			users := protected.Group("/users")
			{
				users.GET("/me", handlers.HandleGetMe)
				users.GET(":id", handlers.HandleGetUser)
				users.PATCH(":id", handlers.HandleUpdateUser)
			}

			// ? Teams
			teams := protected.Group("/teams")
			{
				teams.GET("", handlers.HandleGetTeams)
				teams.POST("", handlers.HandleCreateTeam)
				teams.GET(":id", handlers.HandleGetTeam)
				teams.POST(":id/members", handlers.HandleAddTeamMember)
				teams.DELETE(":id/members/:userId", handlers.HandleRemoveTeamMember)
			}

			// ? Boards
			boards := protected.Group("/boards")
			{
				boards.GET("", handlers.HandleGetBoards)
				boards.POST("", handlers.HandleCreateBoard)
				boards.GET(":id", handlers.HandleGetBoard)
				boards.PUT(":id", handlers.HandleUpdateBoard)
				boards.DELETE(":id", handlers.HandleDeleteBoard)
			}

			// ? Columns
			columns := protected.Group("/columns")
			{
				columns.POST("", handlers.HandleCreateColumn)
				columns.GET("/board/:boardId", handlers.HandleGetColumns)
				columns.PUT(":id", handlers.HandleUpdateColumn)
				columns.DELETE(":id", handlers.HandleDeleteColumn)
			}

			// ? Tasks
			tasks := protected.Group("/tasks")
			{
				tasks.POST("", handlers.HandleCreateTask)
				tasks.GET(":id", handlers.HandleGetTask)
				tasks.GET("/board/:boardId", handlers.HandleGetTasksByBoard)
				tasks.PUT(":id", handlers.HandleUpdateTask)
				tasks.DELETE(":id", handlers.HandleDeleteTask)
				tasks.PUT(":id/move", handlers.HandleMoveTask)
				tasks.POST(":id/comment", handlers.HandleAddComment)
				tasks.GET(":id/comments", handlers.HandleGetComments)
			}
		}
	}

	// ? WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		handleWebSocket(c, hub)
	})

	// ? Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("API Gateway is running on http: // * localhost:%s/api", port)
	log.Printf("WebSocket is available on ws: // * localhost:%s/ws", port)
	log.Printf("Swagger documentation is available on http: // * localhost:%s/swagger/index.html", port)
	router.Run(":" + port)
}
