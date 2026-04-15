package list

import (
	"context"
	"errors"
	"testing"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Execute(t *testing.T) {
	userID := uuid.New()
	branchID := uuid.New()
	otherBranchID := uuid.New()

	students := []*domain.Student{
		{ID: uuid.New(), FirstName: "John", LastName: "Doe", Status: domain.StatusActive},
		{ID: uuid.New(), FirstName: "Jane", LastName: "Smith", Status: domain.StatusArchived},
		{ID: uuid.New(), FirstName: "Alice", LastName: "Johnson", Status: domain.StatusActive},
	}

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerTeacher := domain.Caller{UserID: userID, Role: domain.RoleTeacher, BranchIDs: []uuid.UUID{branchID}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}}
	callerTeacherNoAccess := domain.Caller{UserID: userID, Role: domain.RoleTeacher, BranchIDs: []uuid.UUID{otherBranchID}}

	teacherStudents := []*domain.Student{students[0]}

	errDB := errors.New("db error")

	tests := []struct {
		name       string
		caller     domain.Caller
		req        Request
		setupMock  func(sr *mocks.StudentRepository)
		wantCount  int
		wantErr    error
		assertRepo func(t *testing.T, sr *mocks.StudentRepository)
	}{
		{
			name:   "success as SUPERADMIN - no filters",
			caller: callerSuper,
			req:    Request{BranchID: branchID},
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return(students, nil)
			},
			wantCount: 3,
		},
		{
			name:   "filter by status ACTIVE",
			caller: callerSuper,
			req:    Request{BranchID: branchID, Status: "ACTIVE"},
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return([]*domain.Student{students[0], students[2]}, nil)
			},
			wantCount: 2,
		},
		{
			name:   "search by name",
			caller: callerSuper,
			req:    Request{BranchID: branchID, Search: "john"},
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return(students, nil)
			},
			wantCount: 2,
		},
		{
			name:   "search and status combined",
			caller: callerSuper,
			req:    Request{BranchID: branchID, Search: "john", Status: "ACTIVE"},
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return([]*domain.Student{students[0], students[2]}, nil)
			},
			wantCount: 2,
		},
		{
			name:   "ADMIN with branch access",
			caller: callerAdmin,
			req:    Request{BranchID: branchID},
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return(students, nil)
			},
			wantCount: 3,
		},
		{
			name:   "TEACHER with branch access",
			caller: callerTeacher,
			req:    Request{BranchID: branchID},
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByBranchIDAndTeacherID", mock.Anything, branchID, userID, mock.Anything).Return(teacherStudents, nil)
			},
			wantCount: 1,
		},
		{
			name:    "TEACHER without branch access",
			caller:  callerTeacherNoAccess,
			req:     Request{BranchID: branchID},
			wantErr: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, sr *mocks.StudentRepository) {
				sr.AssertNotCalled(t, "GetByBranchIDAndTeacherID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name:    "ADMIN without branch access",
			caller:  callerAdminNoAccess,
			req:     Request{BranchID: branchID},
			wantErr: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, sr *mocks.StudentRepository) {
				sr.AssertNotCalled(t, "GetByBranchID", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name:    "missing branch_id",
			caller:  callerSuper,
			req:     Request{},
			wantErr: ErrBranchIDRequired,
		},
		{
			name:   "db error",
			caller: callerSuper,
			req:    Request{BranchID: branchID},
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return(nil, errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := new(mocks.StudentRepository)
			if tt.setupMock != nil {
				tt.setupMock(sr)
			}

			uc := NewUseCase(sr)
			res, err := uc.Execute(context.Background(), tt.caller, tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Len(t, res.Students, tt.wantCount)
			}

			if tt.assertRepo != nil {
				tt.assertRepo(t, sr)
				return
			}

			sr.AssertExpectations(t)
		})
	}
}
