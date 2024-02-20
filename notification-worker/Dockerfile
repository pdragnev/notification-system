# Start from the official Golang base image
FROM golang:1.21.4 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY common /app/common
COPY notification-worker /app/notification-worker

WORKDIR /app/notification-worker

RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main ./cmd/worker/main.go

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /main .

# Command to run the executable
CMD ["./main"]
