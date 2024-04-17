-- +goose Up
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets
    DROP COLUMN IF EXISTS "wallet_uuid";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT ('no chance for rollback data') as "message";
-- +goose StatementEnd
