-- name: SaveVerificationCode :one
INSERT INTO verification_codes(user_id,verfication_type,expires_at)
VALUES ($1,$2,$3)
RETURNING user_id,verfication_type, created_at, expires_at;