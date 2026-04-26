package branch_statistics

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	BranchID  uuid.UUID  `query:"branch_id" validate:"required"`
	StartDate *time.Time `query:"start_date"`
	EndDate   *time.Time `query:"end_date"`
}

type Response struct {
	ActiveStudents       int     `json:"active_students"`
	CompletedLessons     int     `json:"completed_lessons"`
	CancelledLessons     int     `json:"cancelled_lessons"`
	AttendancePercentage float64 `json:"attendance_percentage"`
}
