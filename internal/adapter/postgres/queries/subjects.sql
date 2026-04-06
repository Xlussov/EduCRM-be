-- name: CreateSubject :one
INSERT INTO subjects (branch_id, name, description)
VALUES ($1, $2, $3)
RETURNING id, branch_id, name, description, status, created_at, updated_at;

-- name: UpdateSubjectStatus :exec
UPDATE subjects
SET status = $1, updated_at = NOW()
WHERE id = $2;

-- name: GetSubject :one
SELECT id, branch_id, name, description, status, created_at, updated_at
FROM subjects
WHERE id = $1 AND branch_id = $2
LIMIT 1;

-- name: ListSubjects :many
SELECT id, branch_id, name, description, status, created_at, updated_at
FROM subjects
WHERE branch_id = $1
ORDER BY name ASC;

-- name: UpdateSubject :one
UPDATE subjects
SET branch_id = $1, name = $2, description = $3, updated_at = NOW()
WHERE id = $4 RETURNING *;
