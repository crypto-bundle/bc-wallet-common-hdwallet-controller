-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX wallet_uuid__idx ON mnemonic_wallets (
    "wallet_uuid"
);

CREATE INDEX is_hot__idx ON mnemonic_wallets (
    "is_hot"
);

CREATE UNIQUE INDEX hash__idx ON mnemonic_wallets ("hash");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX wallet_uuid__idx;
DROP INDEX is_hot__idx;
-- +goose StatementEnd
