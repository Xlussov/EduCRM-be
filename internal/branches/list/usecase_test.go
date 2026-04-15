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
	userID := uuid.New()
	branchID := uuid.New()

	branch := &domain.Branch{
		ID:        branchID,
		Name:      "Test",
		Address:   "Test",
		City:      "Test City",
		Status:    domain.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	errDB := errors.New("db err")

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.BranchRepository)
		expectedErr   error
		expectedCount int
		expectedID    uuid.UUID
	}{
		{
			name:   "superadmin - gets all",
			caller: domain.Caller{UserID: userID, Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetAll", mock.Anything, (*domain.EntityStatus)(nil)).Return([]*domain.Branch{branch}, nil)
			},
			expectedCount: 1,
			expectedID:    branchID,
		},
		{
			name:   "admin - gets by param",
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin},
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByUserID", mock.Anything, userID, (*domain.EntityStatus)(nil)).Return([]*domain.Branch{branch}, nil)
			},
			expectedCount: 1,
			expectedID:    branchID,
		},
		{
			name:   "db error",
			caller: domain.Caller{UserID: userID, Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetAll", mock.Anything, (*domain.EntityStatus)(nil)).Return(nil, errDB)
			},
			expectedErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.BranchRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, Request{})

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tt.expectedCount)
				if tt.expectedCount > 0 {
					assert.Equal(t, tt.expectedID, res[0].ID)
				}
			}

			repo.AssertExpectations(t)
		})
	}
}
