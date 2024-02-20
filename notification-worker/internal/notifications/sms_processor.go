package notifications

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pdragnev/notification-system/common"
	"github.com/pdragnev/notification-system/notification-worker/internal/db"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type SmsProcessor struct {
	BaseProcessor
	client *twilio.RestClient
}

func NewSmsProcessor(userRepo db.UserRepository) *SmsProcessor {
	param := twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACC_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	}
	client := twilio.NewRestClientWithParams(param)
	return &SmsProcessor{
		BaseProcessor: BaseProcessor{UserRepo: userRepo},
		client:        client,
	}
}

func (p *SmsProcessor) Process(notificationMsg common.NotificationMessage) error {
	notification := notificationMsg.Notification
	userPhoneNumbers, err := p.UserRepo.GetUserPhonesByIds(context.Background(), notification.To)
	if err != nil || len(userPhoneNumbers) == 0 {
		return fmt.Errorf("failed to fetch user phone numbers: %v", err)
	}

	params := &api.CreateMessageParams{}
	params.SetBody(notification.Content)
	params.SetFrom(notification.From)

	for i := 0; i < len(userPhoneNumbers); i++ {
		params.SetTo(userPhoneNumbers[i])
		resp, err := p.client.Api.CreateMessage(params)
		if err != nil {
			return err
		} else {
			if resp.Sid != nil {
				log.Print(*resp.Sid)
			} else {
				log.Println(resp.Sid)
			}
		}
	}

	return nil
}
