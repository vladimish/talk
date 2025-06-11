-- name: GetUserByForeignID :one
SELECT *
FROM users
WHERE foreign_id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (foreign_id, language, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUserCurrentStep :exec
UPDATE users
SET current_step = $2, updated_at = NOW()
WHERE id = $1;