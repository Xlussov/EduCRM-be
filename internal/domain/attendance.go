package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	ID        uuid.UUID
	LessonID  uuid.UUID
	StudentID uuid.UUID
	IsPresent bool
	Notes     *string
	CreatedAt time.Time
}

type LessonAttendanceStudent struct {
	StudentID uuid.UUID
	FirstName string
	LastName  string
	Status    EntityStatus
	IsPresent *bool
	Notes     *string
}

type AttendanceRepository interface {
	UpsertAttendance(ctx context.Context, attendance []Attendance) error
	GetLessonAttendance(ctx context.Context, lessonID uuid.UUID) ([]LessonAttendanceStudent, error)
	GetStudentAttendance(ctx context.Context, studentID uuid.UUID, from, to time.Time) ([]Attendance, error)
}
