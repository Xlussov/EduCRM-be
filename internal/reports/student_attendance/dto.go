package student_attendance

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	StudentID uuid.UUID  `query:"student_id" validate:"required"`
	StartDate *time.Time `query:"start_date"`
	EndDate   *time.Time `query:"end_date"`
	SubjectID *uuid.UUID `query:"subject_id"`
}

type ReportItem struct {
	Date        time.Time `json:"date"`
	Time        time.Time `json:"time"`
	SubjectName string    `json:"subject_name"`
	IsPresent   bool      `json:"is_present"`
	Notes       string    `json:"notes"`
}

type Summary struct {
	TotalLessons         int     `json:"total_lessons"`
	Attended             int     `json:"attended"`
	Missed               int     `json:"missed"`
	AttendancePercentage float64 `json:"attendance_percentage"`
}

type Response struct {
	Items   []ReportItem `json:"items"`
	Summary Summary      `json:"summary"`
}
