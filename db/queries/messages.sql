-- name: CreateMessage :one
INSERT INTO messages (
    message_type,
    user_id,
    sent_by,
    conversation_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetMessagesByUserID :many
SELECT * FROM messages
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetMessagesByConversationID :many
SELECT * FROM messages
WHERE conversation_id = $1
ORDER BY created_at ASC;

-- name: GetLatestMessageByConversationID :one
SELECT * FROM messages
WHERE conversation_id = $1
ORDER BY created_at DESC
LIMIT 1;