-- +goose Up
-- +goose StatementBegin
ALTER TABLE payments ALTER COLUMN amount TYPE BIGINT;
ALTER TABLE transactions ALTER COLUMN amount TYPE BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE payments ALTER COLUMN amount TYPE INTEGER;
ALTER TABLE transactions ALTER COLUMN amount TYPE INTEGER;
-- +goose StatementEnd
