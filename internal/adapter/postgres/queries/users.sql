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

-- name: CheckTeacherInBranch :one
SELECT EXISTS (
	SELECT 1
	FROM user_branches ub
	JOIN users u ON u.id = ub.user_id
	WHERE ub.user_id = $1
	  AND ub.branch_id = $2
	  AND u.role = 'TEACHER'
);

-- name: GetUserByID :one
SELECT id, phone, password_hash, first_name, last_name, role, is_active, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetAdmins :many
SELECT
	u.id,
	u.phone,
	u.first_name,
	u.last_name,
	u.role,
	u.is_active,
	u.created_at,
	u.updated_at,
	b.id AS branch_id,
	b.name AS branch_name
FROM users u
LEFT JOIN user_branches ub ON ub.user_id = u.id
LEFT JOIN branches b ON b.id = ub.branch_id
WHERE u.role = 'ADMIN'
ORDER BY u.created_at DESC, b.name ASC;

-- name: GetTeachers :many
SELECT
	u.id,
	u.phone,
	u.first_name,
	u.last_name,
	u.role,
	u.is_active,
	u.created_at,
	u.updated_at,
	b.id AS branch_id,
	b.name AS branch_name
FROM users u
LEFT JOIN user_branches ub ON ub.user_id = u.id
LEFT JOIN branches b ON b.id = ub.branch_id
WHERE u.role = 'TEACHER'
  AND (
	  $1::uuid[] IS NULL
	  OR ub.branch_id = ANY($1::uuid[])
  )
ORDER BY u.created_at DESC, b.name ASC;

-- name: GetUserWithBranchesByID :many
SELECT
	u.id,
	u.phone,
	u.first_name,
	u.last_name,
	u.role,
	u.is_active,
	u.created_at,
	u.updated_at,
	b.id AS branch_id,
	b.name AS branch_name
FROM users u
LEFT JOIN user_branches ub ON ub.user_id = u.id
LEFT JOIN branches b ON b.id = ub.branch_id
WHERE u.id = $1
ORDER BY b.name ASC;

-- name: UpdateUser :execrows
UPDATE users
SET
	phone = $2,
	first_name = $3,
	last_name = $4,
	updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserStatus :execrows
UPDATE users
SET
	is_active = $2,
	updated_at = NOW()
WHERE id = $1;

-- name: DeleteUserBranches :exec
DELETE FROM user_branches
WHERE user_id = $1;