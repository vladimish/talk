-- name: CreateMessage :one
INSERT INTO messages (
    message_type,
    user_id,
    sent_by,
    conversation
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetMessagesByUserID :many
SELECT * FROM messages
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetMessagesByConversation :many
SELECT * FROM messages
WHERE conversation = $1
ORDER BY created_at ASC;

-- name: GetConversationsByUserID :many
SELECT DISTINCT conversation, MIN(created_at) as first_message_at
FROM messages
WHERE user_id = $1 AND conversation IS NOT NULL
GROUP BY conversation
ORDER BY first_message_at DESC;