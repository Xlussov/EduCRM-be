-- name: CreateLesson :one
INSERT INTO lessons (
    branch_id, template_id, teacher_id, subject_id, student_id, group_id, date, start_time, end_time, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: CreateTemplate :one
INSERT INTO lesson_templates (
    branch_id, teacher_id, subject_id, student_id, group_id, day_of_week, start_time, end_time, start_date, end_date, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: BulkCreateLessons :copyfrom
INSERT INTO lessons (
    branch_id, template_id, teacher_id, subject_id, student_id, group_id, date, start_time, end_time, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
);

-- name: UpdateLessonStatus :exec
UPDATE lessons
SET status = $2
WHERE id = $1;

-- name: GetLessonByID :one
SELECT * FROM lessons WHERE id = $1;

-- name: CheckTeacherConflict :one
SELECT EXISTS (
    SELECT 1 FROM lessons
    WHERE teacher_id = $1
      AND date = $2
      AND status != 'CANCELLED'
      AND start_time < $4
      AND end_time > $3
);

-- name: CheckStudentConflict :one
SELECT EXISTS (
    SELECT 1 FROM lessons l
    LEFT JOIN student_groups sg ON l.group_id = sg.group_id 
        AND sg.student_id = $1 
        AND sg.joined_at <= NOW() 
        AND (sg.left_at IS NULL OR sg.left_at > NOW())
    WHERE (l.student_id = $1 OR sg.student_id = $1)
      AND l.date = $2
      AND l.status != 'CANCELLED'
      AND l.start_time < $4
      AND l.end_time > $3
);

-- name: GetTeacherSchedule :many
SELECT * FROM lessons 
WHERE teacher_id = $1 
  AND date >= $2 
  AND date <= $3
ORDER BY date ASC, start_time ASC;
