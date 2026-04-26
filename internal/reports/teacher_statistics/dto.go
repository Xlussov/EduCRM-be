package teacher_statistics

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	TeacherID uuid.UUID  `query:"teacher_id"`
	StartDate *time.Time `query:"start_date"`
	EndDate   *time.Time `query:"end_date"`
}

type Response struct {
	ScheduledLessons int `json:"scheduled_lessons"`
	CompletedLessons int `json:"completed_lessons"`
	CancelledLessons int `json:"cancelled_lessons"`
}
