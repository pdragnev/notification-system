package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("APP_ENV") == "development" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}
}

type UserRepository interface {
	GetUserEmailsByIds(userIds []string) ([]string, error)
}

func Connect() (*pgxpool.Pool, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	return pgxpool.Connect(context.Background(), databaseUrl)
}
