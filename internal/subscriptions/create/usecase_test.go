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
	studentID := uuid.New()

	req := Request{
		PlanID:    uuid.New(),
		SubjectID: uuid.New(),
		StartDate: time.Now(),
	}

	tests := []struct {
		name        string
		caller      domain.Caller
		req         Request
		setupMocks  func(mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository)
		expectedErr error
	}{
		{
			name:   "Success_SUPERADMIN",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req:    req,
			setupMocks: func(mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockSR.On("ValidatePlanSubject", mock.Anything, req.PlanID, req.SubjectID).Return(true, nil).Once()
				mockSR.On("GetSubscriptionBranchIDs", mock.Anything, studentID, req.PlanID, req.SubjectID).Return(&domain.SubscriptionBranchIDs{
					StudentBranchID: branch1,
					PlanBranchID:    branch1,
					SubjectBranchID: branch1,
				}, nil).Once()
				mockSR.On("AssignToStudent", mock.Anything, mock.AnythingOfType("*domain.StudentSubscription")).Run(func(args mock.Arguments) {
					sub := args.Get(1).(*domain.StudentSubscription)
					sub.ID = uuid.New()
				}).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Success_ADMIN_HasAccess",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branch1}},
			req:    req,
			setupMocks: func(mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockStdR.On("GetBranchID", mock.Anything, studentID).Return(branch1, nil).Once()
				mockSR.On("ValidatePlanSubject", mock.Anything, req.PlanID, req.SubjectID).Return(true, nil).Once()
				mockSR.On("GetSubscriptionBranchIDs", mock.Anything, studentID, req.PlanID, req.SubjectID).Return(&domain.SubscriptionBranchIDs{
					StudentBranchID: branch1,
					PlanBranchID:    branch1,
					SubjectBranchID: branch1,
				}, nil).Once()
				mockSR.On("AssignToStudent", mock.Anything, mock.AnythingOfType("*domain.StudentSubscription")).Run(func(args mock.Arguments) {
					sub := args.Get(1).(*domain.StudentSubscription)
					sub.ID = uuid.New()
				}).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Forbidden_ADMIN_NoAccess",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branch2}},
			req:    req,
			setupMocks: func(mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockStdR.On("GetBranchID", mock.Anything, studentID).Return(branch1, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "Error_ArchivedOrInvalidReference",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req:    req,
			setupMocks: func(mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockSR.On("ValidatePlanSubject", mock.Anything, req.PlanID, req.SubjectID).Return(false, nil).Once()
			},
			expectedErr: domain.ErrArchivedReference,
		},
		{
			name:   "Error_CrossBranchData",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req:    req,
			setupMocks: func(mockSR *mocks.SubscriptionRepository, mockStdR *mocks.StudentRepository) {
				mockSR.On("ValidatePlanSubject", mock.Anything, req.PlanID, req.SubjectID).Return(true, nil).Once()
				mockSR.On("GetSubscriptionBranchIDs", mock.Anything, studentID, req.PlanID, req.SubjectID).Return(&domain.SubscriptionBranchIDs{
					StudentBranchID: branch1,
					PlanBranchID:    branch1,
					SubjectBranchID: branch2,
				}, nil).Once()
			},
			expectedErr: ErrCrossBranchData,
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
			_, err := uc.Execute(context.Background(), tt.caller, studentID, tt.req)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			mockSR.AssertExpectations(t)
			mockStdR.AssertExpectations(t)
		})
	}
}
