-- +goose Up
-- +goose StatementBegin
-- Drop the old conversation column and add conversation_id
ALTER TABLE messages DROP COLUMN conversation;
ALTER TABLE messages ADD COLUMN conversation_id BIGINT REFERENCES conversations(id) ON DELETE SET NULL;

CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE messages DROP COLUMN conversation_id;
ALTER TABLE messages ADD COLUMN conversation UUID;
-- +goose StatementEnd
