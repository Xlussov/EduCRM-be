package archive

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
	teacherID := uuid.New()
	branchID := uuid.New()
	otherBranchID := uuid.New()
	errDB := errors.New("db error")

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.UserRepository)
		expectedError error
		expectedMsg   string
	}{
		{
			name:   "superadmin_success",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: true,
					Branches: []domain.UserBranch{{ID: branchID, Name: "Main"}},
				}, nil).Once()
				repo.On("UpdateUserStatus", mock.Anything, teacherID, false).Return(nil).Once()
			},
			expectedError: nil,
			expectedMsg:   "success",
		},
		{
			name:   "admin_no_access",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: true,
					Branches: []domain.UserBranch{{ID: branchID, Name: "Main"}},
				}, nil).Once()
			},
			expectedError: domain.ErrBranchAccessDenied,
		},
		{
			name:   "already_archived",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{ID: teacherID, Role: domain.RoleTeacher, IsActive: false}, nil).Once()
			},
			expectedError: domain.ErrAlreadyArchived,
		},
		{
			name:   "not_teacher",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{ID: teacherID, Role: domain.RoleAdmin, IsActive: true}, nil).Once()
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name:   "repo_get_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return((*domain.UserWithBranches)(nil), errDB).Once()
			},
			expectedError: errDB,
		},
		{
			name:   "repo_update_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{ID: teacherID, Role: domain.RoleTeacher, IsActive: true}, nil).Once()
				repo.On("UpdateUserStatus", mock.Anything, teacherID, false).Return(errDB).Once()
			},
			expectedError: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, teacherID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, res.Message)
			}

			repo.AssertExpectations(t)
		})
	}
}
