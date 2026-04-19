package update

import "github.com/google/uuid"

type Request struct {
	Date      string    `json:"date" validate:"required,datetime=2006-01-02"`
	StartTime string    `json:"start_time" validate:"required"`
	EndTime   string    `json:"end_time" validate:"required"`
	TeacherID uuid.UUID `json:"teacher_id" validate:"required"`
	SubjectID uuid.UUID `json:"subject_id" validate:"required"`
}

type Response struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}
