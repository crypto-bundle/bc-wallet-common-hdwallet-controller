-- +goose Up
-- +goose StatementBegin
ALTER TABLE access_tokens ADD COLUMN IF NOT EXISTS raw_data bytea DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE access_tokens DROP COLUMN IF EXISTS raw_data;
-- +goose StatementEnd
