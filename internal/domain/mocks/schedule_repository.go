package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type ScheduleRepository struct {
	mock.Mock
}

func (m *ScheduleRepository) CreateLesson(ctx context.Context, lesson *domain.Lesson) error {
	args := m.Called(ctx, lesson)
	return args.Error(0)
}

func (m *ScheduleRepository) CreateTemplate(ctx context.Context, template *domain.Template) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *ScheduleRepository) BulkCreateLessons(ctx context.Context, lessons []domain.Lesson) error {
	args := m.Called(ctx, lessons)
	return args.Error(0)
}

func (m *ScheduleRepository) UpdateLessonStatus(ctx context.Context, id uuid.UUID, status domain.LessonStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *ScheduleRepository) UpdateLesson(ctx context.Context, lesson *domain.Lesson) error {
	args := m.Called(ctx, lesson)
	return args.Error(0)
}

func (m *ScheduleRepository) GetLessonByID(ctx context.Context, id uuid.UUID) (*domain.Lesson, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Lesson), args.Error(1)
}

func (m *ScheduleRepository) GetTemplateByID(ctx context.Context, id uuid.UUID) (*domain.Template, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Template), args.Error(1)
}

func (m *ScheduleRepository) DeactivateTemplate(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ScheduleRepository) CancelFutureLessonsByTemplate(ctx context.Context, templateID uuid.UUID) error {
	args := m.Called(ctx, templateID)
	return args.Error(0)
}

func (m *ScheduleRepository) CheckTeacherConflict(ctx context.Context, teacherID uuid.UUID, date time.Time, start, end time.Time) (bool, error) {
	args := m.Called(ctx, teacherID, date, start, end)
	return args.Bool(0), args.Error(1)
}

func (m *ScheduleRepository) CheckTeacherConflictExcludingLesson(ctx context.Context, teacherID uuid.UUID, date time.Time, start, end time.Time, lessonID uuid.UUID) (bool, error) {
	args := m.Called(ctx, teacherID, date, start, end, lessonID)
	return args.Bool(0), args.Error(1)
}

func (m *ScheduleRepository) CheckStudentConflict(ctx context.Context, studentID uuid.UUID, date time.Time, start, end time.Time) (bool, error) {
	args := m.Called(ctx, studentID, date, start, end)
	return args.Bool(0), args.Error(1)
}

func (m *ScheduleRepository) CheckStudentConflictExcludingLesson(ctx context.Context, studentID uuid.UUID, date time.Time, start, end time.Time, lessonID uuid.UUID) (bool, error) {
	args := m.Called(ctx, studentID, date, start, end, lessonID)
	return args.Bool(0), args.Error(1)
}

func (m *ScheduleRepository) CheckTeacherFutureLessonsInBranch(ctx context.Context, teacherID, branchID uuid.UUID) (bool, error) {
	args := m.Called(ctx, teacherID, branchID)
	return args.Bool(0), args.Error(1)
}

func (m *ScheduleRepository) CheckTeacherActiveTemplatesInBranch(ctx context.Context, teacherID, branchID uuid.UUID) (bool, error) {
	args := m.Called(ctx, teacherID, branchID)
	return args.Bool(0), args.Error(1)
}

func (m *ScheduleRepository) ListLessons(ctx context.Context, from, to time.Time, teacherID, studentID, groupID *uuid.UUID, branchIDs []uuid.UUID) ([]domain.LessonDetails, error) {
	args := m.Called(ctx, from, to, teacherID, studentID, groupID, branchIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.LessonDetails), args.Error(1)
}

func (m *ScheduleRepository) GetTeacherSchedule(ctx context.Context, teacherID uuid.UUID, from, to time.Time) ([]domain.Lesson, error) {
	args := m.Called(ctx, teacherID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Lesson), args.Error(1)
}
