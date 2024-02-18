package notifications

import (
	"encoding/json"
	"os"

	"github.com/pdragnev/notification-system/common"
	"github.com/pdragnev/notification-system/notification-api/internal/queue"
)

type NotificationService struct {
	QueueClient *queue.RabbitMQClient
}

func NewNotificationService() *NotificationService {
	queueClient, _ := queue.NewRabbitMQClient()
	return &NotificationService{
		QueueClient: queueClient,
	}
}

func (s *NotificationService) SendNotification(notification common.Notification) error {
	notificationMessage := common.NotificationMessage{
		Notification: notification,
		RetryCount:   0,
	}
	notificationMessageBytes, err := json.Marshal(notificationMessage)
	if err != nil {
		return err
	}

	return s.QueueClient.PublishMessage(os.Getenv("RABBITMQ_NOTIFICATION_QUEUE_NAME"), notificationMessageBytes)
}
