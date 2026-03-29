-- name: CreateSubject :one
INSERT INTO subjects (name, description)
VALUES ($1, $2)
RETURNING id;

-- name: UpdateSubjectStatus :exec
UPDATE subjects
SET status = $1, updated_at = NOW()
WHERE id = $2;
