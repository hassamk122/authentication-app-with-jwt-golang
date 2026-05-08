-- name: CreateUserSession :one
INSERT INTO user_sessions(user_id)
VALUES ($1)
RETURNING id,user_id,created_at, expires_at;