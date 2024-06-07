-- +goose Up
-- +goose StatementBegin
CREATE TABLE pow_proofs
(
    id   serial PRIMARY KEY,
    uuid uuid      NOT NULL,

    access_token_uuid uuid NOT NULL,

    message_check_nonce bigint NOT NULL check (message_check_nonce >= 0),
    message_hash char(64) NOT NULL,
    message_data bytea NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS pow_proofs__uuid_idx ON pow_proofs (
    "uuid"
);
CREATE INDEX IF NOT EXISTS pow_proofs__access_token_uuid__idx ON pow_proofs (
    "access_token_uuid"
);
CREATE INDEX IF NOT EXISTS pow_proofs__message_hash__idx ON pow_proofs (
    "message_hash"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pow_proofs
-- +goose StatementEnd
