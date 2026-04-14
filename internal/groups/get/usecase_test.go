package get

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
		groupID     uuid.UUID
		setupMocks  func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository)
		expectedErr error
	}{
		{
			name:    "Success_SUPERADMIN",
			role:    "SUPERADMIN",
			groupID: groupID,
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Name: "A1"}, nil).Once()
				mockGR.On("GetStudents", mock.Anything, groupID).Return([]*domain.GroupStudent{
					{ID: uuid.New(), FirstName: "John", LastName: "Doe", Status: domain.StatusActive},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:    "Success_ADMIN_HasAccess",
			role:    "ADMIN",
			groupID: groupID,
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Name: "A1"}, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branchID}, nil).Once()
				mockGR.On("GetStudents", mock.Anything, groupID).Return([]*domain.GroupStudent{}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:    "Forbidden_ADMIN_NoAccess",
			role:    "ADMIN",
			groupID: groupID,
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Name: "A1"}, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{uuid.New()}, nil).Once()
			},
			expectedErr: ErrBranchAccessDenied,
		},
		{
			name:    "Error_GroupNotFound",
			role:    "SUPERADMIN",
			groupID: groupID,
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
			_, err := uc.Execute(context.Background(), userID, tt.role, tt.groupID)

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
