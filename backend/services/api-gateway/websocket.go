package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/wetask/backend/pkg/common"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	rooms  map[string]bool
	userID uint
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	// Subscribe to RabbitMQ events
	msgs, err := common.RabbitMQChannel.Consume(
		"events_queue", // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Println("Failed to consume events:", err)
		return
	}

	go func() {
		for msg := range msgs {
			var eventData map[string]interface{}
			if err := json.Unmarshal(msg.Body, &eventData); err != nil {
				continue
			}

			// Broadcast to appropriate rooms
			eventType := msg.RoutingKey
			h.broadcastToRoom(eventType, eventData)
		}
	}()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) broadcastToRoom(eventType string, data map[string]interface{}) {
	var room string
	switch eventType {
	case common.TaskCreated, common.TaskUpdated, common.TaskDeleted:
		if boardID, ok := data["boardId"].(float64); ok {
			room = "board:" + fmt.Sprintf("%.0f", boardID)
		}
	case common.BoardUpdated:
		if teamID, ok := data["teamId"].(float64); ok {
			room = "team:" + fmt.Sprintf("%.0f", teamID)
		}
	case common.TeamMemberAdded, common.TeamMemberRemoved:
		if teamID, ok := data["teamId"].(float64); ok {
			room = "team:" + fmt.Sprintf("%.0f", teamID)
		}
	}

	if room == "" {
		return
	}

	event := map[string]interface{}{
		"type": eventType,
		"data": data,
	}
	message, _ := json.Marshal(event)

	for client := range h.clients {
		if client.rooms[room] {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

func handleWebSocket(c *gin.Context, hub *Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		hub:   hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		rooms:  make(map[string]bool),
		userID: 0,
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		eventType, _ := msg["type"].(string)
		data, _ := msg["data"].(map[string]interface{})

		switch eventType {
		case "join:board":
			boardID, _ := data["boardId"].(float64)
			room := "board:" + fmt.Sprintf("%.0f", boardID)
			c.rooms[room] = true
		case "leave:board":
			boardID, _ := data["boardId"].(float64)
			room := "board:" + fmt.Sprintf("%.0f", boardID)
			delete(c.rooms, room)
		case "join:team":
			teamID, _ := data["teamId"].(float64)
			userID, _ := data["userId"].(float64)
			room := "team:" + fmt.Sprintf("%.0f", teamID)
			c.rooms[room] = true
			c.userID = uint(userID)
		}
	}
}

func (c *Client) writePump() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

