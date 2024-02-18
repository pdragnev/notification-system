package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pdragnev/notification-system/notification-worker/internal/db"
	"github.com/pdragnev/notification-system/notification-worker/internal/queue"
	"github.com/pdragnev/notification-system/notification-worker/internal/workers"
)

func main() {
	//Connection to DB
	pool, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepository := db.NewUserRepository(pool)

	//Connection to RabbitMQ
	rabbitMQClient := queue.NewRabbitMQClient()
	defer rabbitMQClient.Connection.Close()

	notificationWorker := workers.NewNotificationWorker(rabbitMQClient, userRepository)

	// Setup channel to listen for interrupt or terminate signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Start your worker
	go func() {
		log.Println("Worker started. Press Ctrl+C to stop.")
		notificationWorker.Start()
	}()

	// Wait for interrupt or terminate signal
	<-stopChan
	log.Println("Shutdown signal received, exiting...")

	// Attempt a graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()

	log.Println("Worker shutdown gracefully")
}
