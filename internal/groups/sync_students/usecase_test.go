package syncstudents

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
	studentA := uuid.New()
	studentB := uuid.New()
	studentC := uuid.New()
	otherBranchID := uuid.New()

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}}

	tests := []struct {
		name        string
		caller      domain.Caller
		req         Request
		setupMocks  func(mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository)
		expectedErr error
	}{
		{
			name:   "success superadmin add and remove",
			caller: callerSuper,
			req:    Request{StudentIDs: []uuid.UUID{studentA, studentC}},
			setupMocks: func(mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockSR.On("GetBranchID", mock.Anything, studentA).Return(branchID, nil).Once()
				mockSR.On("GetBranchID", mock.Anything, studentC).Return(branchID, nil).Once()
				mockGR.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{studentA, studentB}, nil).Once()
				mockGR.On("RemoveStudents", mock.Anything, groupID, []uuid.UUID{studentB}, mock.AnythingOfType("time.Time")).Return(nil).Once()
				mockGR.On("AddStudents", mock.Anything, groupID, []uuid.UUID{studentC}, mock.AnythingOfType("time.Time")).Return(nil).Once()
			},
		},
		{
			name:   "success admin no-op diff",
			caller: callerAdmin,
			req:    Request{StudentIDs: []uuid.UUID{studentA, studentB}},
			setupMocks: func(mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockSR.On("GetBranchID", mock.Anything, studentA).Return(branchID, nil).Once()
				mockSR.On("GetBranchID", mock.Anything, studentB).Return(branchID, nil).Once()
				mockGR.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{studentA, studentB}, nil).Once()
			},
		},
		{
			name:   "success clear all with empty list",
			caller: callerSuper,
			req:    Request{StudentIDs: []uuid.UUID{}},
			setupMocks: func(mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockGR.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{studentA, studentB}, nil).Once()
				mockGR.On("RemoveStudents", mock.Anything, groupID, []uuid.UUID{studentA, studentB}, mock.AnythingOfType("time.Time")).Return(nil).Once()
			},
		},
		{
			name:   "error admin no branch access",
			caller: callerAdminNoAccess,
			req:    Request{StudentIDs: []uuid.UUID{studentA}},
			setupMocks: func(mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "error student not found",
			caller: callerSuper,
			req:    Request{StudentIDs: []uuid.UUID{studentA}},
			setupMocks: func(mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockSR.On("GetBranchID", mock.Anything, studentA).Return(uuid.Nil, errors.New("no rows in result set")).Once()
			},
			expectedErr: ErrStudentNotFound,
		},
		{
			name:   "error student branch mismatch",
			caller: callerSuper,
			req:    Request{StudentIDs: []uuid.UUID{studentA}},
			setupMocks: func(mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository) {
				mockGR.On("GetByID", mock.Anything, groupID).Return(&domain.Group{ID: groupID, BranchID: branchID}, nil).Once()
				mockSR.On("GetBranchID", mock.Anything, studentA).Return(uuid.New(), nil).Once()
			},
			expectedErr: ErrStudentBranchMismatch,
		},
		{
			name:   "error student_ids required",
			caller: callerSuper,
			req:    Request{},
			setupMocks: func(mockGR *mocks.GroupRepository, mockSR *mocks.StudentRepository) {
			},
			expectedErr: ErrStudentIDsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGR := new(mocks.GroupRepository)
			mockSR := new(mocks.StudentRepository)
			mockTx := new(mocks.MockTxManager)

			if tt.setupMocks != nil {
				tt.setupMocks(mockGR, mockSR)
			}

			uc := NewUseCase(mockGR, mockSR, mockTx)
			res, err := uc.Execute(context.Background(), tt.caller, groupID, tt.req)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, "success", res.Message)
			}

			mockGR.AssertExpectations(t)
			mockSR.AssertExpectations(t)
		})
	}
}
