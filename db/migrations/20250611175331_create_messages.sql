-- +goose Up
-- +goose StatementBegin
CREATE TYPE message_sender AS ENUM ('user', 'bot');

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    message_type JSONB NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sent_by message_sender NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_user_id ON messages(user_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
DROP TYPE message_sender;
-- +goose StatementEnd
