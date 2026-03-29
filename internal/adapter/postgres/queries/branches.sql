-- name: CreateBranch :one
INSERT INTO branches (name, address, city, status)
VALUES ($1, $2, $3, 'ACTIVE')
RETURNING id;

-- name: UpdateBranchStatus :exec
UPDATE branches
SET status = $1, updated_at = NOW()
WHERE id = $2;

-- name: GetAllBranches :many
SELECT id, name, address, city, status, created_at, updated_at
FROM branches;

-- name: GetBranchesByUserID :many
SELECT b.id, b.name, b.address, b.city, b.status, b.created_at, b.updated_at
FROM branches b
JOIN user_branches ub ON b.id = ub.branch_id
WHERE ub.user_id = $1;

-- name: GetBranchByID :one
SELECT id, name, address, city, status, created_at, updated_at
FROM branches
WHERE id = $1;

-- name: UpdateBranch :exec
UPDATE branches
SET name = $1, address = $2, city = $3, updated_at = NOW()
WHERE id = $4;
