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
	adminID := uuid.New()
	branchID := uuid.New()
	errDB := errors.New("db error")

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.UserRepository)
		expectedError error
	}{
		{
			name:   "success",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, adminID).Return(&domain.UserWithBranches{
					ID:        adminID,
					Phone:     "+998901234567",
					FirstName: "John",
					LastName:  "Admin",
					Role:      domain.RoleAdmin,
					IsActive:  true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Branches: []domain.UserBranch{
						{ID: branchID, Name: "Main"},
					},
				}, nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "repo_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetWithBranchesByID", mock.Anything, adminID).Return((*domain.UserWithBranches)(nil), errDB).Once()
			},
			expectedError: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, Request{ID: adminID})

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, adminID, res.ID)
				assert.Equal(t, "ACTIVE", res.Status)
				assert.Len(t, res.Branches, 1)
			}

			repo.AssertExpectations(t)
		})
	}
}
