-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN web_search_enabled BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN web_search_enabled;
-- +goose StatementEnd
