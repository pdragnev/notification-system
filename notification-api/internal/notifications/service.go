package notifications

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pdragnev/notification-system/common"
	"github.com/pdragnev/notification-system/notification-api/internal/queue"
)

type NotificationService struct {
	QueueClient       *queue.RabbitMQClient
	NotificationQueue string
}

func NewNotificationService(notificationQueueName string, config queue.RabbitMQConfig) (*NotificationService, error) {
	queueClient, err := queue.NewRabbitMQClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize RabbitMQ client: %w", err)
	}
	return &NotificationService{
		QueueClient:       queueClient,
		NotificationQueue: notificationQueueName,
	}, nil
}

func (s *NotificationService) SendNotification(notification common.Notification) error {
	notificationMessage := common.NotificationMessage{
		Notification: notification,
		RetryCount:   0,
	}
	notificationMessageBytes, err := json.Marshal(notificationMessage)
	if err != nil {
		log.Printf("Error marshaling notification message: %v", err)
		return err
	}

	if err := s.QueueClient.PublishMessage(s.NotificationQueue, notificationMessageBytes); err != nil {
		log.Printf("Error publishing notification message: %v", err)
		return err
	}

	log.Printf("Notification message enqueued successfully")
	return nil
}
