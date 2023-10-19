-- name: CreateApproval :one
INSERT INTO approval(schema_id, uuid, pr_id, sha, is_approved)
VALUES 
  ($1, $2, $3, $4, $5)
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

-- name: GetCreateStatusInputsFromApprovalUuid :one
SELECT p.installation_id, u.login, r.name, a.sha
FROM approval a, pull_request p, repo r, gh_user u
WHERE a.uuid = $1 AND 
  a.pr_id = p.pr_id AND 
  p.opened_by = u.id AND 
  p.repo_id = r.id;
