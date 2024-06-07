-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS wsat__token_uuid__idx ON wallet_sessions_access_tokens (
    "token_uuid"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS wsat__token_uuid__idx;
-- +goose StatementEnd
