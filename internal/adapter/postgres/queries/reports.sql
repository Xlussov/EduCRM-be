-- name: GetStudentAttendanceHistory :many
SELECT
    l.date,
    l.start_time,
    sub.name AS subject_name,
    COALESCE(a.is_present, false)::boolean AS is_present,
    COALESCE(a.notes, '')::text AS notes
FROM lessons l
JOIN subjects sub ON l.subject_id = sub.id
LEFT JOIN attendance a ON l.id = a.lesson_id AND a.student_id = sqlc.arg(student_id)
WHERE (l.student_id = sqlc.arg(student_id) OR EXISTS (
    SELECT 1 FROM student_groups sg 
    WHERE sg.group_id = l.group_id 
    AND sg.student_id = sqlc.arg(student_id)
    AND sg.joined_at::date <= l.date
    AND (sg.left_at IS NULL OR sg.left_at::date > l.date)
))
AND l.status = 'COMPLETED'
AND (sqlc.narg(start_date)::date IS NULL OR l.date >= sqlc.narg(start_date)::date)
AND (sqlc.narg(end_date)::date IS NULL OR l.date <= sqlc.narg(end_date)::date)
AND (sqlc.narg(subject_id)::uuid IS NULL OR l.subject_id = sqlc.narg(subject_id)::uuid)
ORDER BY l.date DESC, l.start_time DESC;

-- name: CountActiveStudentsByBranch :one
SELECT COUNT(id) 
FROM students 
WHERE branch_id = sqlc.arg(branch_id) 
AND status = 'ACTIVE';

-- name: CountCompletedLessonsByBranch :one
SELECT COUNT(id) 
FROM lessons 
WHERE branch_id = sqlc.arg(branch_id) 
AND status = 'COMPLETED'
AND (sqlc.narg(start_date)::date IS NULL OR date >= sqlc.narg(start_date)::date)
AND (sqlc.narg(end_date)::date IS NULL OR date <= sqlc.narg(end_date)::date);

-- name: CountCancelledLessonsByBranch :one
SELECT COUNT(id) 
FROM lessons 
WHERE branch_id = sqlc.arg(branch_id) 
AND status = 'CANCELLED'
AND (sqlc.narg(start_date)::date IS NULL OR date >= sqlc.narg(start_date)::date)
AND (sqlc.narg(end_date)::date IS NULL OR date <= sqlc.narg(end_date)::date);

-- name: GetBranchAttendanceStats :one
SELECT 
    COUNT(a.lesson_id)::bigint AS total_attendance_records,
    COUNT(a.lesson_id) FILTER (WHERE a.is_present)::bigint AS total_present_records
FROM attendance a
JOIN lessons l ON a.lesson_id = l.id
WHERE l.branch_id = sqlc.arg(branch_id)
AND (sqlc.narg(start_date)::date IS NULL OR l.date >= sqlc.narg(start_date)::date)
AND (sqlc.narg(end_date)::date IS NULL OR l.date <= sqlc.narg(end_date)::date);

-- name: GetTeacherStatistics :one
SELECT 
    COUNT(*) FILTER (WHERE status = 'SCHEDULED')::int AS scheduled_lessons,
    COUNT(*) FILTER (WHERE status = 'COMPLETED')::int AS completed_lessons,
    COUNT(*) FILTER (WHERE status = 'CANCELLED')::int AS cancelled_lessons
FROM lessons
WHERE teacher_id = sqlc.arg(teacher_id)
AND (sqlc.narg(start_date)::date IS NULL OR date >= sqlc.narg(start_date)::date)
AND (sqlc.narg(end_date)::date IS NULL OR date <= sqlc.narg(end_date)::date);