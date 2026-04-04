package list

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
	planID := uuid.New()

	req := Request{
		BranchID: branch1,
	}

	planDetails := []*domain.PlanDetails{
		{
			Plan: domain.Plan{
				ID:       planID,
				BranchID: branch1,
				Name:     "Test Plan",
				Type:     "INDIVIDUAL",
				Status:   "ACTIVE",
			},
			Subjects: []*domain.SubjectBase{
				{
					ID:   uuid.New(),
					Name: "Math",
				},
			},
			PricingGrid: []*domain.PricingGrid{
				{
					ID:              uuid.New(),
					PlanID:          planID,
					LessonsPerMonth: 4,
					PricePerLesson:  100.0,
				},
			},
		},
	}

	tests := []struct {
		name        string
		role        string
		req         Request
		setupMocks  func(mockUR *mocks.UserRepository, mockPR *mocks.SubscriptionRepository)
		expectedErr error
	}{
		{
			name: "Success_SUPERADMIN",
			role: "SUPERADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockPR *mocks.SubscriptionRepository) {
				mockPR.On("GetPlansByBranchID", mock.Anything, branch1).Return(planDetails, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Success_ADMIN_HasAccess",
			role: "ADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockPR *mocks.SubscriptionRepository) {
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch1}, nil).Once()
				mockPR.On("GetPlansByBranchID", mock.Anything, branch1).Return(planDetails, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Forbidden_ADMIN_NoAccess",
			role: "ADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockPR *mocks.SubscriptionRepository) {
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch2}, nil).Once()
			},
			expectedErr: ErrBranchAccessDenied,
		},
		{
			name: "Error_BranchIDRequired",
			role: "SUPERADMIN",
			req:  Request{},
			setupMocks: func(mockUR *mocks.UserRepository, mockPR *mocks.SubscriptionRepository) {
			},
			expectedErr: ErrBranchIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUR := new(mocks.UserRepository)
			mockPR := new(mocks.SubscriptionRepository)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUR, mockPR)
			}

			uc := NewUseCase(mockPR, mockUR)
			res, err := uc.Execute(context.Background(), userID, tt.role, tt.req)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Len(t, res, 1)
				require.Equal(t, planID, res[0].ID)
				require.Len(t, res[0].Subjects, 1)
				require.Len(t, res[0].PricingGrid, 1)
			}

			mockUR.AssertExpectations(t)
			mockPR.AssertExpectations(t)
		})
	}
}
