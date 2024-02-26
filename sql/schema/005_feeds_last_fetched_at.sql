-- +goose Up
ALTER TABLE feeds
    ADD COLUMN last_fetched_at VARCHAR(64) UNIQUE NOT NULL DEFAULT encode(sha256(random()::text::bytea), 'hex');


-- +goose Down
ALTER TABLE users
    DROP COLUMN api_key;
