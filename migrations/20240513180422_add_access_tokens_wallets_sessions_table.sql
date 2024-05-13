-- +goose Up
-- +goose StatementBegin
CREATE TABLE access_tokens_wallet_sessions
(
    serial_number bigint NOT NULL,

    token_uuid uuid NOT NULL,
    wallet_session_uuid uuid NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS atws__serial_number__token_uuid__idx ON access_tokens_wallet_sessions (
    "serial_number", "token_uuid"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_tokens_wallet_sessions
-- +goose StatementEnd
