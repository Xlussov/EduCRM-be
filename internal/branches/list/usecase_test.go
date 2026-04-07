package list

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Execute(t *testing.T) {
	userID := uuid.New()
	branchID := uuid.New()

	branch := &domain.Branch{
		ID:        branchID,
		Name:      "Test",
		Address:   "Test",
		City:      "Test City",
		Status:    domain.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("superadmin - gets all", func(t *testing.T) {
		repo := new(mocks.BranchRepository)
		repo.On("GetAll", mock.Anything, (*domain.EntityStatus)(nil)).Return([]*domain.Branch{branch}, nil)

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), userID, "SUPERADMIN", Request{})

		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, branchID, res[0].ID)
		repo.AssertExpectations(t)
	})

	t.Run("admin - gets by param", func(t *testing.T) {
		repo := new(mocks.BranchRepository)
		repo.On("GetByUserID", mock.Anything, userID, (*domain.EntityStatus)(nil)).Return([]*domain.Branch{branch}, nil)

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), userID, "ADMIN", Request{})

		assert.NoError(t, err)
		assert.Len(t, res, 1)
		repo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		repo := new(mocks.BranchRepository)
		repo.On("GetAll", mock.Anything, (*domain.EntityStatus)(nil)).Return(nil, errors.New("db err"))

		uc := NewUseCase(repo)
		_, err := uc.Execute(context.Background(), userID, "SUPERADMIN", Request{})

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
