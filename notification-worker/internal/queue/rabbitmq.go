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

type RabbitMQConfig struct {
	URL               string
	NotificationQueue string
}

type RabbitMQClient struct {
	Connection *amqp091.Connection
	config     *RabbitMQConfig
}

func NewRabbitMQClient(config RabbitMQConfig) (*RabbitMQClient, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("RabbitMQ URL must not be empty")
	}
	if config.NotificationQueue == "" {
		return nil, fmt.Errorf("NotificationQueue name must not be empty")
	}

	conn, err := amqp091.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	return &RabbitMQClient{
		Connection: conn,
		config:     &config,
	}, nil
}

func (client *RabbitMQClient) StartConsuming(handler func(amqp091.Delivery) error) {
	ch, err := client.Connection.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		client.config.NotificationQueue,
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
				client.handleProcessingError(err, d)
			} else {
				d.Ack(false)
			}
		}(d)
	}
}

func (client *RabbitMQClient) handleProcessingError(err error, d amqp091.Delivery) {
	switch e := err.(type) {
	case *models.RetryError:
		updatedMessageBytes, _ := json.Marshal(e.UpdatedMessage)
		if requeueErr := client.requeueMessage(updatedMessageBytes); requeueErr != nil {
			log.Printf("Failed to requeue message: %v", requeueErr)
		}
		d.Ack(false)
	case *models.DeserializingMsgError, *models.ProcessingTypeError, *models.MaxRetryError:
		d.Nack(false, false) // Send to DLQ
	default:
		d.Nack(false, true) // Requeue for temporary issues
	}
}

func (client *RabbitMQClient) requeueMessage(updatedMessage []byte) error {
	ch, err := client.Connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		"",                              // exchange
		client.config.NotificationQueue, // routing key (queue name)
		false,                           // mandatory
		false,                           // immediate
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
