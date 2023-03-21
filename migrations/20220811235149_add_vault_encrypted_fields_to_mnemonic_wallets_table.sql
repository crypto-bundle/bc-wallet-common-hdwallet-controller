-- +goose Up
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets ADD COLUMN vault_encrypted bytea NOT NULL;
ALTER TABLE mnemonic_wallets ADD COLUMN vault_encrypted_hash varchar NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets DROP COLUMN vault_encrypted;
ALTER TABLE mnemonic_wallets DROP COLUMN vault_encrypted_hash;
-- +goose StatementEnd
