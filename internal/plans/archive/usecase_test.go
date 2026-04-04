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
	userID := uuid.New()
	branchID := uuid.New()

	plan := &domain.Plan{
		ID:       planID,
		BranchID: branchID,
	}

	tests := []struct {
		name        string
		role        string
		setupMocks  func(mockPR *mocks.SubscriptionRepository, mockUR *mocks.UserRepository)
		expectedErr error
	}{
		{
			name: "Success_SUPERADMIN",
			role: "SUPERADMIN",
			setupMocks: func(mockPR *mocks.SubscriptionRepository, mockUR *mocks.UserRepository) {
				mockPR.On("GetPlanByID", mock.Anything, planID).Return(plan, nil).Once()
				mockPR.On("UpdatePlanStatus", mock.Anything, planID, domain.StatusArchived).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Success_ADMIN_HasAccess",
			role: "ADMIN",
			setupMocks: func(mockPR *mocks.SubscriptionRepository, mockUR *mocks.UserRepository) {
				mockPR.On("GetPlanByID", mock.Anything, planID).Return(plan, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branchID}, nil).Once()
				mockPR.On("UpdatePlanStatus", mock.Anything, planID, domain.StatusArchived).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Forbidden_ADMIN_NoAccess",
			role: "ADMIN",
			setupMocks: func(mockPR *mocks.SubscriptionRepository, mockUR *mocks.UserRepository) {
				mockPR.On("GetPlanByID", mock.Anything, planID).Return(plan, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{uuid.New()}, nil).Once()
			},
			expectedErr: ErrBranchAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPR := new(mocks.SubscriptionRepository)
			mockUR := new(mocks.UserRepository)
			if tt.setupMocks != nil {
				tt.setupMocks(mockPR, mockUR)
			}

			uc := NewUseCase(mockPR, mockUR)
			res, err := uc.Execute(context.Background(), userID, tt.role, planID, Request{Status: string(domain.StatusArchived)})

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, "success", res.Message)
			}

			mockPR.AssertExpectations(t)
			mockUR.AssertExpectations(t)
		})
	}
}
