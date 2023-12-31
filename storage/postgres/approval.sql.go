// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: approval.sql

package postgres

import (
	"context"
)

const createApproval = `-- name: CreateApproval :one
INSERT INTO approval(schema_id, uuid, pr_id, sha, is_approved)
VALUES 
  ($1, $2, $3, $4, $5)
RETURNING id, schema_id, uuid, pr_id, sha, is_approved, last_updated
`

type CreateApprovalParams struct {
	SchemaID   int32  `json:"schema_id"`
	Uuid       string `json:"uuid"`
	PrID       int32  `json:"pr_id"`
	Sha        string `json:"sha"`
	IsApproved bool   `json:"is_approved"`
}

func (q *Queries) CreateApproval(ctx context.Context, arg CreateApprovalParams) (Approval, error) {
	row := q.db.QueryRow(ctx, createApproval,
		arg.SchemaID,
		arg.Uuid,
		arg.PrID,
		arg.Sha,
		arg.IsApproved,
	)
	var i Approval
	err := row.Scan(
		&i.ID,
		&i.SchemaID,
		&i.Uuid,
		&i.PrID,
		&i.Sha,
		&i.IsApproved,
		&i.LastUpdated,
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
SELECT id, schema_id, uuid, pr_id, sha, is_approved, last_updated FROM approval
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetApprovalById(ctx context.Context, id int32) (Approval, error) {
	row := q.db.QueryRow(ctx, getApprovalById, id)
	var i Approval
	err := row.Scan(
		&i.ID,
		&i.SchemaID,
		&i.Uuid,
		&i.PrID,
		&i.Sha,
		&i.IsApproved,
		&i.LastUpdated,
	)
	return i, err
}

const getApprovalByPrIDSha = `-- name: GetApprovalByPrIDSha :one
SELECT id, schema_id, uuid, pr_id, sha, is_approved, last_updated FROM approval
WHERE pr_id = $1 AND sha = $2 LIMIT 1
`

type GetApprovalByPrIDShaParams struct {
	PrID int32  `json:"pr_id"`
	Sha  string `json:"sha"`
}

func (q *Queries) GetApprovalByPrIDSha(ctx context.Context, arg GetApprovalByPrIDShaParams) (Approval, error) {
	row := q.db.QueryRow(ctx, getApprovalByPrIDSha, arg.PrID, arg.Sha)
	var i Approval
	err := row.Scan(
		&i.ID,
		&i.SchemaID,
		&i.Uuid,
		&i.PrID,
		&i.Sha,
		&i.IsApproved,
		&i.LastUpdated,
	)
	return i, err
}

const getApprovalByUuid = `-- name: GetApprovalByUuid :one
SELECT id, schema_id, uuid, pr_id, sha, is_approved, last_updated FROM approval
WHERE uuid = $1 LIMIT 1
`

func (q *Queries) GetApprovalByUuid(ctx context.Context, uuid string) (Approval, error) {
	row := q.db.QueryRow(ctx, getApprovalByUuid, uuid)
	var i Approval
	err := row.Scan(
		&i.ID,
		&i.SchemaID,
		&i.Uuid,
		&i.PrID,
		&i.Sha,
		&i.IsApproved,
		&i.LastUpdated,
	)
	return i, err
}

const getApprovals = `-- name: GetApprovals :many
SELECT id, schema_id, uuid, pr_id, sha, is_approved, last_updated FROM approval
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
			&i.SchemaID,
			&i.Uuid,
			&i.PrID,
			&i.Sha,
			&i.IsApproved,
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

const getCreateStatusInputsFromApprovalUuid = `-- name: GetCreateStatusInputsFromApprovalUuid :one
SELECT p.installation_id, u.login, r.name, a.sha
FROM approval a, pull_request p, repo r, gh_user u
WHERE a.uuid = $1 AND 
  a.pr_id = p.pr_id AND 
  p.opened_by = u.id AND 
  p.repo_id = r.id
`

type GetCreateStatusInputsFromApprovalUuidRow struct {
	InstallationID int32  `json:"installation_id"`
	Login          string `json:"login"`
	Name           string `json:"name"`
	Sha            string `json:"sha"`
}

func (q *Queries) GetCreateStatusInputsFromApprovalUuid(ctx context.Context, uuid string) (GetCreateStatusInputsFromApprovalUuidRow, error) {
	row := q.db.QueryRow(ctx, getCreateStatusInputsFromApprovalUuid, uuid)
	var i GetCreateStatusInputsFromApprovalUuidRow
	err := row.Scan(
		&i.InstallationID,
		&i.Login,
		&i.Name,
		&i.Sha,
	)
	return i, err
}

const updateApprovalByUuid = `-- name: UpdateApprovalByUuid :exec
UPDATE approval SET is_approved = $2
WHERE uuid = $1
`

type UpdateApprovalByUuidParams struct {
	Uuid       string `json:"uuid"`
	IsApproved bool   `json:"is_approved"`
}

func (q *Queries) UpdateApprovalByUuid(ctx context.Context, arg UpdateApprovalByUuidParams) error {
	_, err := q.db.Exec(ctx, updateApprovalByUuid, arg.Uuid, arg.IsApproved)
	return err
}
