-- +goose Up
-- +goose StatementBegin
ALTER TABLE access_tokens ADD COLUMN role smallint DEFAULT 99 check (role >= 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE access_tokens DROP COLUMN role;
-- +goose StatementEnd
