-- name: CreateForeignMessage :one
INSERT INTO foreign_messages (message_id, foreign_message_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetForeignMessageByMessageID :one
SELECT * FROM foreign_messages
WHERE message_id = $1;

-- name: GetForeignMessageByForeignID :one
SELECT * FROM foreign_messages
WHERE foreign_message_id = $1;

-- name: DeleteForeignMessage :exec
DELETE FROM foreign_messages
WHERE message_id = $1;