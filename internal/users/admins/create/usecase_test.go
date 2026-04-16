package create

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
	branchID := uuid.New()

	tests := []struct {
		name          string
		caller        domain.Caller
		req           Request
		mockSetup     func(userRepo *mocks.UserRepository)
		expectedError string
	}{
		{
			name:   "Success path",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req: Request{
				Phone:     "123456",
				Password:  "password123",
				FirstName: "Admin",
				LastName:  "Test",
				BranchIDs: []uuid.UUID{branchID},
			},
			mockSetup: func(ur *mocks.UserRepository) {
				ur.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.Phone == "123456" && u.FirstName == "Admin" && u.Role == domain.RoleAdmin
				})).Return(nil).Run(func(args mock.Arguments) {
					u := args.Get(1).(*domain.User)
					u.ID = uuid.New() // simulate DB setting ID
				})
				ur.On("CountActiveBranchesByIDs", mock.Anything, []uuid.UUID{branchID}).Return(1, nil)
				ur.On("AssignToBranches", mock.Anything, mock.AnythingOfType("uuid.UUID"), []uuid.UUID{branchID}).Return(nil)
			},
			expectedError: "",
		},
		{
			name:   "Create user error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req: Request{
				Phone:    "123456",
				Password: "pw",
			},
			mockSetup: func(ur *mocks.UserRepository) {
				ur.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.UserRepository)
			tt.mockSetup(userRepo)

			uc := NewUseCase(userRepo, &mocks.MockTxManager{})
			res, err := uc.Execute(context.Background(), tt.caller, tt.req)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, uuid.Nil, res.ID)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, res.ID)
			}
			userRepo.AssertExpectations(t)
		})
	}
}
