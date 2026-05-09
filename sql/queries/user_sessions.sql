-- name: CreateUserSession :one
INSERT INTO user_sessions(user_id,expires_at)
VALUES ($1,$2)
RETURNING id,user_id,created_at, expires_at;