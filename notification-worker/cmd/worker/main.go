package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pdragnev/notification-system/notification-worker/internal/db"
	"github.com/pdragnev/notification-system/notification-worker/internal/queue"
	"github.com/pdragnev/notification-system/notification-worker/internal/workers"
)

func main() {
	pool, err := db.Connect(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer pool.Close()

	userRepository := db.NewUserRepository(pool)

	//Connection to RabbitMQ
	rabbitMQConfig := queue.RabbitMQConfig{
		URL:               os.Getenv("RABBITMQ_URL"),
		NotificationQueue: os.Getenv("RABBITMQ_NOTIFICATION_QUEUE_NAME"),
	}
	rabbitMQClient, err := queue.NewRabbitMQClient(rabbitMQConfig)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ client: %v", err)
	}
	defer func() {
		if err := rabbitMQClient.Connection.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

	notificationWorker := workers.NewNotificationWorker(rabbitMQClient, userRepository)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Worker started. Press Ctrl+C to stop.")
		notificationWorker.Start()
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received, initiating graceful shutdown...")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Worker shutdown gracefully")
}
