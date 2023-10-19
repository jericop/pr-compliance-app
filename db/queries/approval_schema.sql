-- name: GetDefaultApprovalSchema :one
SELECT * FROM approval_schema 
WHERE id = (SELECT schema_id FROM default_approval_schema LIMIT 1);

-- name: GetApprovalSchemaById :one
SELECT * FROM approval_schema WHERE id = $1;
