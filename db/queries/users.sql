-- name: GetUserByForeignID :one
SELECT *
FROM users
WHERE foreign_id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (foreign_id, language, current_step, selected_model, conversation_list_offset, web_search_enabled, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateUserCurrentStep :exec
UPDATE users
SET current_step = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserSelectedModel :exec
UPDATE users
SET selected_model = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserLanguage :exec
UPDATE users
SET language = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserConversationListOffset :exec
UPDATE users
SET conversation_list_offset = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserCurrentConversationID :exec
UPDATE users
SET current_conversation = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserWebSearchEnabled :exec
UPDATE users
SET web_search_enabled = $2, updated_at = NOW()
WHERE id = $1;