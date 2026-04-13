-- name: CreateGroup :one
INSERT INTO groups (branch_id, name, status)
VALUES ($1, $2, 'ACTIVE')
RETURNING id;

-- name: GetGroupsByBranchID :many
SELECT 
    g.id, 
    g.name, 
    g.status, 
    COALESCE(COUNT(sg.student_id), 0)::int as students_count
FROM groups g
LEFT JOIN student_groups sg ON g.id = sg.group_id AND sg.left_at IS NULL
WHERE g.branch_id = $1
    AND (sqlc.narg(status)::entity_status IS NULL OR g.status = sqlc.narg(status)::entity_status)
GROUP BY g.id, g.name, g.status
ORDER BY g.created_at DESC;

-- name: GetGroupByID :one
SELECT id, branch_id, name, status, created_at
FROM groups
WHERE id = $1;

-- name: UpdateGroupName :one
UPDATE groups
SET name = $1
WHERE id = $2
RETURNING *;

-- name: AddStudentToGroup :exec
INSERT INTO student_groups (student_id, group_id, joined_at)
VALUES ($1, $2, $3);

-- name: RemoveStudentFromGroup :exec
UPDATE student_groups
SET left_at = $3
WHERE student_id = $1 AND group_id = $2 AND left_at IS NULL;

-- name: GetGroupActiveStudentIDs :many
SELECT student_id
FROM student_groups
WHERE group_id = $1 AND left_at IS NULL;

-- name: GetGroupStudents :many
SELECT s.id, s.first_name, s.last_name
FROM students s
JOIN student_groups sg ON s.id = sg.student_id
WHERE sg.group_id = $1 AND sg.left_at IS NULL
ORDER BY s.last_name, s.first_name;

-- name: GetGroupBranchID :one
SELECT branch_id FROM groups WHERE id = $1;

-- name: UpdateGroupStatus :exec
UPDATE groups
SET status = $1
WHERE id = $2;