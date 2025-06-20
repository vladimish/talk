// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: transactions.sql

package generated

import (
	"context"
	"database/sql"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (user_id, token_type, amount, transaction_type, model_used, description)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, token_type, amount, transaction_type, model_used, description, created_at
`

type CreateTransactionParams struct {
	UserID          int64
	TokenType       string
	Amount          int64
	TransactionType string
	ModelUsed       sql.NullString
	Description     sql.NullString
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.UserID,
		arg.TokenType,
		arg.Amount,
		arg.TransactionType,
		arg.ModelUsed,
		arg.Description,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TokenType,
		&i.Amount,
		&i.TransactionType,
		&i.ModelUsed,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const getUserTokenBalance = `-- name: GetUserTokenBalance :one
SELECT 
    COALESCE(SUM(CASE WHEN token_type = 'premium' THEN amount ELSE 0 END), 0) as premium_balance,
    COALESCE(SUM(CASE WHEN token_type = 'regular' THEN amount ELSE 0 END), 0) as regular_balance
FROM transactions 
WHERE user_id = $1
`

type GetUserTokenBalanceRow struct {
	PremiumBalance interface{}
	RegularBalance interface{}
}

func (q *Queries) GetUserTokenBalance(ctx context.Context, userID int64) (GetUserTokenBalanceRow, error) {
	row := q.db.QueryRowContext(ctx, getUserTokenBalance, userID)
	var i GetUserTokenBalanceRow
	err := row.Scan(&i.PremiumBalance, &i.RegularBalance)
	return i, err
}

const getUserTokenBalanceByType = `-- name: GetUserTokenBalanceByType :one
SELECT COALESCE(SUM(amount), 0) as balance
FROM transactions 
WHERE user_id = $1 AND token_type = $2
`

type GetUserTokenBalanceByTypeParams struct {
	UserID    int64
	TokenType string
}

func (q *Queries) GetUserTokenBalanceByType(ctx context.Context, arg GetUserTokenBalanceByTypeParams) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getUserTokenBalanceByType, arg.UserID, arg.TokenType)
	var balance interface{}
	err := row.Scan(&balance)
	return balance, err
}

const getUserTransactionHistory = `-- name: GetUserTransactionHistory :many
SELECT id, user_id, token_type, amount, transaction_type, model_used, description, created_at FROM transactions 
WHERE user_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type GetUserTransactionHistoryParams struct {
	UserID int64
	Limit  int32
	Offset int32
}

func (q *Queries) GetUserTransactionHistory(ctx context.Context, arg GetUserTransactionHistoryParams) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, getUserTransactionHistory, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.TokenType,
			&i.Amount,
			&i.TransactionType,
			&i.ModelUsed,
			&i.Description,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
