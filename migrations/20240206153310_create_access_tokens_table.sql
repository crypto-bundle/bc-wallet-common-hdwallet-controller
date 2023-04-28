-- +goose Up
-- +goose StatementBegin
CREATE TABLE access_tokens
(
    id serial PRIMARY KEY,
    uuid uuid NOT NULL,

    mnemonic_wallet_uuid uuid NOT NULL,
    token varchar NOT NULL,


    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS token_access_permissions
-- +goose StatementEnd
