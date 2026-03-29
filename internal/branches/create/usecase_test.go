package create

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
	userId := uuid.New()
	branchId := uuid.New()

	req := Request{
		Name:    "Test Branch",
		Address: "123 Main St",
		City:    "New York",
	}

	expectedBranch := &domain.Branch{
		Name:    req.Name,
		Address: req.Address,
		City:    req.City,
		Status:  domain.StatusActive,
	}

	t.Run("success", func(t *testing.T) {
		branchRepo := new(mocks.BranchRepository)
		userRepo := new(mocks.UserRepository)

		branchRepo.On("Create", mock.Anything, mock.MatchedBy(func(b *domain.Branch) bool {
			// Populate ID to mock DB behavior
			b.ID = branchId
			b.CreatedAt = time.Now()
			b.UpdatedAt = time.Now()
			return b.Name == expectedBranch.Name && b.Address == expectedBranch.Address && b.Status == expectedBranch.Status
		})).Return(nil)

		userRepo.On("AssignToBranches", mock.Anything, userId, []uuid.UUID{branchId}).Return(nil)

		uc := NewUseCase(branchRepo, userRepo)
		res, err := uc.Execute(context.Background(), userId, req)

		assert.NoError(t, err)
		assert.Equal(t, branchId, res.ID)

		branchRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("failed to create branch", func(t *testing.T) {
		branchRepo := new(mocks.BranchRepository)
		userRepo := new(mocks.UserRepository)

		branchRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Branch")).Return(errors.New("db error"))

		uc := NewUseCase(branchRepo, userRepo)
		_, err := uc.Execute(context.Background(), userId, req)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())

		branchRepo.AssertExpectations(t)
		userRepo.AssertNotCalled(t, "AssignToBranches", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("failed to assign user to branch", func(t *testing.T) {
		branchRepo := new(mocks.BranchRepository)
		userRepo := new(mocks.UserRepository)

		branchRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Branch")).Return(nil).Run(func(args mock.Arguments) {
			b := args.Get(1).(*domain.Branch)
			b.ID = branchId
		})

		userRepo.On("AssignToBranches", mock.Anything, userId, []uuid.UUID{branchId}).Return(errors.New("db user error"))

		uc := NewUseCase(branchRepo, userRepo)
		_, err := uc.Execute(context.Background(), userId, req)

		assert.Error(t, err)
		assert.Equal(t, "db user error", err.Error())

		branchRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})
}
