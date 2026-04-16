package unarchive

import (
	"context"
	"errors"
	"testing"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_Execute(t *testing.T) {
	planID := uuid.New()
	branchID := uuid.New()
	otherBranchID := uuid.New()
	errDB := errors.New("db error")

	archivedPlan := &domain.Plan{ID: planID, BranchID: branchID, Status: domain.StatusArchived}
	activePlan := &domain.Plan{ID: planID, BranchID: branchID, Status: domain.StatusActive}

	tests := []struct {
		name        string
		caller      domain.Caller
		mockSetup   func(repo *mocks.SubscriptionRepository)
		expectedErr error
	}{
		{
			name:   "superadmin_success",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.SubscriptionRepository) {
				repo.On("GetPlanByID", mock.Anything, planID).Return(archivedPlan, nil).Once()
				repo.On("UpdatePlanStatus", mock.Anything, planID, domain.StatusActive).Return(nil).Once()
			},
		},
		{
			name:   "admin_success",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(repo *mocks.SubscriptionRepository) {
				repo.On("GetPlanByID", mock.Anything, planID).Return(archivedPlan, nil).Once()
				repo.On("UpdatePlanStatus", mock.Anything, planID, domain.StatusActive).Return(nil).Once()
			},
		},
		{
			name:   "admin_no_access",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}},
			mockSetup: func(repo *mocks.SubscriptionRepository) {
				repo.On("GetPlanByID", mock.Anything, planID).Return(archivedPlan, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "already_active",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.SubscriptionRepository) {
				repo.On("GetPlanByID", mock.Anything, planID).Return(activePlan, nil).Once()
			},
			expectedErr: domain.ErrAlreadyActive,
		},
		{
			name:   "repo_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.SubscriptionRepository) {
				repo.On("GetPlanByID", mock.Anything, planID).Return((*domain.Plan)(nil), errDB).Once()
			},
			expectedErr: errDB,
		},
		{
			name:   "update_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.SubscriptionRepository) {
				repo.On("GetPlanByID", mock.Anything, planID).Return(archivedPlan, nil).Once()
				repo.On("UpdatePlanStatus", mock.Anything, planID, domain.StatusActive).Return(errDB).Once()
			},
			expectedErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.SubscriptionRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, planID)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, "success", res.Message)
			}

			repo.AssertExpectations(t)
		})
	}
}
