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
			req:  Request{BranchID: branch1},
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockGR.On("GetByBranchID", mock.Anything, branch1).Return([]*domain.GroupWithCount{
					{
						Group: domain.Group{
							ID:     uuid.New(),
							Name:   "Test Group",
							Status: domain.StatusActive,
						},
						StudentsCount: 5,
					},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Success_ADMIN_HasAccess",
			role: "ADMIN",
			req:  Request{BranchID: branch1},
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch1}, nil).Once()
				mockGR.On("GetByBranchID", mock.Anything, branch1).Return([]*domain.GroupWithCount{}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Forbidden_ADMIN_NoAccess",
			role: "ADMIN",
			req:  Request{BranchID: branch1},
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branch2}, nil).Once()
			},
			expectedErr: ErrBranchAccessDenied,
		},
		{
			name:        "Error_MissingBranchID",
			role:        "SUPERADMIN",
			req:         Request{BranchID: uuid.Nil},
			setupMocks:  func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {},
			expectedErr: ErrBranchIDRequired,
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
