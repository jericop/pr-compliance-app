// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: approval.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createApproval = `-- name: CreateApproval :one
INSERT INTO approval(uuid, pr_id, sha, approved_on)
VALUES 
  ($1, $2, $3, $4)
RETURNING id, uuid, pr_id, sha, approved_on
`

type CreateApprovalParams struct {
	Uuid       string           `json:"uuid"`
	PrID       pgtype.Int4      `json:"pr_id"`
	Sha        string           `json:"sha"`
	ApprovedOn pgtype.Timestamp `json:"approved_on"`
}

func (q *Queries) CreateApproval(ctx context.Context, arg CreateApprovalParams) (Approval, error) {
	row := q.db.QueryRow(ctx, createApproval,
		arg.Uuid,
		arg.PrID,
		arg.Sha,
		arg.ApprovedOn,
	)
	var i Approval
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.PrID,
		&i.Sha,
		&i.ApprovedOn,
	)
	return i, err
}

const deleteApproval = `-- name: DeleteApproval :exec
DELETE FROM approval
WHERE id = $1
`

func (q *Queries) DeleteApproval(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteApproval, id)
	return err
}

const getApprovalById = `-- name: GetApprovalById :one
SELECT id, uuid, pr_id, sha, approved_on FROM approval
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetApprovalById(ctx context.Context, id int32) (Approval, error) {
	row := q.db.QueryRow(ctx, getApprovalById, id)
	var i Approval
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.PrID,
		&i.Sha,
		&i.ApprovedOn,
	)
	return i, err
}

const getApprovalByUuid = `-- name: GetApprovalByUuid :one
SELECT id, uuid, pr_id, sha, approved_on FROM approval
WHERE uuid = $1 LIMIT 1
`

func (q *Queries) GetApprovalByUuid(ctx context.Context, uuid string) (Approval, error) {
	row := q.db.QueryRow(ctx, getApprovalByUuid, uuid)
	var i Approval
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.PrID,
		&i.Sha,
		&i.ApprovedOn,
	)
	return i, err
}

const getApprovals = `-- name: GetApprovals :many
SELECT id, uuid, pr_id, sha, approved_on FROM approval
`

func (q *Queries) GetApprovals(ctx context.Context) ([]Approval, error) {
	rows, err := q.db.Query(ctx, getApprovals)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Approval{}
	for rows.Next() {
		var i Approval
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.PrID,
			&i.Sha,
			&i.ApprovedOn,
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
