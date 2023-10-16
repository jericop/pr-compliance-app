-- name: CreateApproval :one
INSERT INTO approval(uuid, pr_id, sha, is_approved)
VALUES 
  ($1, $2, $3, $4)
RETURNING *;

-- name: GetApprovalById :one
SELECT * FROM approval
WHERE id = $1 LIMIT 1;

-- name: GetApprovalByUuid :one
SELECT * FROM approval
WHERE uuid = $1 LIMIT 1;

-- name: UpdateApprovalByUuid :exec
UPDATE approval SET is_approved = $2
WHERE uuid = $1;

-- name: GetApprovalByPrIDSha :one
SELECT * FROM approval
WHERE pr_id = $1 AND sha = $2 LIMIT 1;

-- name: GetApprovals :many
SELECT * FROM approval;

-- name: DeleteApproval :exec
DELETE FROM approval
WHERE id = $1;
