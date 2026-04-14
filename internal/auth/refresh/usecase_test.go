package refresh

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/Xlussov/EduCRM-be/pkg/hash"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Execute(t *testing.T) {
	reqToken := "some-refresh-token"
	hashedToken := hash.SHA256Token(reqToken)

	activeUser := &domain.User{
		ID:       uuid.New(),
		Role:     domain.RoleAdmin,
		IsActive: true,
	}

	validToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    activeUser.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: false,
	}

	expiredToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    activeUser.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(-time.Hour),
		IsRevoked: false,
	}

	reusedToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    activeUser.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: true,
	}

	tests := []struct {
		name          string
		req           Request
		mockSetup     func(userRepo *mocks.UserRepository, authRepo *mocks.AuthRepository)
		expectedError string
	}{
		{
			name: "Success path",
			req:  Request{RefreshToken: reqToken},
			mockSetup: func(ur *mocks.UserRepository, ar *mocks.AuthRepository) {
				ar.On("GetRefreshToken", mock.Anything, hashedToken).Return(validToken, nil)
				ar.On("RevokeRefreshToken", mock.Anything, validToken.ID).Return(nil)
				ur.On("GetByID", mock.Anything, validToken.UserID).Return(activeUser, nil)
				ur.On("GetUserBranchIDs", mock.Anything, validToken.UserID).Return([]uuid.UUID{uuid.New()}, nil)
				ar.On("SaveRefreshToken", mock.Anything, mock.Anything, validToken.UserID, mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "Invalid token",
			req:  Request{RefreshToken: reqToken},
			mockSetup: func(ur *mocks.UserRepository, ar *mocks.AuthRepository) {
				ar.On("GetRefreshToken", mock.Anything, hashedToken).Return(nil, errors.New("not found"))
			},
			expectedError: "invalid refresh token",
		},
		{
			name: "Expired token",
			req:  Request{RefreshToken: reqToken},
			mockSetup: func(ur *mocks.UserRepository, ar *mocks.AuthRepository) {
				ar.On("GetRefreshToken", mock.Anything, hashedToken).Return(expiredToken, nil)
			},
			expectedError: "refresh token expired",
		},
		{
			name: "Token reused",
			req:  Request{RefreshToken: reqToken},
			mockSetup: func(ur *mocks.UserRepository, ar *mocks.AuthRepository) {
				ar.On("GetRefreshToken", mock.Anything, hashedToken).Return(reusedToken, nil)
				ar.On("RevokeAllUserTokens", mock.Anything, activeUser.ID).Return(nil)
			},
			expectedError: "token reused",
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
