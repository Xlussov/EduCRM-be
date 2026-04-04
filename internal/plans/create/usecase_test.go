package create

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
	branch1 := uuid.New()
	branch2 := uuid.New()
	userID := uuid.New()

	req := Request{
		BranchID:   branch1,
		Name:       "Test Plan",
		Type:       "INDIVIDUAL",
		SubjectIDs: []uuid.UUID{uuid.New()},
		PricingGrid: []PricingGridItem{
			{Lessons: 4, Price: 100},
		},
	}

	tests := []struct {
		name        string
		role        string
		req         Request
		setupMocks  func(mockUR *mocks.UserRepository, mockPR *mocks.SubscriptionRepository, mockTx *mocks.MockTxManager)
		expectedErr error
	}{
		{
			name: "Success_SUPERADMIN",
			role: "SUPERADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockPR *mocks.SubscriptionRepository, mockTx *mocks.MockTxManager) {
				mockPR.On("CreatePlan", mock.Anything, mock.AnythingOfType("*domain.Plan"), mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					plan := args.Get(1).(*domain.Plan)
					plan.ID = uuid.New()
				}).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Success_ADMIN_HasAccess",
			role: "ADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockPR *mocks.SubscriptionRepository, mockTx *mocks.MockTxManager) {
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch1}, nil).Once()
				mockPR.On("CreatePlan", mock.Anything, mock.AnythingOfType("*domain.Plan"), mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					plan := args.Get(1).(*domain.Plan)
					plan.ID = uuid.New()
				}).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Forbidden_ADMIN_NoAccess",
			role: "ADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockPR *mocks.SubscriptionRepository, mockTx *mocks.MockTxManager) {
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch2}, nil).Once()
			},
			expectedErr: ErrBranchAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUR := new(mocks.UserRepository)
			mockPR := new(mocks.SubscriptionRepository)
			mockTx := new(mocks.MockTxManager)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUR, mockPR, mockTx)
			}

			uc := NewUseCase(mockTx, mockPR, mockUR)
			_, err := uc.Execute(context.Background(), userID, tt.role, tt.req)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			mockUR.AssertExpectations(t)
			mockPR.AssertExpectations(t)
		})
	}
}
