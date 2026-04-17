package create_template

import "github.com/google/uuid"

type Request struct {
	BranchID  uuid.UUID  `json:"branch_id" validate:"required"`
	TeacherID uuid.UUID  `json:"teacher_id" validate:"required"`
	SubjectID uuid.UUID  `json:"subject_id" validate:"required"`
	StudentID *uuid.UUID `json:"student_id"`
	GroupID   *uuid.UUID `json:"group_id"`
	DayOfWeek int        `json:"day_of_week" validate:"min=0,max=6"`
	StartTime string     `json:"start_time" validate:"required"`
	EndTime   string     `json:"end_time" validate:"required"`
	StartDate string     `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string     `json:"end_date" validate:"required,datetime=2006-01-02"`
}

type Response struct {
	TemplateID          uuid.UUID `json:"template_id"`
	CreatedLessonsCount int       `json:"created_lessons_count"`
	Conflicts           []string  `json:"conflicts"`
}
