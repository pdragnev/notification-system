module github.com/pdragnev/notification-system/notification-api

go 1.21.4

replace github.com/pdragnev/notification-system/common => ../common

require (
	github.com/joho/godotenv v1.5.1
	github.com/pdragnev/notification-system/common v0.0.0-00010101000000-000000000000
	github.com/rabbitmq/amqp091-go v1.9.0
)
