-- +goose Up
-- +goose StatementBegin
ALTER TABLE wallets
    ADD COLUMN is_enabled BOOLEAN NOT NULL DEFAULT true;

CREATE INDEX is_enabled__idx ON wallets (
    "is_enabled"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX is_enabled__idx;
ALTER TABLE wallets DROP COLUMN is_enabled;
-- +goose StatementEnd
