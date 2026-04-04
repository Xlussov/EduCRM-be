package create

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

	req := Request{
		PlanID:    uuid.New(),
		SubjectID: uuid.New(),
		StartDate: time.Now(),
	}

	tests := []struct {
		name        string
		role        string
		req         Request
		setupMocks  func(mockUR *mocks.UserRepository, mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository)
		expectedErr error
	}{
		{
			name: "Success_SUPERADMIN",
			role: "SUPERADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockSR.On("ValidatePlanSubject", mock.Anything, req.PlanID, req.SubjectID).Return(true, nil).Once()
				mockSR.On("AssignToStudent", mock.Anything, mock.AnythingOfType("*domain.StudentSubscription")).Run(func(args mock.Arguments) {
					sub := args.Get(1).(*domain.StudentSubscription)
					sub.ID = uuid.New()
				}).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Success_ADMIN_HasAccess",
			role: "ADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockStdR.On("GetBranchID", mock.Anything, studentID).Return(branch1, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch1}, nil).Once()
				mockSR.On("ValidatePlanSubject", mock.Anything, req.PlanID, req.SubjectID).Return(true, nil).Once()
				mockSR.On("AssignToStudent", mock.Anything, mock.AnythingOfType("*domain.StudentSubscription")).Run(func(args mock.Arguments) {
					sub := args.Get(1).(*domain.StudentSubscription)
					sub.ID = uuid.New()
				}).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Forbidden_ADMIN_NoAccess",
			role: "ADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockStdR.On("GetBranchID", mock.Anything, studentID).Return(branch1, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch2}, nil).Once()
			},
			expectedErr: ErrBranchAccessDenied,
		},
		{
			name: "Error_InvalidSubject",
			role: "SUPERADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockSR.On("ValidatePlanSubject", mock.Anything, req.PlanID, req.SubjectID).Return(false, nil).Once()
			},
			expectedErr: ErrInvalidSubject,
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
			_, err := uc.Execute(context.Background(), userID, studentID, tt.role, tt.req)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			mockUR.AssertExpectations(t)
			mockSR.AssertExpectations(t)
			mockStdR.AssertExpectations(t)
		})
	}
}
