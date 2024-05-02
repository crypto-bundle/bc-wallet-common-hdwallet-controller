-- +goose Up
-- +goose StatementBegin
ALTER TABLE sign_requests
    DROP COLUMN IF EXISTS "derivation_path";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT ('no chance for rollback data') as "message";
-- +goose StatementEnd
