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
)






















































































































































































































































    ""


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
