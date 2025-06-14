-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, payment_id, subscription_type, valid_from, valid_to, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, payment_id, subscription_type, valid_from, valid_to, status, created_at, updated_at;

-- name: GetActiveSubscriptionByUserID :one
SELECT id, user_id, payment_id, subscription_type, valid_from, valid_to, status, created_at, updated_at
FROM subscriptions
WHERE user_id = $1 AND status = 'active' AND valid_to > NOW()
ORDER BY valid_to DESC
LIMIT 1;

-- name: GetUserSubscriptions :many
SELECT id, user_id, payment_id, subscription_type, valid_from, valid_to, status, created_at, updated_at
FROM subscriptions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateSubscriptionStatus :one
UPDATE subscriptions
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, payment_id, subscription_type, valid_from, valid_to, status, created_at, updated_at;

-- name: ExpireOldSubscriptions :exec
UPDATE subscriptions
SET status = 'expired', updated_at = NOW()
WHERE status = 'active' AND valid_to <= NOW();