package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/wetask/backend/pkg/common"
	"github.com/wetask/backend/pkg/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	if err := common.InitPostgreSQL(); err != nil {
		log.Fatal("Failed to initialize PostgreSQL:", err)
	}
	
	// Migrate tasks models
	if err := common.MigrateTasksModels(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize MongoDB
	if err := common.InitMongoDB(); err != nil {
		log.Fatal("Failed to initialize MongoDB:", err)
	}

	// Initialize RabbitMQ
	if err := common.InitRabbitMQ(); err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer common.CloseRabbitMQ()

	// Declare queues
	queues := []string{
		common.TasksCreate,
		common.TasksGetByID,
		common.TasksGetByBoard,
		common.TasksUpdate,
		common.TasksDelete,
		common.TasksMove,
		common.TasksAddComment,
		common.TasksGetComments,
	}

	for _, queue := range queues {
		_, err := common.RabbitMQChannel.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			log.Fatal("Failed to declare queue:", err)
		}
	}

	// Start consuming messages
	for _, queue := range queues {
		msgs, err := common.RabbitMQChannel.Consume(
			queue, // queue
			"",    // consumer
			false, // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			log.Fatal("Failed to register consumer:", err)
		}

		go handleMessages(queue, msgs)
	}

	log.Println("Tasks Service is running...")
	select {} // Keep running
}

func handleMessages(queue string, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		var response common.RPCResponse

		switch queue {
		case common.TasksCreate:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleCreateTask(data)
			}
		case common.TasksGetByID:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				id, _ := data["id"].(float64)
				response = handleGetTaskByID(uint(id))
			}
		case common.TasksGetByBoard:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				boardID, _ := data["boardId"].(float64)
				response = handleGetTasksByBoard(uint(boardID))
			}
		case common.TasksUpdate:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleUpdateTask(data)
			}
		case common.TasksDelete:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				id, _ := data["id"].(float64)
				response = handleDeleteTask(uint(id))
			}
		case common.TasksMove:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleMoveTask(data)
			}
		case common.TasksAddComment:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleAddComment(data)
			}
		case common.TasksGetComments:
			var data map[string]interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				taskID, _ := data["taskId"].(float64)
				response = handleGetComments(uint(taskID))
			}
		}

		// Send response
		responseBody, _ := json.Marshal(response)
		d.Ack(false)
		common.RabbitMQChannel.Publish(
			"",        // exchange
			d.ReplyTo, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          responseBody,
			},
		)
	}
}

func handleCreateTask(data map[string]interface{}) common.RPCResponse {
	title, _ := data["title"].(string)
	description, _ := data["description"].(string)
	columnID, _ := data["columnId"].(float64)
	assignedTo, ok := data["assignedTo"].(float64)
	priority, _ := data["priority"].(string)
	if priority == "" {
		priority = "medium"
	}

	if title == "" {
		return common.RPCResponse{Success: false, Error: "Task title required", StatusCode: 400}
	}

	task := models.Task{
		Title:       title,
		Status:      "todo",
		ColumnID:    uint(columnID),
		Priority:    &priority,
	}

	if description != "" {
		task.Description = &description
	}
	if ok {
		assignedToUint := uint(assignedTo)
		task.AssignedTo = &assignedToUint
	}

	if err := common.DB.Create(&task).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to create task", StatusCode: 500}
	}

	// Get column to find board ID
	var column models.Column
	common.DB.First(&column, task.ColumnID)

	// Publish event
	common.PublishEvent(common.TaskCreated, map[string]interface{}{
		"boardId": column.BoardID,
		"task":    task,
	})

	return common.RPCResponse{
		Success: true,
		Data:    task,
	}
}

func handleGetTaskByID(id uint) common.RPCResponse {
	var task models.Task
	if err := common.DB.Preload("Column").Preload("User").First(&task, id).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Task not found", StatusCode: 404}
	}

	return common.RPCResponse{
		Success: true,
		Data:    task,
	}
}

