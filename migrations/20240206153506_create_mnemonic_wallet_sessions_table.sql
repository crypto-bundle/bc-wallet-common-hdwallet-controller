-- +goose Up
-- +goose StatementBegin
CREATE TABLE mnemonic_wallet_sessions
(
    id                   serial PRIMARY KEY,
    uuid                 uuid      NOT NULL,

    access_token_uuid    uuid      NOT NULL,
    mnemonic_wallet_uuid uuid      NOT NULL,

    status smallint NOT NULL check (status >= 1),

    -- started_at default value = CURRENT_TIMESTAMP + 2s delay
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + (2 * interval '1 seconds'),
    -- expired_at default = started_at + 20 sec (CURRENT_TIMESTAMP + 2s + 20s)
    expired_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + (22 * interval '1 seconds'),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE UNIQUE INDEX mnemonic_wallet_uuid__expired_at ON mnemonic_wallet_sessions (
     "mnemonic_wallet_uuid", "expired_at"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mnemonic_wallet_sessions
-- +goose StatementEnd
