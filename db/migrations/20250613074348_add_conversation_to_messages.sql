-- +goose Up
-- +goose StatementBegin
ALTER TABLE messages ADD COLUMN conversation UUID;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE messages DROP COLUMN conversation;
-- +goose StatementEnd
