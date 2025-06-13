-- +goose Up
-- +goose StatementBegin
-- Update users table to use conversation_id instead of text
ALTER TABLE users DROP COLUMN current_conversation;
ALTER TABLE users ADD COLUMN current_conversation BIGINT REFERENCES conversations(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN current_conversation;
ALTER TABLE users ADD COLUMN current_conversation TEXT;
-- +goose StatementEnd
