package update

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
	req := Request{
		Name:    "Updated Name",
		Address: "Updated Address",
		City:    "Updated City",
	}

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.BranchRepository)
		repo.On("Update", mock.Anything, &domain.Branch{
			ID:      branchID,
			Name:    req.Name,
			Address: req.Address,
			City:    req.City,
		}).Return(nil)

		uc := NewUseCase(repo)
		err := uc.Execute(context.Background(), branchID, req)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		repo := new(mocks.BranchRepository)
		repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db err"))

		uc := NewUseCase(repo)
		err := uc.Execute(context.Background(), branchID, req)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
