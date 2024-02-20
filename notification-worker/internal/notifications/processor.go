package notifications

import (
	"fmt"

	"github.com/pdragnev/notification-system/common"
	"github.com/pdragnev/notification-system/notification-worker/internal/db"
)

type Processor interface {
	Process(notificationMsg common.NotificationMessage) error
}

type BaseProcessor struct {
	UserRepo db.UserRepository
}

func GetProcessorForType(notificationType string, userRepo db.UserRepository) (Processor, error) {
	switch notificationType {
	case "email":
		return NewEmailProcessor(userRepo), nil
	case "sms":
		return NewSmsProcessor(userRepo), nil
	default:
		return nil, fmt.Errorf("unknown notification type: %s", notificationType)
	}
}
