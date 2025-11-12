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
	
	// ? Migrate boards models
	if err := common.MigrateBoardsModels(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// ? Initialize RabbitMQ
	if err := common.InitRabbitMQ(); err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer common.CloseRabbitMQ()

	// ? Declare queues
	queues := []string{
		common.BoardsCreate,
		common.BoardsGetAll,
		common.BoardsGetByID,
		common.BoardsUpdate,
		common.BoardsDelete,
		common.BoardsGetByTeam,
		common.ColumnsCreate,
		common.ColumnsGetByBoard,
		common.ColumnsUpdate,
		common.ColumnsDelete,
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

	log.Println("Boards Service is running...")
	select {} // ? Keep running
}

func handleMessages(queue string, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		var response common.RPCResponse

		switch queue {
		case common.BoardsCreate:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleCreateBoard(data)
			}
		case common.BoardsGetAll:
			var data map[string]any
			json.Unmarshal(d.Body, &data)
			userID, _ := data["userId"].(float64)
			response = handleGetAllBoards(uint(userID))
		case common.BoardsGetByID:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				id, _ := data["id"].(float64)
				response = handleGetBoardByID(uint(id))
			}
		case common.BoardsUpdate:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleUpdateBoard(data)
			}
		case common.BoardsDelete:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				id, _ := data["id"].(float64)
				response = handleDeleteBoard(uint(id))
			}
		case common.BoardsGetByTeam:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				teamID, _ := data["teamId"].(float64)
				response = handleGetBoardsByTeam(uint(teamID))
			}
		case common.ColumnsCreate:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleCreateColumn(data)
			}
		case common.ColumnsGetByBoard:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				boardID, _ := data["boardId"].(float64)
				response = handleGetColumnsByBoard(uint(boardID))
			}
		case common.ColumnsUpdate:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleUpdateColumn(data)
			}
		case common.ColumnsDelete:
			var data map[string]any
			if err := json.Unmarshal(d.Body, &data); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				id, _ := data["id"].(float64)
				response = handleDeleteColumn(uint(id))
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

func handleCreateBoard(data map[string]any) common.RPCResponse {
	title, _ := data["title"].(string)
	teamID, _ := data["teamId"].(float64)

	if title == "" {
		return common.RPCResponse{Success: false, Error: "Board title required", StatusCode: 400}
	}

	board := models.Board{
		Title:  title,
		TeamID: uint(teamID),
	}
	if err := common.DB.Create(&board).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to create board", StatusCode: 500}
	}

	common.DB.Preload("Team").First(&board, board.ID)

	// ? Publish event
	common.PublishEvent(common.BoardUpdated, map[string]any{
		"teamId": teamID,
		"board":  board,
	})

	return common.RPCResponse{
		Success: true,
		Data:    board,
	}
}

func handleGetAllBoards(userID uint) common.RPCResponse {
	// ? Get user's teams
	var members []models.TeamMember
	common.DB.Where("user_id = ?", userID).Find(&members)

	teamIDs := make([]uint, len(members))
	for i, m := range members {
		teamIDs[i] = m.TeamID
	}

	var boards []models.Board
	common.DB.Where("team_id IN ?", teamIDs).Preload("Team").Preload("Columns").Find(&boards)

	return common.RPCResponse{
		Success: true,
		Data:    boards,
	}
}

func handleGetBoardByID(id uint) common.RPCResponse {
	var board models.Board
	if err := common.DB.Preload("Team").Preload("Columns").First(&board, id).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Board not found", StatusCode: 404}
	}

	return common.RPCResponse{
		Success: true,
		Data:    board,
	}
}

func handleUpdateBoard(data map[string]any) common.RPCResponse {
	id, _ := data["id"].(float64)
	title, _ := data["title"].(string)

	var board models.Board
	if err := common.DB.First(&board, uint(id)).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Board not found", StatusCode: 404}
	}

	if title != "" {
		board.Title = title
	}

	if err := common.DB.Save(&board).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to update board", StatusCode: 500}
	}

	common.DB.Preload("Team").First(&board, board.ID)

	// ? Publish event
	common.PublishEvent(common.BoardUpdated, map[string]any{
		"teamId": board.TeamID,
		"board":  board,
	})

	return common.RPCResponse{
		Success: true,
		Data:    board,
	}
}

func handleDeleteBoard(id uint) common.RPCResponse {
	var board models.Board
	if err := common.DB.First(&board, id).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Board not found", StatusCode: 404}
	}

	teamID := board.TeamID
	common.DB.Delete(&board)

	// ? Publish event
	common.PublishEvent(common.BoardUpdated, map[string]any{
		"teamId": teamID,
		"board":  nil,
	})

	return common.RPCResponse{Success: true}
}

func handleGetBoardsByTeam(teamID uint) common.RPCResponse {
	var boards []models.Board
	common.DB.Where("team_id = ?", teamID).Preload("Columns").Find(&boards)

	return common.RPCResponse{
		Success: true,
		Data:    boards,
	}
}

func handleCreateColumn(data map[string]any) common.RPCResponse {
	title, _ := data["title"].(string)
	boardID, _ := data["boardId"].(float64)
	order, _ := data["order"].(float64)

	if title == "" {
		return common.RPCResponse{Success: false, Error: "Column title required", StatusCode: 400}
	}

	column := models.Column{
		Title:   title,
		BoardID: uint(boardID),
		Order:   int(order),
	}
	if err := common.DB.Create(&column).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to create column", StatusCode: 500}
	}

	common.DB.Preload("Board").First(&column, column.ID)

	return common.RPCResponse{
		Success: true,
		Data:    column,
	}
}

func handleGetColumnsByBoard(boardID uint) common.RPCResponse {
	var columns []models.Column
	common.DB.Where("board_id = ?", boardID).Order("\"order\" ASC").Find(&columns)

	return common.RPCResponse{
		Success: true,
		Data:    columns,
	}
}

func handleUpdateColumn(data map[string]any) common.RPCResponse {
	id, _ := data["id"].(float64)
	title, _ := data["title"].(string)
	order, ok := data["order"].(float64)

	var column models.Column
	if err := common.DB.First(&column, uint(id)).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Column not found", StatusCode: 404}
	}

	if title != "" {
		column.Title = title
	}
	if ok {
		column.Order = int(order)
	}

	if err := common.DB.Save(&column).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to update column", StatusCode: 500}
	}

	return common.RPCResponse{
		Success: true,
		Data:    column,
	}
}

func handleDeleteColumn(id uint) common.RPCResponse {
	var column models.Column
	if err := common.DB.First(&column, id).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Column not found", StatusCode: 404}
	}

	common.DB.Delete(&column)

	return common.RPCResponse{Success: true}
}

