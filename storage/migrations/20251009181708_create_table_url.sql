-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url (
	id INT AUTO_INCREMENT PRIMARY KEY,
	code VARCHAR(30) NOT NULL UNIQUE,
	url VARCHAR(1000) NOT NULL,
    user_id BIGINT DEFAULT NULL,
    status TINYINT DEFAULT 1,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	last_visited_at DATETIME DEFAULT NULL,
	expired_at DATETIME DEFAULT NULL,
	count INT DEFAULT 0,
    INDEX idx_status (status)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url;
-- +goose StatementEnd
