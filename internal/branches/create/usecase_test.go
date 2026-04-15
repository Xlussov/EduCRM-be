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
	userID := uuid.New()
	branchId := uuid.New()
	caller := domain.Caller{UserID: userID, Role: domain.RoleAdmin}
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
	errCreate := errors.New("db error")
	errAssign := errors.New("db user error")

	tests := []struct {
		name        string
		caller      domain.Caller
		mockSetup   func(branchRepo *mocks.BranchRepository, userRepo *mocks.UserRepository)
		expectedID  uuid.UUID
		expectedErr error
		assertRepo  func(t *testing.T, branchRepo *mocks.BranchRepository, userRepo *mocks.UserRepository)
	}{
		{
			name:   "success",
			caller: caller,
			mockSetup: func(branchRepo *mocks.BranchRepository, userRepo *mocks.UserRepository) {
				branchRepo.On("Create", mock.Anything, mock.MatchedBy(func(b *domain.Branch) bool {
					b.ID = branchId
					b.CreatedAt = time.Now()
					b.UpdatedAt = time.Now()
					return b.Name == expectedBranch.Name && b.Address == expectedBranch.Address && b.Status == expectedBranch.Status
				})).Return(nil)

				userRepo.On("AssignToBranches", mock.Anything, userID, []uuid.UUID{branchId}).Return(nil)
			},
			expectedID: branchId,
		},
		{
			name:   "failed to create branch",
			caller: caller,
			mockSetup: func(branchRepo *mocks.BranchRepository, userRepo *mocks.UserRepository) {
				branchRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Branch")).Return(errCreate)
			},
			expectedErr: errCreate,
			assertRepo: func(t *testing.T, branchRepo *mocks.BranchRepository, userRepo *mocks.UserRepository) {
				branchRepo.AssertExpectations(t)
				userRepo.AssertNotCalled(t, "AssignToBranches", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name:   "failed to assign user to branch",
			caller: caller,
			mockSetup: func(branchRepo *mocks.BranchRepository, userRepo *mocks.UserRepository) {
				branchRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Branch")).Return(nil).Run(func(args mock.Arguments) {
					b := args.Get(1).(*domain.Branch)
					b.ID = branchId
				})

				userRepo.On("AssignToBranches", mock.Anything, userID, []uuid.UUID{branchId}).Return(errAssign)
			},
			expectedErr: errAssign,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branchRepo := new(mocks.BranchRepository)
			userRepo := new(mocks.UserRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(branchRepo, userRepo)
			}

			uc := NewUseCase(branchRepo, userRepo, &mocks.MockTxManager{})
			res, err := uc.Execute(context.Background(), tt.caller, req)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, res.ID)
			}

			if tt.assertRepo != nil {
				tt.assertRepo(t, branchRepo, userRepo)
				return
			}

			branchRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}
