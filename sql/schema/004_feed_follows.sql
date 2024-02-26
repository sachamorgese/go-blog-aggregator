-- +goose Up
CREATE TABLE users_feeds(user_id UUID REFERENCES users(id),
    feed_id UUID REFERENCES feeds(id),
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT users_feeds_unique UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE users_feeds;