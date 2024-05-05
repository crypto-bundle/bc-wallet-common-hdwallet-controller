-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS uuid__started_at__expired_at__status__idx ON mnemonic_wallet_sessions (
    "uuid", "started_at", "expired_at", "status"
);
CREATE INDEX IF NOT EXISTS started_at__expired_at__status__idx ON mnemonic_wallet_sessions (
   "started_at", "expired_at", "status"
);
CREATE INDEX IF NOT EXISTS mw_uuid__started_at__expired_at__status__idx ON mnemonic_wallet_sessions (
    "mnemonic_wallet_uuid", "started_at", "expired_at", "status"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uuid__started_at__expired_at__status__idx;
DROP INDEX IF EXISTS started_at__expired_at__status__idx;
DROP INDEX IF EXISTS mw_uuid__started_at__expired_at__status__idx;
-- +goose StatementEnd
