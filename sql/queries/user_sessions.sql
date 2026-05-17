-- name: CreateUserSession :one
INSERT INTO user_sessions(user_id,expires_at)
VALUES ($1,$2)
RETURNING id,user_id,created_at, expires_at;

-- name: DeleteUserSession :exec
DELETE FROM user_sessions
WHERE id = $1;