-- UP Migration

ALTER TABLE lesson_templates
ADD COLUMN days_of_week INT[];

UPDATE lesson_templates
SET days_of_week = ARRAY[day_of_week];

ALTER TABLE lesson_templates
ALTER COLUMN days_of_week SET NOT NULL;

ALTER TABLE lesson_templates
DROP COLUMN day_of_week;
