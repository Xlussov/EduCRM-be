package mocks

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type BranchRepository struct {
	mock.Mock
}

func (m *BranchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	args := m.Called(ctx, branch)
	return args.Error(0)
}

func (m *BranchRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.EntityStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *BranchRepository) GetAll(ctx context.Context, status *domain.EntityStatus) ([]*domain.Branch, error) {
	args := m.Called(ctx, status)
	var res []*domain.Branch
	if args.Get(0) != nil {
		res = args.Get(0).([]*domain.Branch)
	}
	return res, args.Error(1)
}

func (m *BranchRepository) GetByUserID(ctx context.Context, userID uuid.UUID, status *domain.EntityStatus) ([]*domain.Branch, error) {
	args := m.Called(ctx, userID, status)
	var res []*domain.Branch
	if args.Get(0) != nil {
		res = args.Get(0).([]*domain.Branch)
	}
	return res, args.Error(1)
}

func (m *BranchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error) {
	args := m.Called(ctx, id)
	var res *domain.Branch
	if args.Get(0) != nil {
		res = args.Get(0).(*domain.Branch)
	}
	return res, args.Error(1)
}

func (m *BranchRepository) Update(ctx context.Context, branch *domain.Branch) (*domain.Branch, error) {
	args := m.Called(ctx, branch)
	var res *domain.Branch
	if args.Get(0) != nil {
		res = args.Get(0).(*domain.Branch)
	}
	return res, args.Error(1)
}

func (m *BranchRepository) IsActive(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *BranchRepository) CountActiveByIDs(ctx context.Context, ids []uuid.UUID) (int, error) {
	args := m.Called(ctx, ids)
	return args.Int(0), args.Error(1)
}
