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
	otherBranchID := uuid.New()

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}}
	callerTeacher := domain.Caller{UserID: userID, Role: domain.RoleTeacher, BranchIDs: []uuid.UUID{branchID}}

	tests := []struct {
		name        string
		caller      domain.Caller
		groupID     uuid.UUID
		setupMocks  func(mockGR *mocks.GroupRepository)
		expectedErr error
	}{
		{
			name:    "Success_SUPERADMIN",
			caller:  callerSuper,
			groupID: groupID,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Name: "A1"}, nil).Once()
				mockGR.On("GetStudents", mock.Anything, groupID).Return([]*domain.GroupStudent{
					{ID: uuid.New(), FirstName: "John", LastName: "Doe", Status: domain.StatusActive},
				}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:    "Success_ADMIN_HasAccess",
			caller:  callerAdmin,
			groupID: groupID,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Name: "A1"}, nil).Once()
				mockGR.On("GetStudents", mock.Anything, groupID).Return([]*domain.GroupStudent{}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:    "Success_TEACHER_HasAccess",
			caller:  callerTeacher,
			groupID: groupID,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Name: "A1"}, nil).Once()
				mockGR.On("IsTeacherGroup", mock.Anything, userID, groupID).Return(true, nil).Once()
				mockGR.On("GetStudents", mock.Anything, groupID).Return([]*domain.GroupStudent{}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:    "Forbidden_TEACHER_NotAssigned",
			caller:  callerTeacher,
			groupID: groupID,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Name: "A1"}, nil).Once()
				mockGR.On("IsTeacherGroup", mock.Anything, userID, groupID).Return(false, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:    "Forbidden_ADMIN_NoAccess",
			caller:  callerAdminNoAccess,
			groupID: groupID,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID, Name: "A1"}, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:    "Error_GroupNotFound",
			caller:  callerSuper,
			groupID: groupID,
			setupMocks: func(mockGR *mocks.GroupRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return((*domain.Group)(nil), errors.New("no rows in result set")).Once()
			},
			expectedErr: ErrGroupNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGR := new(mocks.GroupRepository)
			if tt.setupMocks != nil {
				tt.setupMocks(mockGR)
			}

			uc := NewUseCase(mockGR)
			_, err := uc.Execute(context.Background(), tt.caller, tt.groupID)

			if tt.expectedErr != nil {
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			mockGR.AssertExpectations(t)
		})
	}
}
