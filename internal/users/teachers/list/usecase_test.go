package list

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Execute(t *testing.T) {
	branchID := uuid.New()
	otherBranchID := uuid.New()
	teacherID := uuid.New()
	errDB := errors.New("db error")

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.UserRepository)
		expectedCount int
		expectedError error
	}{
		{
			name:   "superadmin_success",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetTeachers", mock.Anything, []uuid.UUID(nil)).Return([]*domain.UserWithBranches{
					{
						ID:        teacherID,
						Phone:     "+998901234567",
						FirstName: "Jane",
						LastName:  "Teacher",
						Role:      domain.RoleTeacher,
						IsActive:  true,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Branches:  []domain.UserBranch{{ID: branchID, Name: "Main"}},
					},
				}, nil).Once()
			},
			expectedCount: 1,
			expectedError: nil,
		},
		{
			name:   "superadmin_branch_filter",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetTeachers", mock.Anything, []uuid.UUID{branchID}).Return([]*domain.UserWithBranches{}, nil).Once()
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name:   "admin_success_filtered",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetTeachers", mock.Anything, []uuid.UUID{branchID}).Return([]*domain.UserWithBranches{}, nil).Once()
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name:   "admin_branch_filter_allowed",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID, otherBranchID}},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetTeachers", mock.Anything, []uuid.UUID{otherBranchID}).Return([]*domain.UserWithBranches{}, nil).Once()
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name:          "admin_branch_filter_denied",
			caller:        domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup:     func(repo *mocks.UserRepository) {},
			expectedCount: 0,
			expectedError: domain.ErrBranchAccessDenied,
		},
		{
			name:   "admin_no_branches",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{}},
			mockSetup: func(repo *mocks.UserRepository) {
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name:   "repo_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetTeachers", mock.Anything, []uuid.UUID(nil)).Return(nil, errDB).Once()
			},
			expectedCount: 0,
			expectedError: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo)
			var req Request
			switch tt.name {
			case "superadmin_branch_filter":
				req = Request{BranchID: &branchID}
			case "admin_branch_filter_allowed", "admin_branch_filter_denied":
				req = Request{BranchID: &otherBranchID}
			default:
				req = Request{}
			}
			res, err := uc.Execute(context.Background(), tt.caller, req)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tt.expectedCount)
			}

			repo.AssertExpectations(t)
		})
	}
}
