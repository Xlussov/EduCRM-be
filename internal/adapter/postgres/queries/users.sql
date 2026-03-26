-- name: CreateUser :one
INSERT INTO users (phone, password_hash, first_name, last_name, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByPhone :one
SELECT id, phone, password_hash, first_name, last_name, role, is_active, created_at, updated_at
FROM users
WHERE phone = $1 LIMIT 1;

-- name: AssignUserToBranch :exec
INSERT INTO user_branches (user_id, branch_id)
VALUES ($1, $2);

-- name: GetUserBranchIDs :many
SELECT branch_id
FROM user_branches
WHERE user_id = $1;

-- name: GetUserByID :one
SELECT id, phone, password_hash, first_name, last_name, role, is_active, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;