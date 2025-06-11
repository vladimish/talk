-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN selected_model TEXT NOT NULL DEFAULT 'gpt-4o-mini';
ALTER TABLE users ADD COLUMN current_conversation TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN current_conversation;
ALTER TABLE users DROP COLUMN selected_model;
-- +goose StatementEnd
