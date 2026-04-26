package mocks

import (
	"context"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) GetStudentAttendanceHistory(ctx context.Context, studentID uuid.UUID, startDate, endDate *time.Time, subjectID *uuid.UUID) ([]domain.StudentAttendanceReportItem, error) {
	args := m.Called(ctx, studentID, startDate, endDate, subjectID)

	var items []domain.StudentAttendanceReportItem
	if args.Get(0) != nil {
		items = args.Get(0).([]domain.StudentAttendanceReportItem)
	}

	return items, args.Error(1)
}

func (m *MockReportRepository) GetBranchStatistics(ctx context.Context, branchID uuid.UUID, startDate, endDate *time.Time) (*domain.BranchStatisticsReport, error) {
	args := m.Called(ctx, branchID, startDate, endDate)

	var report *domain.BranchStatisticsReport
	if args.Get(0) != nil {
		report = args.Get(0).(*domain.BranchStatisticsReport)
	}

	return report, args.Error(1)
}

func (m *MockReportRepository) GetTeacherStatistics(ctx context.Context, teacherID uuid.UUID, startDate, endDate *time.Time) (*domain.TeacherStatisticsReport, error) {
	args := m.Called(ctx, teacherID, startDate, endDate)

	var report *domain.TeacherStatisticsReport
	if args.Get(0) != nil {
		report = args.Get(0).(*domain.TeacherStatisticsReport)
	}

	return report, args.Error(1)
}
