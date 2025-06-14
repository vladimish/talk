-- name: CreateTransaction :one
INSERT INTO transactions (user_id, token_type, amount, transaction_type, model_used, description)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserTokenBalance :one
SELECT 
    COALESCE(SUM(CASE WHEN token_type = 'premium' THEN amount ELSE 0 END), 0) as premium_balance,
    COALESCE(SUM(CASE WHEN token_type = 'regular' THEN amount ELSE 0 END), 0) as regular_balance
FROM transactions 
WHERE user_id = $1;

-- name: GetUserTransactionHistory :many
SELECT * FROM transactions 
WHERE user_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserTokenBalanceByType :one
SELECT COALESCE(SUM(amount), 0) as balance
FROM transactions 
WHERE user_id = $1 AND token_type = $2;