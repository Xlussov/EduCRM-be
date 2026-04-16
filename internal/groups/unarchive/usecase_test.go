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
	groupID := uuid.New()
	branchID := uuid.New()
	userID := uuid.New()
	otherBranchID := uuid.New()

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}}

	tests := []struct {
		name        string
		caller      domain.Caller
		setupMocks  func(mockGR *mocks.GroupRepository)
		expectedErr error
	}{
		{
			name:   "Success_SUPERADMIN",
			caller: callerSuper,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Status: domain.StatusArchived}, nil).Once()
				mockGR.On("UpdateStatus", mock.Anything, groupID, domain.StatusActive).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Success_ADMIN_HasAccess",
			caller: callerAdmin,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Status: domain.StatusArchived}, nil).Once()
				mockGR.On("UpdateStatus", mock.Anything, groupID, domain.StatusActive).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Error_AlreadyActive",
			caller: callerSuper,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Status: domain.StatusActive}, nil).Once()
			},
			expectedErr: domain.ErrAlreadyActive,
		},
		{
			name:   "Forbidden_ADMIN_NoAccess",
			caller: callerAdminNoAccess,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Status: domain.StatusArchived}, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "Error_GroupNotFound",
			caller: callerSuper,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return((*domain.Group)(nil), errors.New("no rows in result set")).Once()
			},
			expectedErr: errors.New("no rows in result set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGR := new(mocks.GroupRepository)
			if tt.setupMocks != nil {
				tt.setupMocks(mockGR)
			}

			uc := NewUseCase(mockGR)
			_, err := uc.Execute(context.Background(), tt.caller, groupID)

			if tt.expectedErr != nil {
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			mockGR.AssertExpectations(t)
		})
	}
}
