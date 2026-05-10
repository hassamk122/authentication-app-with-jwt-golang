-- name: CreateUser :one
INSERT INTO users(username, email, password,verified)
VALUES ($1, $2, $3,$4)
RETURNING public_id, username, email,verified, created_at, updated_at;


-- name: GetUserByEmail :one
SELECT public_id, username, email,password, created_at, updated_at
FROM users
WHERE email = $1;