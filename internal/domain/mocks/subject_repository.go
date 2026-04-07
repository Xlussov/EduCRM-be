package mocks

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SubjectRepository struct {
	mock.Mock
}

func (m *SubjectRepository) Create(ctx context.Context, subject *domain.Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}

func (m *SubjectRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.EntityStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *SubjectRepository) GetAll(ctx context.Context, branchID uuid.UUID, status *domain.EntityStatus) ([]*domain.Subject, error) {
	args := m.Called(ctx, branchID, status)
	if args.Get(0) != nil {
		return args.Get(0).([]*domain.Subject), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SubjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subject, error) {
	args := m.Called(ctx, id)
	var res *domain.Subject
	if args.Get(0) != nil {
		res = args.Get(0).(*domain.Subject)
	}
	return res, args.Error(1)
}

func (m *SubjectRepository) Update(ctx context.Context, subject *domain.Subject) (*domain.Subject, error) {
	args := m.Called(ctx, subject)
	var res *domain.Subject
	if args.Get(0) != nil {
		res = args.Get(0).(*domain.Subject)
	}
	return res, args.Error(1)
}
