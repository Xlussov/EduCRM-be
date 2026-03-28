package mocks

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	args := m.Called(ctx, phone)
	var r0 *domain.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.User)
	}
	return r0, args.Error(1)
}

func (m *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	var r0 *domain.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.User)
	}
	return r0, args.Error(1)
}

func (m *UserRepository) AssignToBranches(ctx context.Context, userID uuid.UUID, branchIDs []uuid.UUID) error {
	args := m.Called(ctx, userID, branchIDs)
	return args.Error(0)
}

func (m *UserRepository) GetUserBranchIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}
