package update

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
	req := Request{Name: "New Name"}
	otherBranchID := uuid.New()

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}}

	tests := []struct {
		name        string
		caller      domain.Caller
		req         Request
		setupMocks  func(mockGR *mocks.GroupRepository)
		expectedErr error
	}{
		{
			name:   "Success_SUPERADMIN",
			caller: callerSuper,
			req:    req,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockGR.On("UpdateName", mock.Anything, groupID, req.Name).Return(&domain.Group{ID: groupID, BranchID: branchID, Name: req.Name}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Success_ADMIN_HasAccess",
			caller: callerAdmin,
			req:    req,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockGR.On("UpdateName", mock.Anything, groupID, req.Name).Return(&domain.Group{ID: groupID, BranchID: branchID, Name: req.Name}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Forbidden_ADMIN_NoAccess",
			caller: callerAdminNoAccess,
			req:    req,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "Error_GroupNotFound",
			caller: callerSuper,
			req:    req,
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
			_, err := uc.Execute(context.Background(), tt.caller, groupID, tt.req)

			if tt.expectedErr != nil {
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			mockGR.AssertExpectations(t)
		})
	}
}
