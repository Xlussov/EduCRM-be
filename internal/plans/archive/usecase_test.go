package archive

import (
	"context"
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

	plan := &domain.Plan{
		ID:       planID,
		BranchID: branchID,
	}

	tests := []struct {
		name        string
		caller      domain.Caller
		setupMocks  func(mockPR *mocks.SubscriptionRepository)
		expectedErr error
	}{
		{
			name:   "Success_SUPERADMIN",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			setupMocks: func(mockPR *mocks.SubscriptionRepository) {
				mockPR.On("GetPlanByID", mock.Anything, planID).Return(plan, nil).Once()
				mockPR.On("UpdatePlanStatus", mock.Anything, planID, domain.StatusArchived).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Success_ADMIN_HasAccess",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			setupMocks: func(mockPR *mocks.SubscriptionRepository) {
				mockPR.On("GetPlanByID", mock.Anything, planID).Return(plan, nil).Once()
				mockPR.On("UpdatePlanStatus", mock.Anything, planID, domain.StatusArchived).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Forbidden_ADMIN_NoAccess",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}},
			setupMocks: func(mockPR *mocks.SubscriptionRepository) {
				mockPR.On("GetPlanByID", mock.Anything, planID).Return(plan, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPR := new(mocks.SubscriptionRepository)
			if tt.setupMocks != nil {
				tt.setupMocks(mockPR)
			}

			uc := NewUseCase(mockPR)
			res, err := uc.Execute(context.Background(), tt.caller, planID, Request{Status: string(domain.StatusArchived)})

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, "success", res.Message)
			}

			mockPR.AssertExpectations(t)
		})
	}
}
