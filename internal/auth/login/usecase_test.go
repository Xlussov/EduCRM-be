package login

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
	"golang.org/x/crypto/bcrypt"
)

func TestUseCase_Execute(t *testing.T) {
	pwHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)

	activeUser := &domain.User{
		ID:           uuid.New(),
		Phone:        "123456",
		PasswordHash: string(pwHash),
		Role:         domain.RoleAdmin,
		IsActive:     true,
	}

	inactiveUser := &domain.User{
		ID:           uuid.New(),
		Phone:        "123456",
		PasswordHash: string(pwHash),
		IsActive:     false,
	}

	tests := []struct {
		name          string
		req           Request
		mockSetup     func(userRepo *mocks.UserRepository, authRepo *mocks.AuthRepository)
		expectedError string
	}{
		{
			name: "Success path",
			req:  Request{Phone: "123456", Password: "password123"},
			mockSetup: func(ur *mocks.UserRepository, ar *mocks.AuthRepository) {
				ur.On("GetByPhone", mock.Anything, "123456").Return(activeUser, nil)
				ur.On("GetUserBranchIDs", mock.Anything, activeUser.ID).Return([]uuid.UUID{uuid.New()}, nil)
				ar.On("SaveRefreshToken", mock.Anything, mock.Anything, activeUser.ID, mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "User not found",
			req:  Request{Phone: "123456", Password: "password123"},
			mockSetup: func(ur *mocks.UserRepository, ar *mocks.AuthRepository) {
				ur.On("GetByPhone", mock.Anything, "123456").Return(nil, errors.New("not found"))
			},
			expectedError: "invalid credentials",
		},
		{
			name: "Wrong password",
			req:  Request{Phone: "123456", Password: "wrong"},
			mockSetup: func(ur *mocks.UserRepository, ar *mocks.AuthRepository) {
				ur.On("GetByPhone", mock.Anything, "123456").Return(activeUser, nil)
			},
			expectedError: "invalid credentials",
		},
		{
			name: "Inactive user",
			req:  Request{Phone: "123456", Password: "password123"},
			mockSetup: func(ur *mocks.UserRepository, ar *mocks.AuthRepository) {
				ur.On("GetByPhone", mock.Anything, "123456").Return(inactiveUser, nil)
			},
			expectedError: "user is not active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.UserRepository)
			authRepo := new(mocks.AuthRepository)

			tt.mockSetup(userRepo, authRepo)

			uc := NewUseCase(userRepo, authRepo, "test-secret", 15*time.Minute, 720*time.Hour)
			res, err := uc.Execute(context.Background(), tt.req)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Empty(t, res.AccessToken)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, res.AccessToken)
				assert.NotEmpty(t, res.RefreshToken)
			}

			userRepo.AssertExpectations(t)
			authRepo.AssertExpectations(t)
		})
	}
}
