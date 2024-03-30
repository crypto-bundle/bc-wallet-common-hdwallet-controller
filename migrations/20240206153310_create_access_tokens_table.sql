-- +goose Up
-- +goose StatementBegin
CREATE TABLE access_tokens
(
    id serial PRIMARY KEY,
    uuid uuid NOT NULL,

    wallet_uuid uuid NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_tokens
-- +goose StatementEnd
