-- +goose Up
-- +goose StatementBegin
ALTER TABLE sign_requests ADD COLUMN account_data bytea DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sign_requests DROP COLUMN account_data;
-- +goose StatementEnd
