-- +goose Up
-- +goose StatementBegin
ALTER TABLE payments 
ADD COLUMN invoice_payload TEXT,
ADD COLUMN message_id TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE payments 
DROP COLUMN IF EXISTS invoice_payload,
DROP COLUMN IF EXISTS message_id;
-- +goose StatementEnd