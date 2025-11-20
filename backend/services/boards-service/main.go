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
			var req CreateBoardRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleCreateBoard(req)
			}
		case common.BoardsGetAll:
			var req GetAllBoardsRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetAllBoards(req)
			}
		case common.BoardsGetByID:
			var req GetBoardByIDRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetBoardByID(req)
			}
		case common.BoardsUpdate:
			var req UpdateBoardRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleUpdateBoard(req)
			}
		case common.BoardsDelete:
			var req DeleteBoardRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleDeleteBoard(req)
			}
		case common.BoardsGetByTeam:
			var req GetBoardsByTeamRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetBoardsByTeam(req)
			}
		case common.ColumnsCreate:
			var req CreateColumnRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleCreateColumn(req)
			}
		case common.ColumnsGetByBoard:
			var req GetColumnsByBoardRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleGetColumnsByBoard(req)
			}
		case common.ColumnsUpdate:
			var req UpdateColumnRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleUpdateColumn(req)
			}
		case common.ColumnsDelete:
			var req DeleteColumnRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				response = common.RPCResponse{Success: false, Error: "Invalid payload", StatusCode: 400}
			} else {
				response = handleDeleteColumn(req)
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

func handleCreateBoard(req CreateBoardRequest) common.RPCResponse {
	if req.Title == "" {
		return common.RPCResponse{Success: false, Error: "Board title required", StatusCode: 400}
	}
	// ? Validate team exists
	var team models.Team
	if req.UserID == 0 {
		return common.RPCResponse{Success: false, Error: "User ID required", StatusCode: 400}
	}

	rpcResp, err := common.CallRPC(common.TeamsGetByID, map[string]interface{}{"id": req.TeamID})
	if err != nil || rpcResp == nil || !rpcResp.Success {
		return common.RPCResponse{Success: false, Error: "Team not found", StatusCode: 404}
	}

	if dataBytes, marshalErr := json.Marshal(rpcResp.Data); marshalErr == nil {
		_ = json.Unmarshal(dataBytes, &team)
	}

	isMember := false
	for _, m := range team.Members {
		if m.UserID == req.UserID {
			isMember = true
			break
		}
	}
	if !isMember {
		return common.RPCResponse{Success: false, Error: "User is not a team member", StatusCode: 403}
	}

	localTeam := models.Team{}
	if err := common.DB.Where(&models.Team{ID: team.ID}).FirstOrCreate(&localTeam, models.Team{ID: team.ID, Name: team.Name}).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to synchronize team locally", StatusCode: 500}
	}

	board := models.Board{Title: req.Title, TeamID: req.TeamID}
	if err := common.DB.Create(&board).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to create board", StatusCode: 500}
	}
	if rpcResp, rpcErr := common.CallRPC(common.TeamsGetByID, map[string]interface{}{"id": board.TeamID}); rpcErr == nil && rpcResp != nil && rpcResp.Success {
		var team models.Team
		if b, err := json.Marshal(rpcResp.Data); err == nil {
			_ = json.Unmarshal(b, &team)
			board.Team = team
		}
	}

	common.PublishEvent(common.BoardUpdated, map[string]interface{}{
		"teamId": req.TeamID,
		"board":  board,
	})

	return toRPC(boardResponse{Success: true, Data: &board})
}

func handleGetAllBoards(req GetAllBoardsRequest) common.RPCResponse {
	rpcResp, rpcErr := common.CallRPC(common.TeamsGetUserTeams, map[string]interface{}{"userId": req.UserID})
	if rpcErr != nil {
		return common.RPCResponse{Success: false, Error: "Failed to fetch user teams", StatusCode: 500}
	}
	if rpcResp == nil || !rpcResp.Success {
		return toRPC(boardsListResponse{Success: true, Data: []models.Board{}})
	}

	var rawTeams []map[string]interface{}
	if b, err := json.Marshal(rpcResp.Data); err == nil {
		_ = json.Unmarshal(b, &rawTeams)
	}

	teamIDs := make([]uint, 0, len(rawTeams))
	for _, team := range rawTeams {
		if idv, ok := team["id"]; ok {
			switch v := idv.(type) {
			case float64:
				teamIDs = append(teamIDs, uint(v))
			case int:
				teamIDs = append(teamIDs, uint(v))
			case uint:
				teamIDs = append(teamIDs, v)
			}
		}
	}

	if len(teamIDs) == 0 {
		return toRPC(boardsListResponse{Success: true, Data: []models.Board{}})
	}

	var boards []models.Board
	common.DB.Where("team_id IN ?", teamIDs).
		Preload("Columns").
		Find(&boards)

	// ? Attach team details to each board
	teamCache := make(map[uint]models.Team)
	for i := range boards {
		tid := boards[i].TeamID
		if team, ok := teamCache[tid]; ok {
			boards[i].Team = team
			continue
		}
		if rpcResp, rpcErr := common.CallRPC(common.TeamsGetByID, map[string]interface{}{"id": tid}); rpcErr == nil && rpcResp != nil && rpcResp.Success {
			var team models.Team
			if b, err := json.Marshal(rpcResp.Data); err == nil {
				_ = json.Unmarshal(b, &team)
				boards[i].Team = team
				teamCache[tid] = team
			}
		}
	}

	return toRPC(boardsListResponse{Success: true, Data: boards})
}

func handleGetBoardByID(req GetBoardByIDRequest) common.RPCResponse {
	var board models.Board
	if err := common.DB.Preload("Columns").First(&board, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Board not found", StatusCode: 404}
	}
	if rpcResp, rpcErr := common.CallRPC(common.TeamsGetByID, map[string]interface{}{"id": board.TeamID}); rpcErr == nil && rpcResp != nil && rpcResp.Success {
		var team models.Team
		if b, err := json.Marshal(rpcResp.Data); err == nil {
			_ = json.Unmarshal(b, &team)
			board.Team = team
		}
	}
	return toRPC(boardResponse{Success: true, Data: &board})
}

func handleUpdateBoard(req UpdateBoardRequest) common.RPCResponse {
	var board models.Board
	if err := common.DB.First(&board, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Board not found", StatusCode: 404}
	}
	if req.Title != nil {
		board.Title = *req.Title
	}
	if err := common.DB.Save(&board).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to update board", StatusCode: 500}
	}
	if rpcResp, rpcErr := common.CallRPC(common.TeamsGetByID, map[string]interface{}{"id": board.TeamID}); rpcErr == nil && rpcResp != nil && rpcResp.Success {
		var team models.Team
		if b, err := json.Marshal(rpcResp.Data); err == nil {
			_ = json.Unmarshal(b, &team)
			board.Team = team
		}
	}

	common.PublishEvent(common.BoardUpdated, map[string]interface{}{
		"teamId": board.TeamID,
		"board":  board,
	})

	return toRPC(boardResponse{Success: true, Data: &board})
}

func handleDeleteBoard(req DeleteBoardRequest) common.RPCResponse {
	var board models.Board
	if err := common.DB.First(&board, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Board not found", StatusCode: 404}
	}
	teamID := board.TeamID
	common.DB.Delete(&board)

	common.PublishEvent(common.BoardUpdated, map[string]interface{}{
		"teamId": teamID,
		"board":  nil,
	})

	return toRPC(successResponse{Success: true})
}

func handleGetBoardsByTeam(req GetBoardsByTeamRequest) common.RPCResponse {
	var boards []models.Board
	common.DB.Where("team_id = ?", req.TeamID).Preload("Columns").Find(&boards)
	return toRPC(boardsListResponse{Success: true, Data: boards})
}

func handleCreateColumn(req CreateColumnRequest) common.RPCResponse {
	if req.Title == "" {
		return common.RPCResponse{Success: false, Error: "Column title required", StatusCode: 400}
	}
	column := models.Column{Title: req.Title, BoardID: req.BoardID, Order: req.Order}
	if err := common.DB.Create(&column).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to create column", StatusCode: 500}
	}
	common.DB.Preload("Board").First(&column, column.ID)
	return toRPC(columnResponse{Success: true, Data: &column})
}

func handleGetColumnsByBoard(req GetColumnsByBoardRequest) common.RPCResponse {
	var columns []models.Column
	common.DB.Where("board_id = ?", req.BoardID).Order("\"order\" ASC").Find(&columns)
	return toRPC(columnsListResponse{Success: true, Data: columns})
}

func handleUpdateColumn(req UpdateColumnRequest) common.RPCResponse {
	var column models.Column
	if err := common.DB.First(&column, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Column not found", StatusCode: 404}
	}
	if req.Title != nil {
		column.Title = *req.Title
	}
	if req.Order != nil {
		column.Order = *req.Order
	}
	if err := common.DB.Save(&column).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Failed to update column", StatusCode: 500}
	}
	return toRPC(columnResponse{Success: true, Data: &column})
}

func handleDeleteColumn(req DeleteColumnRequest) common.RPCResponse {
	var column models.Column
	if err := common.DB.First(&column, req.ID).Error; err != nil {
		return common.RPCResponse{Success: false, Error: "Column not found", StatusCode: 404}
	}
	common.DB.Delete(&column)
	return toRPC(successResponse{Success: true})
}
