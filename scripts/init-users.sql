CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(20),
    opted_in BOOLEAN DEFAULT TRUE
);

INSERT INTO users (id, email, phone_number, opted_in) VALUES
('80fc203f-3856-43a5-b2d3-b604a640ec54', 'petar@vasilkotsev.com', '+359892091234', TRUE),
('563cfe60-6ed7-49ac-ba33-f05758831980', 'testing@vasilkotsev.com', '+359890123456', TRUE);

