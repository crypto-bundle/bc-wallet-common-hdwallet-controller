-- +goose Up
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets ADD COLUMN unload_interval bigint NOT NULL DEFAULT '15000000000';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets DROP COLUMN unload_interval;
-- +goose StatementEnd
