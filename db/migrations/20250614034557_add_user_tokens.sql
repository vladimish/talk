-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_type VARCHAR(20) NOT NULL CHECK (token_type IN ('premium', 'regular')),
    amount INTEGER NOT NULL, -- positive for credits, negative for debits
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('initial_credit', 'message_cost', 'admin_credit', 'admin_debit')),
    model_used VARCHAR(100), -- model name for message costs
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_user_token_type ON transactions(user_id, token_type);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);

-- Give all existing users initial tokens
INSERT INTO transactions (user_id, token_type, amount, transaction_type, description)
SELECT id, 'premium', 1000, 'initial_credit', 'Initial premium token grant'
FROM users;

INSERT INTO transactions (user_id, token_type, amount, transaction_type, description)
SELECT id, 'regular', 100, 'initial_credit', 'Initial regular token grant'
FROM users;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd
