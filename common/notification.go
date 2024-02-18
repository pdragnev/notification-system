package common

type NotificationService interface {
	SendNotification(notification Notification) error
}

type NotificationType string

const (
	EmailNotificationType NotificationType = "email"
	// more types
)

func IsValidType(t NotificationType) bool {
	switch t {
	case EmailNotificationType:
		return true
	default:
		return false
	}
}

type Notification struct {
	Type    NotificationType `json:"type"`
	To      []string         `json:"to"`
	From    string           `json:"from"`
	Subject string           `json:"subject"`
	Content string           `json:"content"`
}

type NotificationMessage struct {
	Notification Notification `json:"notification"`
	RetryCount   int          `json:"retryCount"`
}
