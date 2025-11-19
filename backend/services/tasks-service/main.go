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

// attachUserToTask populates Task.User by calling the Users service.
func attachUserToTask(t *models.Task) {
	if t == nil || t.AssignedTo == nil {
		return
	}
	if rpcResp, err := common.CallRPC(common.UsersGetByID, map[string]interface{}{"id": *t.AssignedTo}); err == nil && rpcResp != nil && rpcResp.Success {
		if data, ok := rpcResp.Data.(map[string]interface{}); ok {
			var u models.User
			if id, ok := data["id"].(float64); ok {
				u.ID = uint(id)
			}
			if email, ok := data["email"].(string); ok {
				u.Email = email
			}
			if name, ok := data["name"].(string); ok {
				u.Name = name
			}
			t.User = &u
		}
	}
}

func attachUsersToTasks(tasks []models.Task) {
	for i := range tasks {
		attachUserToTask(&tasks[i])
	}
}

func main() {
	// ? Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// ? Initialize PostgreSQL and migrate task models
	if err := common.InitPostgreSQL(); err != nil {
		log.Fatal("Failed to initialize PostgreSQL:", err)
	}
	if err := common.MigrateTasksModels(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// ? Initialize MongoDB (comments storage)
	if err := common.InitMongoDB(); err != nil {
		log.Fatal("Failed to initialize MongoDB:", err)
	}

	// ? Initialize RabbitMQ
	if err := common.InitRabbitMQ(); err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer common.CloseRabbitMQ()

	// ? Declare queues to consume
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

	// ? Start consuming messages for each queue
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

	log.Println("Tasks Service is running...")
	select {} // ? Keep running
}

// ? Handles incoming messages for task-related queues
func handleMessages(queue string, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		var response common.RPCResponse

		switch queue {
		case common.TasksCreate:
			var req CreateTaskRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleCreateTask(req)
			}
		case common.TasksGetByID:
			var req GetTaskByIDRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetTaskByID(req)
			}
		case common.TasksGetByBoard:
			var req GetTasksByBoardRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetTasksByBoard(req)
			}
		case common.TasksUpdate:
			var req UpdateTaskRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleUpdateTask(req)
			}
		case common.TasksDelete:
			var req DeleteTaskRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleDeleteTask(req)
			}
		case common.TasksMove:
			var req MoveTaskRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleMoveTask(req)
			}
		case common.TasksAddComment:
			var req AddCommentRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleAddComment(req)
			}
		case common.TasksGetComments:
			var req GetCommentsRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetComments(req)
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

func handleCreateTask(req CreateTaskRequest) common.RPCResponse {
	if req.Title == "" {
		return common.RPCResponse{Success: false, Error: "Task title required", StatusCode: 400}
	}

	priority := "medium"
	if req.Priority != nil && *req.Priority != "" {
		priority = *req.Priority
	}

	task := models.Task{
		Title:       req.Title,
		Status:      "todo",
		ColumnID:    req.ColumnID,
		Priority:    &priority,
		Description: req.Description,
		AssignedTo:  req.AssignedTo,
	}

	if err := common.DB.Create(&task).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to create task", StatusCode: 500}
	}

	var column models.Column
	common.DB.First(&column, task.ColumnID)

	// Enrich task with assigned user profile (if any)
	attachUserToTask(&task)

	common.PublishEvent(common.TaskCreated, map[string]interface{}{
		"boardId": column.BoardID,
		"task":    task,
	})

	return toRPC(taskResponse{Success: true, Data: &task})
}

func handleGetTaskByID(req GetTaskByIDRequest) common.RPCResponse {
	var task models.Task
	if err := common.DB.Preload("Column").First(&task, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Task not found", StatusCode: 404}
	}
	attachUserToTask(&task)
	return toRPC(taskResponse{Success: true, Data: &task})
}

func handleGetTasksByBoard(req GetTasksByBoardRequest) common.RPCResponse {
	var columns []models.Column
	common.DB.Where("board_id = ?", req.BoardID).Find(&columns)

	columnIDs := make([]uint, len(columns))
	for i, c := range columns {
		columnIDs[i] = c.ID
	}

	var tasks []models.Task
	common.DB.Where("column_id IN ?", columnIDs).Find(&tasks)

	// Enrich tasks with assigned user profiles
	attachUsersToTasks(tasks)

	return toRPC(tasksListResponse{Success: true, Data: tasks})
}

func handleUpdateTask(req UpdateTaskRequest) common.RPCResponse {
	var task models.Task
	if err := common.DB.First(&task, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Task not found", StatusCode: 404}
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.Priority != nil && *req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.AssignedTo != nil {
		task.AssignedTo = req.AssignedTo
	}

	if err := common.DB.Save(&task).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to update task", StatusCode: 500}
	}

	var column models.Column
	common.DB.First(&column, task.ColumnID)

	// Enrich task with assigned user profile
	attachUserToTask(&task)

	common.PublishEvent(common.TaskUpdated, map[string]interface{}{
		"boardId": column.BoardID,
		"task":    task,
	})

	return toRPC(taskResponse{Success: true, Data: &task})
}

func handleDeleteTask(req DeleteTaskRequest) common.RPCResponse {
	var task models.Task
	if err := common.DB.First(&task, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Task not found", StatusCode: 404}
	}

	var column models.Column
	common.DB.First(&column, task.ColumnID)
	boardID := column.BoardID

	common.DB.Delete(&task)

	common.PublishEvent(common.TaskDeleted, map[string]interface{}{
		"boardId": boardID,
		"taskId":  req.ID,
	})

	return toRPC(successResponse{Success: true})
}

func handleMoveTask(req MoveTaskRequest) common.RPCResponse {
	var task models.Task
	if err := common.DB.First(&task, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Task not found", StatusCode: 404}
	}

	var oldColumn models.Column
	common.DB.First(&oldColumn, task.ColumnID)
	boardID := oldColumn.BoardID

	task.ColumnID = req.ColumnID
	if err := common.DB.Save(&task).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to move task", StatusCode: 500}
	}

	// Enrich task before publishing
	attachUserToTask(&task)

	common.PublishEvent(common.TaskUpdated, map[string]interface{}{
		"boardId": boardID,
		"task":    task,
	})

	return toRPC(taskResponse{Success: true, Data: &task})
}

func handleAddComment(req AddCommentRequest) common.RPCResponse {
	if req.Message == "" {
		return common.RPCResponse{Success: false, Error: "Comment message required", StatusCode: 400}
	}

	comment := models.Comment{
		ID:        primitive.NewObjectID().Hex(),
		TaskID:    req.TaskID,
		UserID:    req.UserID,
		Message:   req.Message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := common.MongoDB.Collection("comments")
	_, err := collection.InsertOne(context.Background(), comment)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to add comment", StatusCode: 500}
	}

	return toRPC(commentResponse{Success: true, Data: &comment})
}

func handleGetComments(req GetCommentsRequest) common.RPCResponse {
	collection := common.MongoDB.Collection("comments")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{"taskId": req.TaskID},
		options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}}),
	)
	if err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to get comments", StatusCode: 500}
	}
	defer cursor.Close(context.Background())

	var comments []models.Comment
	if err := cursor.All(context.Background(), &comments); err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to decode comments", StatusCode: 500}
	}

	return toRPC(commentsListResponse{Success: true, Data: comments})
}
