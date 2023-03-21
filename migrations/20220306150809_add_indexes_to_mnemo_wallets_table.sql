-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX uuid__idx ON mnemonic_wallets (
     "uuid"
);
CREATE UNIQUE INDEX wallet_uuid__idx ON mnemonic_wallets (
    "wallet_uuid"
);

CREATE UNIQUE INDEX wallet_uuid__mnemonic_wallet_uuid__idx ON mnemonic_wallets (
    "wallet_uuid", "uuid"
);

CREATE INDEX is_hot__idx ON mnemonic_wallets (
    "is_hot"
);

CREATE INDEX wallet_uuid__is_hot__idx ON mnemonic_wallets (
   "wallet_uuid", "is_hot"
);

CREATE UNIQUE INDEX hash__idx ON mnemonic_wallets ("hash");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uuid__idx;
DROP INDEX wallet_uuid__idx;
DROP INDEX wallet_uuid__mnemonic_wallet_uuid__idx;
DROP INDEX is_hot__idx;
DROP INDEX wallet_uuid__is_hot__idx;
DROP INDEX hash__idx;
-- +goose StatementEnd
