-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets
(
    id serial PRIMARY KEY,
    uuid uuid NOT NULL,
    title varchar NOT NULL,
    purpose varchar NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets
-- +goose StatementEnd