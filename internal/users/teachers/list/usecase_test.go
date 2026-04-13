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
	teacherID := uuid.New()

	tests := []struct {
		name          string
		role          string
		branchIDs     []uuid.UUID
		mockSetup     func(repo *mocks.UserRepository)
		expectedCount int
		expectedError error
	}{
		{
			name:      "superadmin_success",
			role:      "SUPERADMIN",
			branchIDs: nil,
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
			name:      "admin_success_filtered",
			role:      "ADMIN",
			branchIDs: []uuid.UUID{branchID},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetTeachers", mock.Anything, []uuid.UUID{branchID}).Return([]*domain.UserWithBranches{}, nil).Once()
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name:      "admin_no_branches",
			role:      "ADMIN",
			branchIDs: []uuid.UUID{},
			mockSetup: func(repo *mocks.UserRepository) {
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name:      "repo_error",
			role:      "SUPERADMIN",
			branchIDs: nil,
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetTeachers", mock.Anything, []uuid.UUID(nil)).Return(nil, errors.New("db error")).Once()
			},
			expectedCount: 0,
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.role, tt.branchIDs, Request{})

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tt.expectedCount)
			}

			repo.AssertExpectations(t)
		})
	}
}
