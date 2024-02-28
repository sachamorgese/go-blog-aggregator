-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(255) NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL
);


-- +goose Down
DROP TABLE users;