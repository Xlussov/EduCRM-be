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
