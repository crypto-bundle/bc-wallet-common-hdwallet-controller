-- +goose Up
-- +goose StatementBegin
CREATE TABLE access_tokens_permissions
(
    access_token_id serial PRIMARY KEY,
    permission_type smallint NOT NULL check (permission_type >= 1),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_tokens_permissions
-- +goose StatementEnd
