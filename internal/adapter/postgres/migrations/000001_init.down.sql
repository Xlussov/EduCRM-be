-- DOWN Migration

DROP INDEX IF EXISTS idx_lessons_group_date;
DROP INDEX IF EXISTS idx_lessons_student_date;
DROP INDEX IF EXISTS idx_lessons_teacher_date;
DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_subjects_branch_id;

DROP TABLE IF EXISTS attendance;
DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS lesson_templates;
DROP TABLE IF EXISTS student_subscriptions;
DROP TABLE IF EXISTS plan_pricing_grid;
DROP TABLE IF EXISTS plan_subjects;
DROP TABLE IF EXISTS subscription_plans;
DROP TABLE IF EXISTS student_groups;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS user_branches;
DROP TABLE IF EXISTS branches;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS plan_type;
DROP TYPE IF EXISTS lesson_status;
DROP TYPE IF EXISTS entity_status;
DROP TYPE IF EXISTS user_role;

DROP EXTENSION IF EXISTS "uuid-ossp";