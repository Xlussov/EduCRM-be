package archive

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

	tests := []struct {
		name        string
		role        string
		setupMocks  func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository)
		expectedErr error
	}{
		{
			name: "Success_SUPERADMIN",
			role: "SUPERADMIN",
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Status: domain.StatusActive}, nil).Once()
				mockGR.On("UpdateStatus", mock.Anything, groupID, domain.StatusArchived).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Success_ADMIN_HasAccess",
			role: "ADMIN",
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Status: domain.StatusActive}, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branchID}, nil).Once()
				mockGR.On("UpdateStatus", mock.Anything, groupID, domain.StatusArchived).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Error_AlreadyArchived",
			role: "SUPERADMIN",
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Status: domain.StatusArchived}, nil).Once()
			},
			expectedErr: domain.ErrAlreadyArchived,
		},
		{
			name: "Forbidden_ADMIN_NoAccess",
			role: "ADMIN",
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Status: domain.StatusActive}, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{uuid.New()}, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name: "Error_GroupNotFound",
			role: "SUPERADMIN",
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return((*domain.Group)(nil), errors.New("no rows in result set")).Once()
			},
			expectedErr: errors.New("no rows in result set"),
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
			_, err := uc.Execute(context.Background(), userID, tt.role, groupID)

			if tt.expectedErr != nil {
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			mockUR.AssertExpectations(t)
			mockGR.AssertExpectations(t)
		})
	}
}
