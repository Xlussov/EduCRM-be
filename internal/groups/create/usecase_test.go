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

	tests := []struct {
		name        string
		role        string
		req         Request
		setupMocks  func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository)
		expectedErr error
	}{
		{
			name: "Success_SUPERADMIN",
			role: "SUPERADMIN",
			req:  Request{BranchID: branch1, Name: "A1"},
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockGR.On("Create", mock.Anything, mock.AnythingOfType("*domain.Group")).Run(func(args mock.Arguments) {
					group := args.Get(1).(*domain.Group)
					group.ID = uuid.New()
				}).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Success_ADMIN_HasAccess",
			role: "ADMIN",
			req:  Request{BranchID: branch1, Name: "A1"},
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch1}, nil).Once()
				mockGR.On("Create", mock.Anything, mock.AnythingOfType("*domain.Group")).Run(func(args mock.Arguments) {
					group := args.Get(1).(*domain.Group)
					group.ID = uuid.New()
				}).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Forbidden_ADMIN_NoAccess",
			role: "ADMIN",
			req:  Request{BranchID: branch1, Name: "A1"},
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch2}, nil).Once()
			},
			expectedErr: ErrBranchAccessDenied,
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
			_, err := uc.Execute(context.Background(), userID, tt.role, tt.req)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			mockUR.AssertExpectations(t)
			mockGR.AssertExpectations(t)
		})
	}
}
