-- name: UpsertAttendance :exec
INSERT INTO attendance (lesson_id, student_id, is_present, notes)
SELECT
    $1::uuid,
    unnest($2::uuid[]),
    unnest($3::boolean[]),
    unnest($4::text[])
ON CONFLICT (lesson_id, student_id)
DO UPDATE SET
    is_present = EXCLUDED.is_present,
    notes = EXCLUDED.notes;

-- name: GetLessonAttendance :many
SELECT DISTINCT
    s.id AS student_id,
    s.first_name,
    s.last_name,
    s.status,
    a.is_present,
    a.notes
FROM lessons l
LEFT JOIN student_groups sg
    ON sg.group_id = l.group_id
    AND sg.joined_at::date <= l.date
    AND (sg.left_at IS NULL OR sg.left_at::date > l.date)
JOIN students s
    ON (l.student_id = s.id OR sg.student_id = s.id)
LEFT JOIN attendance a
    ON a.lesson_id = l.id AND a.student_id = s.id
WHERE l.id = $1
ORDER BY s.last_name, s.first_name;
