package models

type MailchimpEmailResponse struct {
	Email        string `json:"email"`
	Status       string `json:"status"`
	RejectReason string `json:"reject_reason"`
	ID           string `json:"_id"`
}
