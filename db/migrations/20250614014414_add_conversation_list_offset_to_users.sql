-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN conversation_list_offset INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN conversation_list_offset;
-- +goose StatementEnd
