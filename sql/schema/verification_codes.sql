CREATE TABLE IF NOT EXISTS verification_codes(
    user_id INT NOT NULL,
    verfication_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(users_id) ON DELETE CASCADE
);