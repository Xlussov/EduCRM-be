package get

import "github.com/google/uuid"

type Request struct {
	ID uuid.UUID
}

type Response struct {
	ID         uuid.UUID  `json:"id"`
	BranchID   uuid.UUID  `json:"branch_id"`
	TemplateID *uuid.UUID `json:"template_id"`
	TeacherID  uuid.UUID  `json:"teacher_id"`
	SubjectID  uuid.UUID  `json:"subject_id"`
	StudentID  *uuid.UUID `json:"student_id,omitempty"`
	GroupID    *uuid.UUID `json:"group_id,omitempty"`
	Date       string     `json:"date"`
	StartTime  string     `json:"start_time"`
	EndTime    string     `json:"end_time"`
	Status     string     `json:"status"`
}
