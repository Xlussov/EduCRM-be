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

	students := []*domain.Student{
		{ID: uuid.New(), FirstName: "John", LastName: "Doe", Status: domain.StatusActive},
		{ID: uuid.New(), FirstName: "Jane", LastName: "Smith", Status: domain.StatusArchived},
		{ID: uuid.New(), FirstName: "Alice", LastName: "Johnson", Status: domain.StatusActive},
	}

	tests := []struct {
		name      string
		role      string
		req       Request
		setupMock func(sr *mocks.StudentRepository, ur *mocks.UserRepository)
		wantCount int
		wantErr   error
	}{
		{
			name: "success as SUPERADMIN - no filters",
			role: "SUPERADMIN",
			req:  Request{BranchID: branchID},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return(students, nil)
			},
			wantCount: 3,
		},
		{
			name: "filter by status ACTIVE",
			role: "SUPERADMIN",
			req:  Request{BranchID: branchID, Status: "ACTIVE"},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return([]*domain.Student{students[0], students[2]}, nil)
			},
			wantCount: 2,
		},
		{
			name: "search by name",
			role: "SUPERADMIN",
			req:  Request{BranchID: branchID, Search: "john"},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return(students, nil)
			},
			wantCount: 2,
		},
		{
			name: "search and status combined",
			role: "SUPERADMIN",
			req:  Request{BranchID: branchID, Search: "john", Status: "ACTIVE"},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return([]*domain.Student{students[0], students[2]}, nil)
			},
			wantCount: 2,
		},
		{
			name: "ADMIN with branch access",
			role: "ADMIN",
			req:  Request{BranchID: branchID},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				ur.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branchID}, nil)
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return(students, nil)
			},
			wantCount: 3,
		},
		{
			name: "ADMIN without branch access",
			role: "ADMIN",
			req:  Request{BranchID: branchID},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				ur.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{uuid.New()}, nil)
			},
			wantErr: ErrBranchAccessDenied,
		},
		{
			name:      "missing branch_id",
			role:      "SUPERADMIN",
			req:       Request{},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {},
			wantErr:   ErrBranchIDRequired,
		},
		{
			name: "db error",
			role: "SUPERADMIN",
			req:  Request{BranchID: branchID},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("GetByBranchID", mock.Anything, branchID, mock.Anything).Return(nil, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := new(mocks.StudentRepository)
			ur := new(mocks.UserRepository)
			tt.setupMock(sr, ur)

			uc := NewUseCase(sr, ur)
			res, err := uc.Execute(context.Background(), userID, tt.role, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Len(t, res.Students, tt.wantCount)
			}

			sr.AssertExpectations(t)
			ur.AssertExpectations(t)
		})
	}
}
