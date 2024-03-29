// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: feed_follow.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
INSERT INTO users_feeds(id, user_id, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id, feed_id, id, created_at, updated_at
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.NullUUID `json:"user_id"`
	FeedID    uuid.NullUUID `json:"feed_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (UsersFeed, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.UserID,
		arg.FeedID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i UsersFeed
	err := row.Scan(
		&i.UserID,
		&i.FeedID,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteFeedFollow = `-- name: DeleteFeedFollow :exec
DELETE FROM users_feeds
WHERE user_id = $1 AND feed_id = $2
RETURNING user_id, feed_id, id, created_at, updated_at
`

type DeleteFeedFollowParams struct {
	UserID uuid.NullUUID `json:"user_id"`
	FeedID uuid.NullUUID `json:"feed_id"`
}

func (q *Queries) DeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollow, arg.UserID, arg.FeedID)
	return err
}

const getAllFeedFollowsForUser = `-- name: GetAllFeedFollowsForUser :many
SELECT user_id, feed_id, id, created_at, updated_at FROM users_feeds
WHERE user_id = $1
`

func (q *Queries) GetAllFeedFollowsForUser(ctx context.Context, userID uuid.NullUUID) ([]UsersFeed, error) {
	rows, err := q.db.QueryContext(ctx, getAllFeedFollowsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UsersFeed
	for rows.Next() {
		var i UsersFeed
		if err := rows.Scan(
			&i.UserID,
			&i.FeedID,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
