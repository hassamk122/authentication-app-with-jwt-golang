-- name: CreateUserSession :one
INSERT INTO user_sessions(user_id,expires_at)
VALUES ($1,$2)
RETURNING id,user_id,created_at, expires_at;

-- name: DeleteUserSession :exec
DELETE FROM user_sessions
WHERE id = $1;

-- name: FindById :one
SELECT id, user_id, created_at, expires_at
FROM user_sessions
WHERE id = $1;

-- name: UpdateSession :one
UPDATE user_sessions
SET expires_at = $1
WHERE id = $2
RETURNING *;