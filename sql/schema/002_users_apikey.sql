-- +goose Up
ALTER TABLE users
ADD COLUMN api_key TIMESTAMP;


-- +goose Down
ALTER TABLE users
DROP COLUMN api_key;
