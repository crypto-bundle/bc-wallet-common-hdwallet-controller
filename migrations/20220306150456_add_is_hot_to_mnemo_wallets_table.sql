-- +goose Up
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets
    ADD COLUMN is_hot BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets
    DROP COLUMN is_hot;
-- +goose StatementEnd
