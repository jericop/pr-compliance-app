-- name: CreatePullRequest :one
INSERT INTO pull_request(repo_id, pr_id, pr_number, opened_by, installation_id, is_merged)
VALUES 
  ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPullRequestById :one
SELECT * FROM pull_request
WHERE id = $1 LIMIT 1;

-- name: GetPullRequestByRepoIdPrId :one
SELECT * FROM pull_request
WHERE repo_id = $1 AND pr_id = $2 LIMIT 1;

-- name: GetPullRequests :many
SELECT * FROM pull_request;

-- name: GetPullRequestForUpdate :one
SELECT * FROM pull_request
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdatePullRequestIsMerged :one
UPDATE pull_request
SET is_merged = $2
WHERE id = $1
RETURNING *;

-- name: DeletePullRequest :exec
DELETE FROM pull_request
WHERE id = $1;
