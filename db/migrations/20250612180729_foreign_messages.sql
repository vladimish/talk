-- +goose Up
-- +goose StatementBegin
CREATE TABLE foreign_messages (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    foreign_message_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(message_id),
    UNIQUE(foreign_message_id)
);

CREATE INDEX idx_foreign_messages_message_id ON foreign_messages(message_id);
CREATE INDEX idx_foreign_messages_foreign_message_id ON foreign_messages(foreign_message_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS foreign_messages;
-- +goose StatementEnd
