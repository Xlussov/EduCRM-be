package get

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

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.BranchRepository)
		repo.On("GetByID", mock.Anything, branchID).Return(branch, nil)

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), branchID)

		assert.NoError(t, err)
		assert.Equal(t, branchID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("db err", func(t *testing.T) {
		repo := new(mocks.BranchRepository)
		repo.On("GetByID", mock.Anything, branchID).Return(nil, errors.New("db err"))

		uc := NewUseCase(repo)
		_, err := uc.Execute(context.Background(), branchID)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
