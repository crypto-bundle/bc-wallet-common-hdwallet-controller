-- +goose Up
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets RENAME COLUMN hash TO mnemonic_hash;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets RENAME COLUMN mnemonic_hash TO hash;
-- +goose StatementEnd
