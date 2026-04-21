-- name: CreateLesson :one
INSERT INTO lessons (
    branch_id, template_id, teacher_id, subject_id, student_id, group_id, date, start_time, end_time, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: CreateTemplate :one
INSERT INTO lesson_templates (
    branch_id, teacher_id, subject_id, student_id, group_id, days_of_week, start_time, end_time, start_date, end_date, is_active
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

-- name: UpdateLesson :one
UPDATE lessons
SET date = $2,
    start_time = $3,
    end_time = $4,
    teacher_id = $5,
    subject_id = $6
WHERE id = $1
RETURNING *;

-- name: GetLessonByID :one
SELECT * FROM lessons WHERE id = $1;

-- name: GetTemplateByID :one
SELECT * FROM lesson_templates WHERE id = $1;

-- name: CheckTeacherConflict :one
SELECT EXISTS (
    SELECT 1 FROM lessons
    WHERE teacher_id = sqlc.arg(teacher_id)
      AND date = sqlc.arg(date)
      AND start_time < sqlc.arg(new_end_time)
      AND sqlc.arg(new_start_time) < end_time
      AND status != 'CANCELLED'
);

-- name: CheckTeacherConflictExcludingLesson :one
SELECT EXISTS (
        SELECT 1 FROM lessons
        WHERE teacher_id = sqlc.arg(teacher_id)
            AND date = sqlc.arg(date)
            AND start_time < sqlc.arg(new_end_time)
            AND sqlc.arg(new_start_time) < end_time
            AND status != 'CANCELLED'
            AND id != sqlc.arg(exclude_lesson_id)
);

-- name: CheckStudentConflict :one
SELECT EXISTS (
    SELECT 1 FROM lessons l
    LEFT JOIN student_groups sg ON l.group_id = sg.group_id 
        AND sg.student_id = sqlc.arg(student_id) 
        AND sg.joined_at <= NOW() 
        AND (sg.left_at IS NULL OR sg.left_at > NOW())
    WHERE (l.student_id = sqlc.arg(student_id) OR sg.student_id = sqlc.arg(student_id))
      AND l.date = sqlc.arg(date)
      AND l.start_time < sqlc.arg(new_end_time)
      AND sqlc.arg(new_start_time) < l.end_time
      AND l.status != 'CANCELLED'
);

-- name: CheckStudentConflictExcludingLesson :one
SELECT EXISTS (
    SELECT 1 FROM lessons l
    LEFT JOIN student_groups sg ON l.group_id = sg.group_id 
        AND sg.student_id = sqlc.arg(student_id) 
        AND sg.joined_at <= NOW() 
        AND (sg.left_at IS NULL OR sg.left_at > NOW())
    WHERE (l.student_id = sqlc.arg(student_id) OR sg.student_id = sqlc.arg(student_id))
      AND l.date = sqlc.arg(date)
      AND l.start_time < sqlc.arg(new_end_time)
      AND sqlc.arg(new_start_time) < l.end_time
      AND l.status != 'CANCELLED'
      AND l.id != sqlc.arg(exclude_lesson_id)
);

-- name: GetTeacherSchedule :many
SELECT * FROM lessons 
WHERE teacher_id = $1 
  AND date >= $2 
  AND date <= $3
ORDER BY date ASC, start_time ASC;

-- name: ListLessons :many
SELECT
        l.id,
        l.branch_id,
        l.template_id,
        l.teacher_id,
        t.first_name AS teacher_first_name,
        t.last_name AS teacher_last_name,
        l.subject_id,
        s.name AS subject_name,
        l.student_id,
        st.first_name AS student_first_name,
        st.last_name AS student_last_name,
        l.group_id,
        g.name AS group_name,
        l.date,
        l.start_time,
        l.end_time,
        l.status,
        l.created_at
FROM lessons l
JOIN users t ON t.id = l.teacher_id
JOIN subjects s ON s.id = l.subject_id
LEFT JOIN students st ON st.id = l.student_id
LEFT JOIN groups g ON g.id = l.group_id
WHERE l.date >= $1
    AND l.date <= $2
    AND ($3::uuid IS NULL OR l.teacher_id = $3)
    AND ($4::uuid IS NULL OR l.student_id = $4)
    AND ($5::uuid IS NULL OR l.group_id = $5)
    AND ($6::uuid[] IS NULL OR l.branch_id = ANY($6))
ORDER BY l.date ASC, l.start_time ASC;

-- name: DeactivateTemplate :exec
UPDATE lesson_templates
SET is_active = FALSE
WHERE id = $1;

-- name: CancelFutureLessonsByTemplate :exec
UPDATE lessons
SET status = 'CANCELLED'
WHERE template_id = $1
    AND date >= CURRENT_DATE
    AND status = 'SCHEDULED';

-- name: CheckTeacherFutureLessonsInBranch :one
SELECT EXISTS (
        SELECT 1 FROM lessons
        WHERE teacher_id = $1
            AND branch_id = $2
            AND status = 'SCHEDULED'
            AND date >= CURRENT_DATE
);

-- name: CheckTeacherActiveTemplatesInBranch :one
SELECT EXISTS (
        SELECT 1 FROM lesson_templates
        WHERE teacher_id = $1
            AND branch_id = $2
            AND is_active = TRUE
            AND end_date >= CURRENT_DATE
);
