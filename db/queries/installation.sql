-- name: CreateInstallation :one
INSERT INTO installation(id)
VALUES 
  ($1)
RETURNING *;

-- name: GetInstallation :one
SELECT * FROM installation
WHERE id = $1 LIMIT 1;

-- name: GetInstallations :many
SELECT * FROM installation;

-- name: DeleteInstallation :exec
DELETE FROM installation
WHERE id = $1;
