-- name: CreateConversation :one
INSERT INTO conversations (name, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetConversationsByUserID :many
SELECT * FROM conversations
WHERE user_id = $1
ORDER BY updated_at DESC;

-- name: GetConversationByID :one
SELECT * FROM conversations
WHERE id = $1
LIMIT 1;

-- name: UpdateConversationName :exec
UPDATE conversations
SET name = $2, updated_at = NOW()
WHERE id = $1;