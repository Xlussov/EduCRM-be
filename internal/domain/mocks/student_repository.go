package mocks

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type StudentRepository struct {
	mock.Mock
}

func (m *StudentRepository) Create(ctx context.Context, student *domain.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *StudentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.EntityStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *StudentRepository) GetBranchID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *StudentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Student, error) {
	args := m.Called(ctx, id)
	var res *domain.Student
	if args.Get(0) != nil {
		res = args.Get(0).(*domain.Student)
	}
	return res, args.Error(1)
}

func (m *StudentRepository) Update(ctx context.Context, student *domain.Student) (*domain.Student, error) {
	args := m.Called(ctx, student)
	var res *domain.Student
	if args.Get(0) != nil {
		res = args.Get(0).(*domain.Student)
	}
	return res, args.Error(1)
}

func (m *StudentRepository) GetByBranchID(ctx context.Context, branchID uuid.UUID) ([]*domain.Student, error) {
	args := m.Called(ctx, branchID)
	var res []*domain.Student
	if args.Get(0) != nil {
		res = args.Get(0).([]*domain.Student)
	}
	return res, args.Error(1)
}
