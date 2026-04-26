package mocks

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) Create(ctx context.Context, student *domain.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *MockStudentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.EntityStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockStudentRepository) GetBranchID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockStudentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Student, error) {
	args := m.Called(ctx, id)
	var student *domain.Student
	if args.Get(0) != nil {
		student = args.Get(0).(*domain.Student)
	}
	return student, args.Error(1)
}

func (m *MockStudentRepository) Update(ctx context.Context, student *domain.Student) (*domain.Student, error) {
	args := m.Called(ctx, student)
	var updated *domain.Student
	if args.Get(0) != nil {
		updated = args.Get(0).(*domain.Student)
	}
	return updated, args.Error(1)
}

func (m *MockStudentRepository) GetByBranchID(ctx context.Context, branchID uuid.UUID, status *domain.EntityStatus) ([]*domain.Student, error) {
	args := m.Called(ctx, branchID, status)
	var students []*domain.Student
	if args.Get(0) != nil {
		students = args.Get(0).([]*domain.Student)
	}
	return students, args.Error(1)
}

func (m *MockStudentRepository) GetByBranchIDAndTeacherID(ctx context.Context, branchID, teacherID uuid.UUID, status *domain.EntityStatus) ([]*domain.Student, error) {
	args := m.Called(ctx, branchID, teacherID, status)
	var students []*domain.Student
	if args.Get(0) != nil {
		students = args.Get(0).([]*domain.Student)
	}
	return students, args.Error(1)
}

func (m *MockStudentRepository) IsTeacherStudent(ctx context.Context, teacherID, studentID uuid.UUID) (bool, error) {
	args := m.Called(ctx, teacherID, studentID)
	return args.Bool(0), args.Error(1)
}
