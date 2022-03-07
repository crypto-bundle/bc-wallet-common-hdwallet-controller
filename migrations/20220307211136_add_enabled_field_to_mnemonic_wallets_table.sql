-- +goose Up
-- +goose StatementBegin
ALTER TABLE mnemonic_wallets
    ADD COLUMN is_enabled BOOLEAN NOT NULL DEFAULT true;

CREATE INDEX is_enabled__idx ON mnemonic_wallets (
    "is_enabled"
);

CREATE INDEX is_hot__is_enabled__idx ON mnemonic_wallets ("is_hot", "is_enabled");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX is_enabled__idx;
DROP INDEX is_hot__is_enabled__idx;
ALTER TABLE mnemonic_wallets DROP COLUMN is_enabled;
-- +goose StatementEnd
