
-- name: GetSortedApprovalYesNoQuestionsBySchemaId :many
SELECT id, question_text FROM approval_yesno_question
WHERE schema_id = (SELECT id FROM approval_schema WHERE name = $1) 
GROUP BY id, question_text;


-- name: GetSortedApprovalYesNoQuestionAnswersByUuid :many
WITH false_answers AS (
  SELECT q.id, q.question_text
  FROM approval_yesno_question q, approval a
  WHERE 
    a.uuid = $1 AND
    q.schema_id = a.schema_id
), true_answers AS (
  SELECT ya.question_id
  FROM approval a, approval_yes_answer ya
  WHERE 
    a.uuid = $1 AND 
    ya.approval_id = a.id
)
SELECT 
  f.id, 
  f.question_text, 
  CASE WHEN EXISTS (SELECT 1 FROM true_answers WHERE question_id=f.id) THEN true ELSE false END as answered_yes
FROM false_answers f
ORDER BY f.id;
