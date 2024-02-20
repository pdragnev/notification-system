# Start from the official Golang base image
FROM golang:1.21.4 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY common /app/common
COPY notification-api /app/notification-api

WORKDIR /app/notification-api

RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main ./cmd/server/main.go

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add curl

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /main .

# Command to run the executable
CMD ["./main"]
