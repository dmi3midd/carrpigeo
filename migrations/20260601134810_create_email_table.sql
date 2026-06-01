-- +goose Up
CREATE TABLE emails (
    id VARCHAR(20) PRIMARY KEY,
    sender VARCHAR(255) NOT NULL,
    reciever VARCHAR(255) NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
);

-- +goose Down
SELECT 'down SQL query';
