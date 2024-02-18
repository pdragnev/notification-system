package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pdragnev/notification-system/common"
	"github.com/pdragnev/notification-system/notification-worker/internal/db"
)

type EmailProcessor struct {
	BaseProcessor
}

func NewEmailProcessor(userRepo db.UserRepository) *EmailProcessor {
	return &EmailProcessor{
		BaseProcessor: BaseProcessor{UserRepo: userRepo},
	}
}

func (p *EmailProcessor) Process(notificationMsg common.NotificationMessage) error {
	notification := notificationMsg.Notification
	userEmails, err := p.UserRepo.GetUserEmailsByIds(notification.To)
	if err != nil {
		return fmt.Errorf("failed to fetch user emails: %v", err)
	}

	messagePayload := map[string]interface{}{
		"key": os.Getenv("MAILCHIMP_API_KEY"),
		"message": map[string]interface{}{
			"from_email": notification.From,
			"subject":    notification.Subject,
			"text":       notification.Content,
			"to":         formatRecipients(userEmails),
		},
	}

	payloadBytes, err := json.Marshal(messagePayload)
	if err != nil {
		return fmt.Errorf("error marshalling message payload: %v", err)
	}

	// Send the email
	resp, err := http.Post("https://mandrillapp.com/api/1.0/messages/send", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("email sending failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func formatRecipients(emails []string) []map[string]string {
	recipients := make([]map[string]string, len(emails))
	for i, email := range emails {
		recipients[i] = map[string]string{"email": email, "type": "to"}
	}
	return recipients
}
