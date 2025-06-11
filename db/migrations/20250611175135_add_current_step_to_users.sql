-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN current_step TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN current_step;
-- +goose StatementEnd
