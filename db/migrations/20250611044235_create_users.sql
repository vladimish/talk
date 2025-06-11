-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    foreign_id BIGINT UNIQUE NOT NULL,
    language   TEXT          NOT NULL,
    created_at TIMESTAMP     NOT NULL,
    updated_at TIMESTAMP     NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users
-- +goose StatementEnd
