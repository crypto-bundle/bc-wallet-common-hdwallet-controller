-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallet_sessions_access_tokens_counters (
  id SERIAL PRIMARY KEY,
  token_uuid uuid NOT NULL,
  counter_value BIGINT NOT NULL check (counter_value >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS wallet_sessions_access_tokens_counters__token_uuid__idx
    ON wallet_sessions_access_tokens_counters (token_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallet_sessions_access_tokens_counters;
-- +goose StatementEnd