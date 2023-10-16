// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: pull_request_event.sql

package postgres

import (
	"context"
	"time"
)

const createPullRequestEvent = `-- name: CreatePullRequestEvent :one
INSERT INTO pull_request_event(pr_id, action, sha, is_merged, last_updated)
VALUES 
  ($1, $2, $3, $4, $5)
RETURNING id, pr_id, action, sha, is_merged, last_updated
`

type CreatePullRequestEventParams struct {
	PrID        int32     `json:"pr_id"`
	Action      string    `json:"action"`
	Sha         string    `json:"sha"`
	IsMerged    bool      `json:"is_merged"`
	LastUpdated time.Time `json:"last_updated"`
}

func (q *Queries) CreatePullRequestEvent(ctx context.Context, arg CreatePullRequestEventParams) (PullRequestEvent, error) {
	row := q.db.QueryRow(ctx, createPullRequestEvent,
		arg.PrID,
		arg.Action,
		arg.Sha,
		arg.IsMerged,
		arg.LastUpdated,
	)
	var i PullRequestEvent
	err := row.Scan(
		&i.ID,
		&i.PrID,
		&i.Action,
		&i.Sha,
		&i.IsMerged,
		&i.LastUpdated,
	)
	return i, err
}

const deletePullRequestEvent = `-- name: DeletePullRequestEvent :exec
DELETE FROM pull_request_event
WHERE id = $1
`

func (q *Queries) DeletePullRequestEvent(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deletePullRequestEvent, id)
	return err
}

const getPullRequestEvent = `-- name: GetPullRequestEvent :one
SELECT id, pr_id, action, sha, is_merged, last_updated FROM pull_request_event
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPullRequestEvent(ctx context.Context, id int32) (PullRequestEvent, error) {
	row := q.db.QueryRow(ctx, getPullRequestEvent, id)
	var i PullRequestEvent
	err := row.Scan(
		&i.ID,
		&i.PrID,
		&i.Action,
		&i.Sha,
		&i.IsMerged,
		&i.LastUpdated,
	)
	return i, err
}

const getPullRequestEvents = `-- name: GetPullRequestEvents :many
SELECT id, pr_id, action, sha, is_merged, last_updated FROM pull_request_event
`

func (q *Queries) GetPullRequestEvents(ctx context.Context) ([]PullRequestEvent, error) {
	rows, err := q.db.Query(ctx, getPullRequestEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PullRequestEvent{}
	for rows.Next() {
		var i PullRequestEvent
		if err := rows.Scan(
			&i.ID,
			&i.PrID,
			&i.Action,
			&i.Sha,
			&i.IsMerged,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
