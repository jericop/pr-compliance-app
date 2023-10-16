-- name: CreatePullRequestAction :one
INSERT INTO pull_request_action(name)
VALUES 
  ($1)
RETURNING *;

-- name: GetPullRequestAction :one
SELECT * FROM pull_request_action
WHERE name = $1 LIMIT 1;

-- name: GetPullRequestActions :many
SELECT * FROM pull_request_action;

-- name: DeletePullRequestAction :exec
DELETE FROM pull_request_action
WHERE name = $1;
