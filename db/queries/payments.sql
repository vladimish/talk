-- name: CreatePayment :one
INSERT INTO payments (user_id, invoice_link, currency, amount, subscription_type, invoice_payload)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, invoice_link, telegram_payment_charge_id, provider_payment_charge_id, currency, amount, status, subscription_type, invoice_payload, message_id, created_at, updated_at;

-- name: GetPaymentByID :one
SELECT id, user_id, invoice_link, telegram_payment_charge_id, provider_payment_charge_id, currency, amount, status, subscription_type, invoice_payload, message_id, created_at, updated_at
FROM payments 
WHERE id = $1;

-- name: GetPaymentByTelegramChargeID :one
SELECT id, user_id, invoice_link, telegram_payment_charge_id, provider_payment_charge_id, currency, amount, status, subscription_type, invoice_payload, message_id, created_at, updated_at
FROM payments 
WHERE telegram_payment_charge_id = $1;

-- name: GetPaymentByInvoicePayload :one
SELECT id, user_id, invoice_link, telegram_payment_charge_id, provider_payment_charge_id, currency, amount, status, subscription_type, invoice_payload, message_id, created_at, updated_at
FROM payments 
WHERE invoice_payload = $1;

-- name: UpdatePaymentStatus :one
UPDATE payments 
SET status = $2, telegram_payment_charge_id = $3, provider_payment_charge_id = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, invoice_link, telegram_payment_charge_id, provider_payment_charge_id, currency, amount, status, subscription_type, invoice_payload, message_id, created_at, updated_at;

-- name: UpdatePaymentWithInvoice :one
UPDATE payments 
SET invoice_link = $2, invoice_payload = $3, message_id = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, invoice_link, telegram_payment_charge_id, provider_payment_charge_id, currency, amount, status, subscription_type, invoice_payload, message_id, created_at, updated_at;

-- name: GetUserPayments :many
SELECT id, user_id, invoice_link, telegram_payment_charge_id, provider_payment_charge_id, currency, amount, status, subscription_type, invoice_payload, message_id, created_at, updated_at
FROM payments 
WHERE user_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPendingPayments :many
SELECT id, user_id, invoice_link, telegram_payment_charge_id, provider_payment_charge_id, currency, amount, status, subscription_type, invoice_payload, message_id, created_at, updated_at
FROM payments 
WHERE status = 'pending' 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;