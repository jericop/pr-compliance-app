-- name: CreateApproval :one
INSERT INTO approval(uuid, pr_id, sha, approved_on)
VALUES 
  ($1, $2, $3, $4)
RETURNING *;

-- name: GetApprovalById :one
SELECT * FROM approval
WHERE id = $1 LIMIT 1;

-- name: GetApprovalByUuid :one
SELECT * FROM approval
WHERE uuid = $1 LIMIT 1;

-- name: GetApprovals :many
SELECT * FROM approval;

-- name: DeleteApproval :exec
DELETE FROM approval
WHERE id = $1;
