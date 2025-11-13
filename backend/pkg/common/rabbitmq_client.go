package common

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var RabbitMQConn *amqp.Connection
var RabbitMQChannel *amqp.Channel

// ? InitRabbitMQ initializes RabbitMQ connection
func InitRabbitMQ() error {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		getEnv("RABBITMQ_USER", "admin"),
		getEnv("RABBITMQ_PASSWORD", "admin123"),
		getEnv("RABBITMQ_HOST", "localhost"),
		getEnv("RABBITMQ_PORT", "5672"),
	)
	var err error
	RabbitMQConn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	RabbitMQChannel, err = RabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	log.Println("RabbitMQ connected successfully")
	return nil
}

// ? CloseRabbitMQ closes RabbitMQ connection
func CloseRabbitMQ() {
	if RabbitMQChannel != nil {
		RabbitMQChannel.Close()
	}
	if RabbitMQConn != nil {
		RabbitMQConn.Close()
	}
}

// ? CallRPC makes an RPC call to a service
func CallRPC(pattern string, payload interface{}) (*RPCResponse, error) {
	if RabbitMQChannel == nil {
		return nil, fmt.Errorf("RabbitMQ channel not initialized")
	}

	// ? Declare queue for response
	replyQueue, err := RabbitMQChannel.QueueDeclare(
		"",    // * name
		false, // * durable
		true,  // * delete when unused
		true,  // * exclusive
		false, // * no-wait
		nil,   // * arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare reply queue: %w", err)
	}

	// ? Consume response with unique consumer tag
	consumerTag := fmt.Sprintf("consumer_%d", time.Now().UnixNano())
	msgs, err := RabbitMQChannel.Consume(
		replyQueue.Name, // * queue
		consumerTag,     // * consumer tag (must be unique per consumer)
		true,            // * auto-ack
		false,           // * exclusive
		false,           // * no-local
		false,           // * no-wait
		nil,             // * args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}
	defer func() {
		// ? Cancell consumer after request done
		if err := RabbitMQChannel.Cancel(consumerTag, false); err != nil {
			log.Printf("Warning: failed to cancel consumer: %v", err)
		}
	}()

	// ? Serialize payload
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// ? Publish request
	corrID := fmt.Sprintf("%d", time.Now().UnixNano())
	err = RabbitMQChannel.Publish(
		"",      // * exchange
		pattern, // * routing key
		false,   // * mandatory
		false,   // * immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       replyQueue.Name,
			Body:          body,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message: %w", err)
	}

	// ? Wait for response with timeout
	timeout := time.After(5 * time.Second)
	for {
		select {
		case msg := <-msgs:
			if msg.CorrelationId == corrID {
				var response RPCResponse
				if err := json.Unmarshal(msg.Body, &response); err != nil {
					return nil, fmt.Errorf("failed to unmarshal response: %w", err)
				}
				return &response, nil
			}
		case <-timeout:
			return nil, fmt.Errorf("RPC call timeout")
		}
	}
}

// ? PublishEvent publishes an event to RabbitMQ
func PublishEvent(event string, data interface{}) error {
	if RabbitMQChannel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	err = RabbitMQChannel.Publish(
		"events", // * exchange
		event,    // * routing key
		false,    // * mandatory
		false,    // * immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}
