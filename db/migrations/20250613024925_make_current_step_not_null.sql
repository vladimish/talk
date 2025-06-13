-- +goose Up
-- +goose StatementBegin
-- Set default value for existing NULL current_step entries
UPDATE users SET current_step = 'menu' WHERE current_step IS NULL;

-- Make current_step NOT NULL with default value
ALTER TABLE users 
ALTER COLUMN current_step SET NOT NULL,
ALTER COLUMN current_step SET DEFAULT 'menu';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove NOT NULL constraint and default value
ALTER TABLE users 
ALTER COLUMN current_step DROP NOT NULL,
ALTER COLUMN current_step DROP DEFAULT;
-- +goose StatementEnd
