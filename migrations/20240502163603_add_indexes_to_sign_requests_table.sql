-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS uuid__status__idx ON sign_requests (
    "uuid", "status"
);
CREATE INDEX IF NOT EXISTS session_uuid__idx ON sign_requests (
    "session_uuid"
);
DROP INDEX IF EXISTS mnemonic_wallet_uuid__expired_at;
CREATE INDEX IF NOT EXISTS mnemonic_wallet_uuid__idx ON sign_requests (
    "mnemonic_wallet_uuid"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uuid__status__idx;
DROP INDEX IF EXISTS session_uuid__status__idx;
CREATE UNIQUE INDEX IF NOT EXISTS mnemonic_wallet_uuid__expired_at ON sign_requests (
    "mnemonic_wallet_uuid"
);
DROP INDEX IF EXISTS mnemonic_wallet_uuid__idx;
-- +goose StatementEnd
