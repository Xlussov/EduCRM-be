package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type AttendanceRepository struct {
	mock.Mock
}

func (m *AttendanceRepository) UpsertAttendance(ctx context.Context, attendance []domain.Attendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *AttendanceRepository) GetLessonAttendance(ctx context.Context, lessonID uuid.UUID) ([]domain.LessonAttendanceStudent, error) {
	args := m.Called(ctx, lessonID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.LessonAttendanceStudent), args.Error(1)
}

func (m *AttendanceRepository) GetStudentAttendance(ctx context.Context, studentID uuid.UUID, from, to time.Time) ([]domain.Attendance, error) {
	args := m.Called(ctx, studentID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Attendance), args.Error(1)
}
