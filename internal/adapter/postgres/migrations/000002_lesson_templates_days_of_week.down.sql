-- DOWN Migration

ALTER TABLE lesson_templates
ADD COLUMN day_of_week INT;

UPDATE lesson_templates
SET day_of_week = days_of_week[1];

ALTER TABLE lesson_templates
ALTER COLUMN day_of_week SET NOT NULL;

ALTER TABLE lesson_templates
DROP COLUMN days_of_week;
