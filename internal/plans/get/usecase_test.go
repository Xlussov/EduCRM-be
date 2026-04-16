package get

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
	subjectID := uuid.New()
	gridID := uuid.New()

	planDetails := &domain.PlanDetails{
		Plan: domain.Plan{
			ID:       planID,
			BranchID: branchID,
			Name:     "Plan A",
			Type:     domain.PlanTypeIndividual,
			Status:   domain.StatusActive,
		},
		Subjects: []*domain.SubjectBase{{ID: subjectID, Name: "Math"}},
		PricingGrid: []*domain.PricingGrid{{
			ID:              gridID,
			PlanID:          planID,
			LessonsPerMonth: 4,
			PricePerLesson:  100,
		}},
	}

	errDB := errors.New("db error")

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
				repo.On("GetPlanDetailsByID", mock.Anything, planID).Return(planDetails, nil).Once()
			},
		},
		{
			name:   "admin_success",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(repo *mocks.SubscriptionRepository) {
				repo.On("GetPlanDetailsByID", mock.Anything, planID).Return(planDetails, nil).Once()
			},
		},
		{
			name:   "admin_no_access",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}},
			mockSetup: func(repo *mocks.SubscriptionRepository) {
				repo.On("GetPlanDetailsByID", mock.Anything, planID).Return(planDetails, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "repo_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.SubscriptionRepository) {
				repo.On("GetPlanDetailsByID", mock.Anything, planID).Return((*domain.PlanDetails)(nil), errDB).Once()
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
				require.Equal(t, planID, res.ID)
				require.Len(t, res.Subjects, 1)
				require.Len(t, res.PricingGrid, 1)
			}

			repo.AssertExpectations(t)
		})
	}
}
