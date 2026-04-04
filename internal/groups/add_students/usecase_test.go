package addstudents

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
	groupID := uuid.New()
	branchID := uuid.New()
	studentID := uuid.New()
	userID := uuid.New()
	req := Request{StudentIDs: []uuid.UUID{studentID}}

	tests := []struct {
		name        string
		role        string
		req         Request
		setupMocks  func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository, mockTx *mocks.MockTxManager)
		expectedErr error
	}{
		{
			name: "Success_SUPERADMIN",
			role: "SUPERADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository, mockTx *mocks.MockTxManager) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockSR.On("GetBranchID", mock.Anything, studentID).Return(branchID, nil).Once()
				mockGR.On("AddStudent", mock.Anything, groupID, studentID, mock.AnythingOfType("time.Time")).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Success_ADMIN_HasAccess",
			role: "ADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository, mockTx *mocks.MockTxManager) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branchID}, nil).Once()
				mockSR.On("GetBranchID", mock.Anything, studentID).Return(branchID, nil).Once()
				mockGR.On("AddStudent", mock.Anything, groupID, studentID, mock.AnythingOfType("time.Time")).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Forbidden_ADMIN_NoAccess",
			role: "ADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository, mockTx *mocks.MockTxManager) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockUR.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{uuid.New()}, nil).Once()
			},
			expectedErr: ErrBranchAccessDenied,
		},
		{
			name: "Error_StudentBranchMismatch",
			role: "SUPERADMIN",
			req:  req,
			setupMocks: func(mockUR *mocks.UserRepository, mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository, mockTx *mocks.MockTxManager) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockSR.On("GetBranchID", mock.Anything, studentID).Return(uuid.New(), nil).Once()
			},
			expectedErr: ErrStudentBranchMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUR := new(mocks.UserRepository)
			mockGR := new(mocks.GroupRepository)
			mockSR := new(mocks.StudentRepository)
			mockTx := new(mocks.MockTxManager)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUR, mockGR, mockSR, mockTx)
			}

			uc := NewUseCase(mockGR, mockUR, mockSR, mockTx)
			res, err := uc.Execute(context.Background(), userID, tt.role, groupID, tt.req)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, "success", res.Message)
			}

			mockUR.AssertExpectations(t)
			mockGR.AssertExpectations(t)
			mockSR.AssertExpectations(t)
		})
	}
}
