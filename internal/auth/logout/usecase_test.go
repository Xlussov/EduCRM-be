package logout

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/Xlussov/EduCRM-be/pkg/hash"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_Execute(t *testing.T) {
	rawToken := "some-raw-refresh-token"
	hashedToken := hash.SHA256Token(rawToken)
	tokenID := uuid.New()
	userID := uuid.New()

	mockDbToken := &domain.RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: false,
	}

	tests := []struct {
		name        string
		req         Request
		setupMocks  func(mockAR *mocks.AuthRepository)
		expectedErr error
		expectedRes Response
	}{
		{
			name: "Success path",
			req:  Request{RefreshToken: rawToken},
			setupMocks: func(mockAR *mocks.AuthRepository) {
				mockAR.On("GetRefreshToken", mock.Anything, hashedToken).Return(mockDbToken, nil).Once()
				mockAR.On("RevokeRefreshToken", mock.Anything, tokenID).Return(nil).Once()
			},
			expectedErr: nil,
			expectedRes: Response{Message: "Successfully logged out"},
		},
		{
			name: "Token not found (idempotent success)",
			req:  Request{RefreshToken: rawToken},
			setupMocks: func(mockAR *mocks.AuthRepository) {
				mockAR.On("GetRefreshToken", mock.Anything, hashedToken).Return(nil, errors.New("no rows in result set")).Once()
			},
			expectedErr: nil,
			expectedRes: Response{Message: "Successfully logged out"},
		},
		{
			name: "Token already revoked (idempotent success)",
			req:  Request{RefreshToken: rawToken},
			setupMocks: func(mockAR *mocks.AuthRepository) {
				revokedToken := *mockDbToken
				revokedToken.IsRevoked = true
				mockAR.On("GetRefreshToken", mock.Anything, hashedToken).Return(&revokedToken, nil).Once()
			},
			expectedErr: nil,
			expectedRes: Response{Message: "Successfully logged out"},
		},
		{
			name: "Database error on get",
			req:  Request{RefreshToken: rawToken},
			setupMocks: func(mockAR *mocks.AuthRepository) {
				mockAR.On("GetRefreshToken", mock.Anything, hashedToken).Return(nil, errors.New("db connection lost")).Once()
			},
			expectedErr: errors.New("db connection lost"),
			expectedRes: Response{},
		},
		{
			name: "Database error on revoke",
			req:  Request{RefreshToken: rawToken},
			setupMocks: func(mockAR *mocks.AuthRepository) {
				mockAR.On("GetRefreshToken", mock.Anything, hashedToken).Return(mockDbToken, nil).Once()
				mockAR.On("RevokeRefreshToken", mock.Anything, tokenID).Return(errors.New("failed to update")).Once()
			},
			expectedErr: errors.New("failed to update"),
			expectedRes: Response{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAR := new(mocks.AuthRepository)
			if tt.setupMocks != nil {
				tt.setupMocks(mockAR)
			}

			uc := NewUseCase(mockAR)
			res, err := uc.Execute(context.Background(), tt.req)

			if tt.expectedErr != nil {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedRes, res)
			}

			mockAR.AssertExpectations(t)
		})
	}
}
