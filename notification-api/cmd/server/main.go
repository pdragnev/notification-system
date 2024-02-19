package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pdragnev/notification-system/common"
	"github.com/pdragnev/notification-system/notification-api/internal/notifications"
	"github.com/pdragnev/notification-system/notification-api/internal/queue"
)

func notificationHandler(notificationService common.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var notification common.Notification
		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			log.Printf("Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if !common.IsValidType(notification.Type) {
			http.Error(w, "Invalid notification type", http.StatusBadRequest)
			return
		}

		if err := notificationService.SendNotification(notification); err != nil {
			log.Printf("Error sending notification: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Notification enqueued successfully"))
	}
}

func main() {
	notificationQueueName := os.Getenv("RABBITMQ_NOTIFICATION_QUEUE_NAME")
	if notificationQueueName == "" {
		log.Fatal("RABBITMQ_NOTIFICATION_QUEUE_NAME must be set")
	}
	config := queue.RabbitMQConfig{
		URL:               os.Getenv("RABBITMQ_URL"),
		DLXExchange:       os.Getenv("DLX_EXCHANGE_NAME"),
		DLXQueue:          os.Getenv("DLX_QUEUE_NAME"),
		NotificationQueue: os.Getenv("RABBITMQ_NOTIFICATION_QUEUE_NAME"),
	}
	notificationService, err := notifications.NewNotificationService(notificationQueueName, config)
	if err := notificationService.QueueClient.SetupQueues(); err != nil {
		log.Fatalf("Failed to setup queues: %v", err)
		return
	}

	if err != nil {
		log.Fatalf("Failed to initialize notification service: %v", err)
	}

	http.HandleFunc("/v1/notification", notificationHandler(notificationService))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server gracefully stopped")
}
