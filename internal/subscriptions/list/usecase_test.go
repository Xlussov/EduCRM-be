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
	userID := uuid.New()
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
		role        string
		setupMocks  func(mockUR *mocks.UserRepository, mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository)
		expectedErr error
	}{
		{
			name: "Success_SUPERADMIN",
			role: "SUPERADMIN",
			setupMocks: func(mockUR *mocks.UserRepository, mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockSR.On("GetStudentSubscriptions", mock.Anything, studentID).Return(subs, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Success_ADMIN_HasAccess",
			role: "ADMIN",
			setupMocks: func(mockUR *mocks.UserRepository, mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockStdR.On("GetBranchID", mock.Anything, studentID).Return(branch1, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch1}, nil).Once()
				mockSR.On("GetStudentSubscriptions", mock.Anything, studentID).Return(subs, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Forbidden_ADMIN_NoAccess",
			role: "ADMIN",
			setupMocks: func(mockUR *mocks.UserRepository, mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockStdR.On("GetBranchID", mock.Anything, studentID).Return(branch1, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch2}, nil).Once()
			},
			expectedErr: ErrBranchAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUR := new(mocks.UserRepository)
			mockSR := new(mocks.SubscriptionRepository)
			mockStdR := new(mocks.StudentRepository)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUR, mockSR, mockStdR)
			}

			uc := NewUseCase(mockSR, mockUR, mockStdR)
			res, err := uc.Execute(context.Background(), userID, studentID, tt.role)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Len(t, res, 1)
				require.Equal(t, subID, res[0].ID)
				require.Equal(t, "Plan A", res[0].Plan.Name)
			}

			mockUR.AssertExpectations(t)
			mockSR.AssertExpectations(t)
			mockStdR.AssertExpectations(t)
		})
	}
}
