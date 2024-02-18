CREATE DATABASE notification_service;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(20),
    opted_in BOOLEAN DEFAULT TRUE
);

INSERT INTO users (email, phone_number, opted_in) VALUES
('vasil@vasilkotsev.com', '0987654321', TRUE),
('petar@vasilkotsev.com', '0987654321', TRUE),
('testing@vasilkotsev.com', '0987654321', TRUE)
