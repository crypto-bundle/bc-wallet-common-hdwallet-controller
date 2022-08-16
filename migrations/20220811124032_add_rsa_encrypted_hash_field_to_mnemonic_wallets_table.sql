-- +goose Up
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets ADD COLUMN rsa_encrypted_hash varchar NOT NULL DEFAULT '';
ALTER TABLE mnemonic_wallets ALTER COLUMN rsa_encrypted_hash DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets DROP COLUMN rsa_encrypted_hash;
-- +goose StatementEnd
