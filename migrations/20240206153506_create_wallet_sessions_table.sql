-- +goose Up
-- +goose StatementBegin
CREATE TABLE mnemonic_wallet_sessions
(
    id serial PRIMARY KEY,
    access_token_id integer NOT NULL check (access_token_id >= 1),
    mnemonic_wallet_id uuid NOT NULL,
    is_closed BOOLEAN NOT NULL DEFAULT true,

    expired_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + (15 * interval '1 seconds'),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mnemonic_wallet_sessions
-- +goose StatementEnd
