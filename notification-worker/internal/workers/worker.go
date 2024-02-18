package workers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pdragnev/notification-system/common"
	"github.com/pdragnev/notification-system/notification-worker/internal/db"
	"github.com/pdragnev/notification-system/notification-worker/internal/models"
	"github.com/pdragnev/notification-system/notification-worker/internal/notifications"
	"github.com/pdragnev/notification-system/notification-worker/internal/queue"
	"github.com/rabbitmq/amqp091-go"
)

type NotificationWorker struct {
	QueueClient *queue.RabbitMQClient
	UserRepo    db.UserRepository
}

func NewNotificationWorker(queueClient *queue.RabbitMQClient, repo db.UserRepository) *NotificationWorker {
	return &NotificationWorker{
		QueueClient: queueClient,
		UserRepo:    repo,
	}
}

func (worker *NotificationWorker) ProcessMessage(message []byte) error {
	var notificationMsg common.NotificationMessage
	err := json.Unmarshal(message, &notificationMsg)
	if err != nil {
		strErr := fmt.Sprintf("Error deserializing message: %v", err)
		log.Print(strErr)
		return models.NewDeserializingMsgError(strErr)
	}

	// Check if retry count has exceeded max retries
	if notificationMsg.RetryCount >= 5 {
		log.Printf("Max retries exceeded for message: %v", notificationMsg)
		return nil
	}

	notification := notificationMsg.Notification

	processor, err := notifications.GetProcessorForType(string(notification.Type), worker.UserRepo)
	if err != nil {
		strErr := fmt.Sprintf("Error getting processor for type %s: %v", notification.Type, err)
		log.Print(strErr)
		return models.NewProcessingTypeError(strErr)
	}

	// Process the notification
	err = processor.Process(notificationMsg)
	if err != nil {
		log.Printf("Error processing notification: %v", err)
		notificationMsg.RetryCount++
		return models.NewRetryError("Retry due to temporary condition", notificationMsg)
	}

	return nil
}

func (worker *NotificationWorker) Start() {
	handler := func(d amqp091.Delivery) error {
		return worker.ProcessMessage(d.Body)
	}

	worker.QueueClient.StartConsuming(os.Getenv("RABBITMQ_NOTIFICATION_QUEUE_NAME"), handler)
}
