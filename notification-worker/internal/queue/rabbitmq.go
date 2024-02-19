package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/pdragnev/notification-system/notification-worker/internal/models"
	"github.com/rabbitmq/amqp091-go"
)

func init() {
	if os.Getenv("APP_ENV") == "development" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}
}

type RabbitMQClient struct {
	Connection *amqp091.Connection
}

func NewRabbitMQClient() *RabbitMQClient {
	conn, err := amqp091.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	return &RabbitMQClient{
		Connection: conn,
	}
}

func (client *RabbitMQClient) StartConsuming(queueName string, handler func(amqp091.Delivery) error) {
	ch, err := client.Connection.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName,
		"",
		false, // we manually ack/nack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	maxWorkersStr := os.Getenv("MAX_WORKERS")
	var maxWorkers int
	if maxWorkersStr != "" {
		maxWorkers, err = strconv.Atoi(maxWorkersStr)
		if err != nil {
			log.Fatalf("Invalid MAX_WORKERS value: %v", err)
			return
		}
	} else {
		maxWorkers = runtime.NumCPU() * 2
	}
	sem := make(chan struct{}, maxWorkers)

	for d := range msgs {
		sem <- struct{}{}
		go func(d amqp091.Delivery) {
			defer func() { <-sem }()
			if err := handler(d); err != nil {
				client.handleProcessingError(err, d, queueName)
			} else {
				d.Ack(false)
			}
		}(d)
	}
}

func (client *RabbitMQClient) handleProcessingError(err error, d amqp091.Delivery, queueName string) {
	switch e := err.(type) {
	case *models.RetryError:
		updatedMessageBytes, _ := json.Marshal(e.UpdatedMessage)
		if requeueErr := client.requeueMessage(queueName, updatedMessageBytes); requeueErr != nil {
			log.Printf("Failed to requeue message: %v", requeueErr)
		}
		d.Ack(false)
	case *models.DeserializingMsgError, *models.ProcessingTypeError, *models.MaxRetryError:
		d.Nack(false, false) // Send to DLQ
	default:
		d.Nack(false, true) // Requeue for temporary issues
	}
}

func (client *RabbitMQClient) requeueMessage(queueName string, updatedMessage []byte) error {
	ch, err := client.Connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "text/plain",
			Body:         updatedMessage,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}
