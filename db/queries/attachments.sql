-- name: CreateAttachment :one
INSERT INTO attachments (message_id, s3_name, content_type, size)
VALUES ($1, $2, $3, $4)
RETURNING *;