-- +goose Up
-- +goose StatementBegin
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    invoice_link TEXT NOT NULL,
    telegram_payment_charge_id TEXT UNIQUE,
    provider_payment_charge_id TEXT,
    currency VARCHAR(10) NOT NULL DEFAULT 'XTR',
    amount INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'pre_checkout', 'paid', 'failed', 'cancelled')),
    subscription_type VARCHAR(50) NOT NULL DEFAULT 'monthly',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_telegram_payment_charge_id ON payments(telegram_payment_charge_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
-- +goose StatementEnd
