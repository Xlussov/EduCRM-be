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

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branch1}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branch2}}

	tests := []struct {
		name        string
		caller      domain.Caller
		req         Request
		setupMocks  func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository)
		expectedErr error
		assertMocks func(t *testing.T, mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository)
	}{
		{
			name:   "Success_SUPERADMIN",
			caller: callerSuper,
			req:    Request{BranchID: branch1, Name: "A1"},
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockUR.On("IsBranchActive", mock.Anything, branch1).Return(true, nil).Once()
				mockGR.On("Create", mock.Anything, mock.AnythingOfType("*domain.Group")).Run(func(args mock.Arguments) {
					group := args.Get(1).(*domain.Group)
					group.ID = uuid.New()
				}).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:   "Success_ADMIN_HasAccess",
			caller: callerAdmin,
			req:    Request{BranchID: branch1, Name: "A1"},
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockUR.On("IsBranchActive", mock.Anything, branch1).Return(true, nil).Once()
				mockGR.On("Create", mock.Anything, mock.AnythingOfType("*domain.Group")).Run(func(args mock.Arguments) {
					group := args.Get(1).(*domain.Group)
					group.ID = uuid.New()
				}).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:        "Forbidden_ADMIN_NoAccess",
			caller:      callerAdminNoAccess,
			req:         Request{BranchID: branch1, Name: "A1"},
			expectedErr: domain.ErrBranchAccessDenied,
			assertMocks: func(t *testing.T, mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockUR.AssertNotCalled(t, "IsBranchActive", mock.Anything, mock.Anything)
				mockGR.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUR := new(mocks.UserRepository)
			mockGR := new(mocks.GroupRepository)
			if tt.setupMocks != nil {
				tt.setupMocks(mockUR, mockGR)
			}

			uc := NewUseCase(mockGR, mockUR)
			_, err := uc.Execute(context.Background(), tt.caller, tt.req)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			if tt.assertMocks != nil {
				tt.assertMocks(t, mockUR, mockGR)
				return
			}

			mockUR.AssertExpectations(t)
			mockGR.AssertExpectations(t)
		})
	}
}
