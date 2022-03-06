-- +goose Up
-- +goose StatementBegin
CREATE TABLE mnemonic_wallets
(
    id serial PRIMARY KEY,
    wallet_uuid uuid NOT NULL,
    hash varchar NOT NULL,
    purpose varchar NOT NULL,
    encrypted_data varchar NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mnemonic_wallets
-- +goose StatementEnd
