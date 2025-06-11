-- name: CreateMessage :one
INSERT INTO messages (
    message_type,
    user_id,
    sent_by
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetMessagesByUserID :many
SELECT * FROM messages
WHERE user_id = $1
ORDER BY created_at ASC;