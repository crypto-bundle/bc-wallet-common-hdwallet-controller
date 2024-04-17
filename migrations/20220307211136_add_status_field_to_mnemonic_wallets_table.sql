-- +goose Up
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets
    ADD COLUMN status smallint NOT NULL check (status >= 1);
CREATE INDEX status__idx ON mnemonic_wallets (
    "status"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX status__idx;
ALTER TABLE mnemonic_wallets DROP COLUMN status;
-- +goose StatementEnd
