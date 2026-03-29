package archive

import (
	"context"
	"errors"
	"testing"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Execute(t *testing.T) {
	branchID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.BranchRepository)
		repo.On("UpdateStatus", mock.Anything, branchID, domain.StatusArchived).Return(nil)

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), branchID)

		assert.NoError(t, err)
		assert.Equal(t, "success", res.Message)
		repo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		repo := new(mocks.BranchRepository)
		repo.On("UpdateStatus", mock.Anything, branchID, domain.StatusArchived).Return(errors.New("db error"))

		uc := NewUseCase(repo)
		_, err := uc.Execute(context.Background(), branchID)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		repo.AssertExpectations(t)
	})
}
