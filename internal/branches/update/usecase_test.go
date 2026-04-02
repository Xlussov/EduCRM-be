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
		updatedDomain := &domain.Branch{
			ID:      branchID,
			Name:    req.Name,
			Address: req.Address,
			City:    req.City,
			Status:  domain.StatusActive,
		}
		repo.On("Update", mock.Anything, &domain.Branch{
			ID:      branchID,
			Name:    req.Name,
			Address: req.Address,
			City:    req.City,
		}).Return(updatedDomain, nil)

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), branchID, req)

		assert.Equal(t, branchID.String(), res.ID)
		assert.Equal(t, req.Name, res.Name)
		assert.Equal(t, req.Address, res.Address)
		assert.Equal(t, req.City, res.City)
		assert.Equal(t, string(domain.StatusActive), res.Status)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		repo := new(mocks.BranchRepository)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("db err"))

		uc := NewUseCase(repo)
		_, err := uc.Execute(context.Background(), branchID, req)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
