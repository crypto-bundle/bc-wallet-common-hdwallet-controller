-- +goose Up
-- +goose StatementBegin
ALTER TABLE wallets ADD COLUMN strategy smallint NOT NULL check (strategy >= 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE wallets DROP COLUMN strategy;
-- +goose StatementEnd
