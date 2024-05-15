-- +goose Up
-- +goose StatementBegin
ALTER TABLE mnemonic_wallet_sessions
    DROP COLUMN IF EXISTS "access_token_uuid";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mnemonic_wallet_sessions
    ADD COLUMN IF NOT EXISTS access_token_uuid  uuid NOT NULL DEFAULT '00000000-00000000-00000000-00000000';
-- +goose StatementEnd
