package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StudentAttendanceReportItem struct {
	Date        time.Time
	Time        time.Time
	SubjectName string
	IsPresent   bool
	Notes       string
}

type StudentAttendanceSummary struct {
	TotalLessons         int
	Attended             int
	Missed               int
	AttendancePercentage float64
}

type BranchStatisticsReport struct {
	ActiveStudents       int
	CompletedLessons     int
	CancelledLessons     int
	AttendancePercentage float64
}

type TeacherStatisticsReport struct {
	ScheduledLessons int
	CompletedLessons int
	CancelledLessons int
}

type ReportRepository interface {
	GetStudentAttendanceHistory(ctx context.Context, studentID uuid.UUID, startDate, endDate *time.Time, subjectID *uuid.UUID) ([]StudentAttendanceReportItem, error)
	GetBranchStatistics(ctx context.Context, branchID uuid.UUID, startDate, endDate *time.Time) (*BranchStatisticsReport, error)
	GetTeacherStatistics(ctx context.Context, teacherID uuid.UUID, startDate, endDate *time.Time) (*TeacherStatisticsReport, error)
}
