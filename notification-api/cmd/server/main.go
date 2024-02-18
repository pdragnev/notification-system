package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pdragnev/notification-system/common"
	"github.com/pdragnev/notification-system/notification-api/internal/notifications"
)

func notificationHandler(notificationService common.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var notification common.Notification
		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if !common.IsValidType(notification.Type) {
			http.Error(w, "Invalid notification type", http.StatusBadRequest)
			return
		}

		if err := notificationService.SendNotification(notification); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Notification enqueued successfully"))
	}
}

func main() {
	notificationService := notifications.NewNotificationService()

	http.HandleFunc("/v1/notification", notificationHandler(notificationService))

	log.Println("API Gateway listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
