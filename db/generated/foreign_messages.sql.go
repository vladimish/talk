// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: foreign_messages.sql

package generated

import (
	"context"
)

const createForeignMessage = `-- name: CreateForeignMessage :one
INSERT INTO foreign_messages (message_id, foreign_message_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING id, message_id, foreign_message_id, created_at
`

type CreateForeignMessageParams struct {
	MessageID        int32
	ForeignMessageID int32
}

func (q *Queries) CreateForeignMessage(ctx context.Context, arg CreateForeignMessageParams) (ForeignMessage, error) {
	row := q.db.QueryRowContext(ctx, createForeignMessage, arg.MessageID, arg.ForeignMessageID)
	var i ForeignMessage
	err := row.Scan(
		&i.ID,
		&i.MessageID,
		&i.ForeignMessageID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteForeignMessage = `-- name: DeleteForeignMessage :exec
DELETE FROM foreign_messages
WHERE message_id = $1
`

func (q *Queries) DeleteForeignMessage(ctx context.Context, messageID int32) error {
	_, err := q.db.ExecContext(ctx, deleteForeignMessage, messageID)
	return err
}

const getForeignMessageByForeignID = `-- name: GetForeignMessageByForeignID :one
SELECT id, message_id, foreign_message_id, created_at FROM foreign_messages
WHERE foreign_message_id = $1
`

func (q *Queries) GetForeignMessageByForeignID(ctx context.Context, foreignMessageID int32) (ForeignMessage, error) {
	row := q.db.QueryRowContext(ctx, getForeignMessageByForeignID, foreignMessageID)
	var i ForeignMessage
	err := row.Scan(
		&i.ID,
		&i.MessageID,
		&i.ForeignMessageID,
		&i.CreatedAt,
	)
	return i, err
}

const getForeignMessageByMessageID = `-- name: GetForeignMessageByMessageID :one
SELECT id, message_id, foreign_message_id, created_at FROM foreign_messages
WHERE message_id = $1
`

func (q *Queries) GetForeignMessageByMessageID(ctx context.Context, messageID int32) (ForeignMessage, error) {
	row := q.db.QueryRowContext(ctx, getForeignMessageByMessageID, messageID)
	var i ForeignMessage
	err := row.Scan(
		&i.ID,
		&i.MessageID,
		&i.ForeignMessageID,
		&i.CreatedAt,
	)
	return i, err
}