func handleGetTasksByBoard(boardID uint) common.RPCResponse {
	// Get all columns for this board
	var columns []models.Column
	common.DB.Where("board_id = ?", boardID).Find(&columns)

	columnIDs := make([]uint, len(columns))
	for i, c := range columns {
		columnIDs[i] = c.ID
	}

	var tasks []models.Task
	common.DB.Where("column_id IN ?", columnIDs).Preload("User").Find(&tasks)

	return common.RPCResponse{
		Success: true,
		Data:    tasks,
	}
}

func handleUpdateTask(data map[string]interface{}) common.RPCResponse {
	id, _ := data["id"].(float64)

	var task models.Task
	if err := common.DB.First(&task, uint(id)).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Task not found", StatusCode: 404}
	}

	if title, ok := data["title"].(string); ok && title != "" {
		task.Title = title
	}
	if description, ok := data["description"].(string); ok {
		task.Description = &description
	}
	if priority, ok := data["priority"].(string); ok && priority != "" {
		task.Priority = &priority
	}
	if assignedTo, ok := data["assignedTo"].(float64); ok {
		assignedToUint := uint(assignedTo)
		task.AssignedTo = &assignedToUint
	}

	if err := common.DB.Save(&task).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to update task", StatusCode: 500}
	}

	// Get column to find board ID
	var column models.Column
	common.DB.First(&column, task.ColumnID)

	// Publish event
	common.PublishEvent(common.TaskUpdated, map[string]interface{}{
		"boardId": column.BoardID,
		"task":    task,
	})

	return common.RPCResponse{
		Success: true,
		Data:    task,
	}
}

func handleDeleteTask(id uint) common.RPCResponse {
	var task models.Task
	if err := common.DB.First(&task, id).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Task not found", StatusCode: 404}
	}

	// Get column to find board ID
	var column models.Column
	common.DB.First(&column, task.ColumnID)
	boardID := column.BoardID

	common.DB.Delete(&task)

	// Publish event
	common.PublishEvent(common.TaskDeleted, map[string]interface{}{
		"boardId": boardID,
		"taskId":  id,
	})

	return common.RPCResponse{Success: true}
}

func handleMoveTask(data map[string]interface{}) common.RPCResponse {
	id, _ := data["id"].(float64)
	columnID, _ := data["columnId"].(float64)

	var task models.Task
	if err := common.DB.First(&task, uint(id)).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Task not found", StatusCode: 404}
	}

	// Get old column for board ID
	var oldColumn models.Column
	common.DB.First(&oldColumn, task.ColumnID)
	boardID := oldColumn.BoardID

	task.ColumnID = uint(columnID)
	if err := common.DB.Save(&task).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to move task", StatusCode: 500}
	}

	// Publish event
	common.PublishEvent(common.TaskUpdated, map[string]interface{}{
		"boardId": boardID,
		"task":    task,
	})

	return common.RPCResponse{
		Success: true,
		Data:    task,
	}
}

func handleAddComment(data map[string]interface{}) common.RPCResponse {
	taskID, _ := data["taskId"].(float64)
	userID, _ := data["userId"].(float64)
	message, _ := data["message"].(string)

	if message == "" {
		return common.RPCResponse{Success: false, Error: "Comment message required", StatusCode: 400}
	}

	comment := models.Comment{
		ID:        primitive.NewObjectID().Hex(),
		TaskID:    uint(taskID),
		UserID:    uint(userID),
		Message:   message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := common.MongoDB.Collection("comments")
	_, err := collection.InsertOne(context.Background(), comment)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to add comment", StatusCode: 500}
	}

	return common.RPCResponse{
		Success: true,
		Data:    comment,
	}
}

func handleGetComments(taskID uint) common.RPCResponse {
	collection := common.MongoDB.Collection("comments")
	cursor, err := collection.Find(context.Background(), bson.M{"taskId": taskID}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}}))
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to get comments", StatusCode: 500}
	}
	defer cursor.Close(context.Background())

	var comments []models.Comment
	if err := cursor.All(context.Background(), &comments); err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to decode comments", StatusCode: 500}
	}

	return common.RPCResponse{
		Success: true,
		Data:    comments,
	}
}

