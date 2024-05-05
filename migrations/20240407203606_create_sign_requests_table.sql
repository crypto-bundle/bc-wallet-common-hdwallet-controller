-- +goose Up
-- +goose StatementBegin
CREATE TABLE sign_requests
(
    id   serial PRIMARY KEY,
    uuid uuid      NOT NULL,

    mnemonic_wallet_uuid uuid NOT NULL,
    session_uuid uuid NOT NULL,
    purpose_uuid uuid NOT NULL,

    derivation_path integer[],

    status smallint NOT NULL check (status >= 1),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS mnemonic_wallet_uuid__expired_at ON sign_requests (
     "mnemonic_wallet_uuid"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mnemonic_wallet_sessionsж
-- +goose StatementEnd
