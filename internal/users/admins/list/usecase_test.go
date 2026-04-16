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
	adminID := uuid.New()
	branchID := uuid.New()
	errDB := errors.New("db error")

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.UserRepository)
		expectedError error
		expectedLen   int
	}{
		{
			name:   "success",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetAdmins", mock.Anything).Return([]*domain.UserWithBranches{
					{
						ID:        adminID,
						Phone:     "+998901112233",
						FirstName: "John",
						LastName:  "Admin",
						Role:      domain.RoleAdmin,
						IsActive:  true,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Branches: []domain.UserBranch{
							{ID: branchID, Name: "Main"},
						},
					},
				}, nil).Once()
			},
			expectedError: nil,
			expectedLen:   1,
		},
		{
			name:   "repo_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetAdmins", mock.Anything).Return(([]*domain.UserWithBranches)(nil), errDB).Once()
			},
			expectedError: errDB,
			expectedLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, Request{})

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tt.expectedLen)
				if tt.expectedLen > 0 {
					assert.Equal(t, adminID, res[0].ID)
					assert.Equal(t, "ACTIVE", res[0].Status)
					assert.Len(t, res[0].Branches, 1)
					assert.Equal(t, branchID, res[0].Branches[0].ID)
				}
			}

			repo.AssertExpectations(t)
		})
	}
}
