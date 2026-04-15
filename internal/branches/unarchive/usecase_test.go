package unarchive

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
	errDB := errors.New("db error")
	callerSuper := domain.Caller{UserID: uuid.New(), Role: domain.RoleSuperadmin}
	callerAdminAllowed := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerAdminDenied := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}}

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.BranchRepository)
		expectedError error
		expectedMsg   string
		assertRepo    func(t *testing.T, repo *mocks.BranchRepository)
	}{
		{
			name:   "superadmin success",
			caller: callerSuper,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return(&domain.Branch{
					ID:     branchID,
					Status: domain.StatusArchived,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, branchID, domain.StatusActive).Return(nil).Once()
			},
			expectedMsg: "success",
		},
		{
			name:          "admin access denied",
			caller:        callerAdminDenied,
			expectedError: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, repo *mocks.BranchRepository) {
				repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "admin access allowed",
			caller: callerAdminAllowed,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return(&domain.Branch{
					ID:     branchID,
					Status: domain.StatusArchived,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, branchID, domain.StatusActive).Return(nil).Once()
			},
			expectedMsg: "success",
		},
		{
			name:   "already active",
			caller: callerSuper,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return(&domain.Branch{
					ID:     branchID,
					Status: domain.StatusActive,
				}, nil).Once()
			},
			expectedError: domain.ErrAlreadyActive,
		},
		{
			name:   "db error on get",
			caller: callerSuper,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return((*domain.Branch)(nil), errDB).Once()
			},
			expectedError: errDB,
		},
		{
			name:   "db error on update",
			caller: callerSuper,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return(&domain.Branch{
					ID:     branchID,
					Status: domain.StatusArchived,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, branchID, domain.StatusActive).Return(errDB).Once()
			},
			expectedError: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.BranchRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, branchID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, res.Message)
			}

			if tt.assertRepo != nil {
				tt.assertRepo(t, repo)
				return
			}

			repo.AssertExpectations(t)
		})
	}
}
