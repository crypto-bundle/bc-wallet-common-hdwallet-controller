-- +goose Up
-- +goose StatementBegin
ALTER TABLE access_tokens ADD COLUMN hash char(64) DEFAULT null;
UPDATE "access_tokens" ac SET hash=substr(sha256(raw_data)::varchar(66),3) WHERE hash IS NULL;
ALTER TABLE access_tokens ALTER COLUMN hash DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE access_tokens DROP COLUMN hash;
-- +goose StatementEnd
