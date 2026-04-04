-- name: CreateStudentSubscription :one
INSERT INTO student_subscriptions (student_id, plan_id, subject_id, start_date)
VALUES ($1, $2, $3, $4)
RETURNING id, student_id, plan_id, subject_id, start_date, created_at;

-- name: GetStudentSubscriptions :many
SELECT 
    ss.id, 
    ss.student_id,
    ss.plan_id,
    ss.subject_id,
    ss.start_date,
    ss.created_at,
    p.name AS plan_name,
    s.name AS subject_name
FROM student_subscriptions ss
JOIN subscription_plans p ON ss.plan_id = p.id
JOIN subjects s ON ss.subject_id = s.id
WHERE ss.student_id = $1
ORDER BY ss.start_date DESC;
