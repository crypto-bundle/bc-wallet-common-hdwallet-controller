-- +goose Up
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets RENAME COLUMN encrypted_data TO rsa_encrypted;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets RENAME COLUMN rsa_encrypted TO encrypted_data;
-- +goose StatementEnd
