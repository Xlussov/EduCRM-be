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
	callerSuper := domain.Caller{UserID: uuid.New(), Role: domain.RoleSuperadmin}
	callerAdminAllowed := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerAdminDenied := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}}
	errDB := errors.New("db err")
	branch := &domain.Branch{
		ID:        branchID,
		Name:      "Test",
		Address:   "Test",
		City:      "Test City",
		Status:    domain.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name        string
		caller      domain.Caller
		mockSetup   func(repo *mocks.BranchRepository)
		expectedErr error
		expectedID  uuid.UUID
		assertRepo  func(t *testing.T, repo *mocks.BranchRepository)
	}{
		{
			name:   "success",
			caller: callerSuper,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return(branch, nil)
			},
			expectedID: branchID,
		},
		{
			name:        "admin access denied",
			caller:      callerAdminDenied,
			expectedErr: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, repo *mocks.BranchRepository) {
				repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "admin access allowed",
			caller: callerAdminAllowed,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return(branch, nil)
			},
			expectedID: branchID,
		},
		{
			name:   "db err",
			caller: callerSuper,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return(nil, errDB)
			},
			expectedErr: errDB,
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

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, res.ID)
			}

			if tt.assertRepo != nil {
				tt.assertRepo(t, repo)
				return
			}

			repo.AssertExpectations(t)
		})
	}
}
