-- name: CreateGithubUser :one
INSERT INTO gh_user(id, login)
VALUES 
  ($1, $2)
RETURNING *;

-- name: GetGithubUser :one
SELECT * FROM gh_user
WHERE id = $1 LIMIT 1;

-- name: GetGithubUsers :many
SELECT * FROM gh_user
ORDER BY login;

-- name: DeleteGithubUser :exec
DELETE FROM gh_user
WHERE id = $1;
