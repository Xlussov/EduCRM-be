-- name: CreateStudent :one
INSERT INTO students (
    branch_id, first_name, last_name, dob, phone, email, address,
    parent_name, parent_phone, parent_email, parent_relationship, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'ACTIVE'
) RETURNING id;

-- name: UpdateStudentStatus :exec
UPDATE students
SET status = $1
WHERE id = $2;

-- name: GetStudentBranchID :one
SELECT branch_id FROM students WHERE id = $1;

-- name: GetStudentByID :one
SELECT id, branch_id, first_name, last_name, dob, phone, email, address, parent_name, parent_phone, parent_email, parent_relationship, status, created_at
FROM students 
WHERE id = $1;

-- name: UpdateStudent :one
UPDATE students
SET first_name = $1, last_name = $2, dob = $3, phone = $4, email = $5,
    address = $6, parent_name = $7, parent_phone = $8, parent_email = $9,
    parent_relationship = $10
WHERE id = $11 RETURNING *;

-- name: GetStudentsByBranchID :many
SELECT id, branch_id, first_name, last_name, dob, phone, email, address, parent_name, parent_phone, parent_email, parent_relationship, status, created_at
FROM students
WHERE branch_id = $1
    AND (sqlc.narg(status)::entity_status IS NULL OR status = sqlc.narg(status)::entity_status)
ORDER BY created_at DESC;

-- name: GetStudentsByBranchIDAndTeacherID :many
SELECT DISTINCT s.id, s.branch_id, s.first_name, s.last_name, s.dob, s.phone, s.email, s.address, s.parent_name, s.parent_phone, s.parent_email, s.parent_relationship, s.status, s.created_at
FROM students s
WHERE s.branch_id = $1
    AND (sqlc.narg(status)::entity_status IS NULL OR s.status = sqlc.narg(status)::entity_status)
    AND (
        EXISTS (
            SELECT 1
            FROM lessons l
            WHERE l.teacher_id = $2 AND l.student_id = s.id
        )
        OR EXISTS (
            SELECT 1
            FROM lessons l
            JOIN student_groups sg ON sg.group_id = l.group_id AND sg.left_at IS NULL
            WHERE l.teacher_id = $2 AND sg.student_id = s.id
        )
    )
ORDER BY s.created_at DESC;

-- name: IsTeacherStudent :one
SELECT (
    EXISTS (
        SELECT 1
        FROM lessons l
        WHERE l.teacher_id = $1 AND l.student_id = $2
    )
    OR EXISTS (
        SELECT 1
        FROM lessons l
        JOIN student_groups sg ON sg.group_id = l.group_id AND sg.left_at IS NULL
        WHERE l.teacher_id = $1 AND sg.student_id = $2
    )
) AS is_teacher_student;
