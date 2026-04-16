package get

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
	teacherID := uuid.New()
	branchID := uuid.New()
	otherBranchID := uuid.New()

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.UserRepository)
		expectedError error
		expectedID    uuid.UUID
	}{
		{
			name:   "superadmin_success",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:        teacherID,
					Phone:     "+998901234567",
					FirstName: "Jane",
					LastName:  "Teacher",
					Role:      domain.RoleTeacher,
					IsActive:  true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Branches:  []domain.UserBranch{{ID: branchID, Name: "Main"}},
				}, nil).Once()
			},
			expectedError: nil,
			expectedID:    teacherID,
		},
		{
			name:   "admin_success_with_shared_branch",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{
					ID:       teacherID,
					Role:     domain.RoleTeacher,
					IsActive: true,
					Branches: []domain.UserBranch{{ID: branchID, Name: "Main"}},
				}, nil).Once()
			},
			expectedError: nil,
			expectedID:    teacherID,
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
			name:   "not_teacher",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(&domain.UserWithBranches{ID: teacherID, Role: domain.RoleAdmin}, nil).Once()
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name:   "repo_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, teacherID).Return(nil, errors.New("db error")).Once()
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, Request{ID: teacherID})

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, domain.ErrBranchAccessDenied) || errors.Is(tt.expectedError, domain.ErrNotFound) {
					assert.ErrorIs(t, err, tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, res.ID)
			}

			repo.AssertExpectations(t)
		})
	}
}
