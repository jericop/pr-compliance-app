-- name: CreateRepo :one
INSERT INTO repo(org, name, id )
VALUES 
  ($1, $2, $3)
RETURNING *;

-- name: GetRepo :one
SELECT * FROM repo
WHERE id = $1 LIMIT 1;

-- name: GetRepos :many
SELECT * FROM repo
ORDER BY org, name;

-- name: UpdateRepoName :one
UPDATE repo
SET name = $2
WHERE id = $1
RETURNING *;

-- name: UpdateRepoOrg :one
UPDATE repo
SET org = $2
WHERE id = $1
RETURNING *;

-- name: DeleteRepo :exec
DELETE FROM repo
WHERE id = $1;
