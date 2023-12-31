// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: pull_request_action.sql

package postgres

import (
	"context"
)

const createPullRequestAction = `-- name: CreatePullRequestAction :one
INSERT INTO pull_request_action(name)
VALUES 
  ($1)
RETURNING name
`

func (q *Queries) CreatePullRequestAction(ctx context.Context, name string) (string, error) {
	row := q.db.QueryRow(ctx, createPullRequestAction, name)
	err := row.Scan(&name)
	return name, err
}

const deletePullRequestAction = `-- name: DeletePullRequestAction :exec
DELETE FROM pull_request_action
WHERE name = $1
`

func (q *Queries) DeletePullRequestAction(ctx context.Context, name string) error {
	_, err := q.db.Exec(ctx, deletePullRequestAction, name)
	return err
}

const getPullRequestAction = `-- name: GetPullRequestAction :one
SELECT name FROM pull_request_action
WHERE name = $1 LIMIT 1
`

func (q *Queries) GetPullRequestAction(ctx context.Context, name string) (string, error) {
	row := q.db.QueryRow(ctx, getPullRequestAction, name)
	err := row.Scan(&name)
	return name, err
}

const getPullRequestActions = `-- name: GetPullRequestActions :many
SELECT name FROM pull_request_action
`

func (q *Queries) GetPullRequestActions(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getPullRequestActions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
