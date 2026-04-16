package list

import (
	"context"
	"testing"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_Execute(t *testing.T) {
	branch1 := uuid.New()
	branch2 := uuid.New()
	studentID := uuid.New()
	subID := uuid.New()
	planID := uuid.New()
	subjectID := uuid.New()

	subs := []*domain.StudentSubscriptionDetails{
		{
			ID:        subID,
			StudentID: studentID,
			Plan:      domain.SubPlanDetails{ID: planID, Name: "Plan A"},
			Subject:   domain.SubSubjectDetails{ID: subjectID, Name: "Math"},
			StartDate: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name        string
		caller      domain.Caller
		setupMocks  func(mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository)
		expectedErr error
	}{
		{
			name:   "Success_SUPERADMIN",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			setupMocks: func(mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockSR.On("GetStudentSubscriptions", mock.Anything, studentID).Return(subs, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Success_ADMIN_HasAccess",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branch1}},
			setupMocks: func(mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockStdR.On("GetBranchID", mock.Anything, studentID).Return(branch1, nil).Once()
				mockSR.On("GetStudentSubscriptions", mock.Anything, studentID).Return(subs, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Forbidden_ADMIN_NoAccess",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branch2}},
			setupMocks: func(mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockStdR.On("GetBranchID", mock.Anything, studentID).Return(branch1, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSR := new(mocks.SubscriptionRepository)
			mockStdR := new(mocks.StudentRepository)

			if tt.setupMocks != nil {
				tt.setupMocks(mockSR, mockStdR)
			}

			uc := NewUseCase(mockSR, mockStdR)
			res, err := uc.Execute(context.Background(), tt.caller, studentID)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Len(t, res, 1)
				require.Equal(t, subID, res[0].ID)
				require.Equal(t, "Plan A", res[0].Plan.Name)
			}

			mockSR.AssertExpectations(t)
			mockStdR.AssertExpectations(t)
		})
	}
}
