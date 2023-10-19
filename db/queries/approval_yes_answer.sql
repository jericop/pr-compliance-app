-- name: CreateApprovalYesAnswer :one
INSERT INTO approval_yes_answer (approval_id, question_id)
VALUES 
  ($1, $2)
RETURNING *;

-- name: DeleteApprovalYesAnswer :exec
DELETE FROM approval_yes_answer
WHERE approval_id = $1 AND question_id = $2;

-- name: CreateApprovalYesAnswerByUuid :one
INSERT INTO approval_yes_answer (approval_id, question_id)
VALUES 
  ((SELECT id from approval WHERE uuid = $1), $2)
RETURNING *;

-- name: DeleteApprovalYesAnswerByUuid :exec
DELETE FROM approval_yes_answer
WHERE question_id = $2 AND
  approval_id = (SELECT id from approval WHERE uuid = $1);
