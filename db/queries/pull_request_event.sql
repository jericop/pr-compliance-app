-- name: CreatePullRequestEvent :one
INSERT INTO pull_request_event(pr_id, action, sha, is_merged, last_updated)
VALUES 
  ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPullRequestEvent :one
SELECT * FROM pull_request_event
WHERE id = $1 LIMIT 1;

-- name: GetPullRequestEvents :many
SELECT * FROM pull_request_event;

-- name: DeletePullRequestEvent :exec
DELETE FROM pull_request_event
WHERE id = $1;
