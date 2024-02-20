# Notification System

This Notification System is designed to send notifications via various channels, initially supporting email and sms notifications. It is built with scalability in mind, utilizing Docker containers and Docker Compose for easy deployment and scaling.

## Architecture Overview

The system comprises several components:

- **Notification API**: A RESTful API service for queuing notifications.
- **Notification Worker**: A background worker that processes queued notifications and sends them out.
- **RabbitMQ**: Message broker for queueing notification requests.
- **PostgreSQL**: Database for storing user information.
- **Common**: Shared library used by both the API and Worker for common data structures and utilities.

## Prerequisites

Before running this project, ensure you have Docker and Docker Compose installed on your machine. This project is developed and tested with:
- Docker v24.0.6
- Docker Compose v2.23.0

## How to Run

1. **Clone the repository**:
```bash
git clone https://github.com/yourusername/notification-system.git
```

2. **Build Docker images**:

Navigate to the root directory of the project and build the Docker images:
```bash
cd notification-system
docker build -f notification-api/Dockerfile -t notification-api:latest .
docker build -f notification-worker/Dockerfile -t notification-worker:latest .
```

3. **Before starting the services, you may want to populate the database with test users to fully test the notification system's capabilities.
To do this, edit the `scripts/init-users.sql` script. Here's an example of how you can add test users:**
```sql
INSERT INTO users (id, email, phone_number, opted_in) VALUES
('uuid', 'john.doe@example.com', '1234567890', TRUE),
('uuid', 'jane.doe@example.com', '0987654321', TRUE);
```

4. **Setup the `.env` with `MAILCHIMP_API_KEY`, `TWILIO_ACC_SID` and `TWILIO_AUTH_TOKEN` and start the services**:
```bash
docker-compose up
```
This command starts all the services defined in `docker-compose.yml`. You can access the Notification API at `http://localhost:8080`.

5. **Scaling the services**:

To scale the services, use the following command:
```
docker-compose up --scale service_name=number_of_instances
```

API Usage
To queue a notification, send a POST request to http://localhost:8080/v1/notification with the following JSON payload:

```json
{
"type": "email",
"to": ["userID1", "userID2"],
"from": "noreply@yourdomain.com",
"subject": "Test Notification",
"content": "This is a test notification."
}
```
or
```json
{
"type": "sms",
"to": ["userID1", "userID2"],
"from": "+16304896680",
"subject": "Test Notification",
"content": "This is a test notification."
}
```
**_NOTE:_**  The email sending functionality is currently restricted to domains registered with MailChimp due to the use of a free trial account.
Similarly, Twilio, SMS notifications can only be sent to verified phone numbers. This is a limitation of the MailChimp/Twilio services for trial accounts.

## Health Checks

- **Notification API**: Access the health check endpoint at `http://localhost:8080/health`.


## Cleanup

To stop and remove all the running containers associated with this project, run:
```
docker-compose down
```
