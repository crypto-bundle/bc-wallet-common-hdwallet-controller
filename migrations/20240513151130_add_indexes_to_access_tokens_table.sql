-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS access_tokens__uuid__wallet_uuid__idx ON access_tokens (
    "uuid", "wallet_uuid"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS access_tokens__uuid__wallet_uuid__idx;
-- +goose StatementEnd
