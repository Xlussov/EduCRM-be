-- name: CreateSubscriptionPlan :one
INSERT INTO subscription_plans (branch_id, name, type, status)
VALUES ($1, $2, $3, $4)
RETURNING id, branch_id, name, type, status, created_at;

-- name: CreatePlanSubject :exec
INSERT INTO plan_subjects (plan_id, subject_id)
VALUES ($1, $2);

-- name: CreatePlanPricingGrid :exec
INSERT INTO plan_pricing_grid (plan_id, lessons_per_month, price_per_lesson)
VALUES ($1, $2, $3);

-- name: GetPlansByBranchID :many
SELECT id, branch_id, name, type, status, created_at
FROM subscription_plans
WHERE branch_id = $1
ORDER BY created_at DESC;

-- name: GetPlanSubjects :many
SELECT s.id, s.name, s.description, s.status, s.created_at, s.updated_at
FROM subjects s
JOIN plan_subjects ps ON s.id = ps.subject_id
WHERE ps.plan_id = $1;

-- name: GetPlanPricingGrids :many
SELECT id, plan_id, lessons_per_month, price_per_lesson
FROM plan_pricing_grid
WHERE plan_id = $1
ORDER BY lessons_per_month ASC;

-- name: ValidatePlanSubject :one
SELECT EXISTS (
    SELECT 1 
    FROM plan_subjects ps
    JOIN subscription_plans p ON p.id = ps.plan_id
    JOIN subjects s ON s.id = ps.subject_id
    WHERE ps.plan_id = $1
      AND ps.subject_id = $2
      AND p.status = 'ACTIVE'
      AND s.status = 'ACTIVE'
);

-- name: UpdatePlanStatus :exec
UPDATE subscription_plans
SET status = $2
WHERE id = $1;

-- name: GetPlanByID :one
SELECT id, branch_id, name, type, status, created_at
FROM subscription_plans
WHERE id = $1 LIMIT 1;

-- name: CountSubjectsInBranch :one
SELECT COUNT(*)::int
FROM unnest($2::uuid[]) AS subject_ids(id)
JOIN subjects s ON s.id = subject_ids.id
WHERE s.branch_id = $1
    AND s.status = 'ACTIVE';