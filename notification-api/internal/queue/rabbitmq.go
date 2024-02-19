package queue

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL               string
	DLXExchange       string
	DLXQueue          string
	NotificationQueue string
}

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

func NewRabbitMQClient(config RabbitMQConfig) (*RabbitMQClient, error) {
	conn, err := amqp091.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	return &RabbitMQClient{Connection: conn}, nil
}

func (client *RabbitMQClient) PublishMessage(queueName string, message []byte) error {
	ch, err := client.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		"",        // Exchange
		queueName, // Routing key (queue name)
		false,     // Mandatory
		false,     // Immediate
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "text/plain",
			Body:         message,
		},
	)
	return err
}

func (client *RabbitMQClient) SetupQueues() error {
	ch, err := client.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	dlxName := os.Getenv("DLX_EXCHANGE_NAME")
	dlqName := os.Getenv("DLX_QUEUE_NAME")
	primaryQueueName := os.Getenv("RABBITMQ_NOTIFICATION_QUEUE_NAME")

	// Ensure DLX exists
	err = ch.ExchangeDeclare(
		dlxName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLX: %v", err)
	}

	_, err = ch.QueueDeclare(
		dlqName,
		true,
		false,
		false,
		false,
		amqp091.Table{"x-dead-letter-exchange": dlxName},
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %v", err)
	}

	// Bind DLQ to DLX
	err = ch.QueueBind(
		dlqName,
		"", // routing key
		dlxName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind DLQ to DLX: %v", err)
	}

	// Create or ensure primary queue exists with DLX configuration
	_, err = ch.QueueDeclare(
		primaryQueueName,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		amqp091.Table{
			"x-dead-letter-exchange": dlxName, // DLX
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare primary queue with DLX: %v", err)
	}

	return nil
}
