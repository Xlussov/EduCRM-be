-- name: CreateBranch :one
INSERT INTO branches (name, address, city, status)
VALUES ($1, $2, $3, 'ACTIVE')
RETURNING id;

-- name: UpdateBranchStatus :exec
UPDATE branches
SET status = $1, updated_at = NOW()
WHERE id = $2;
