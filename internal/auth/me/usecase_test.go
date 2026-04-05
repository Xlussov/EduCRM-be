package me

import (
	"context"
	"errors"
	"testing"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_Execute(t *testing.T) {
	userID := uuid.New()

	mockUser := &domain.User{
		ID:        userID,
		FirstName: "Alexander",
		LastName:  "Dmitriev",
		Phone:     "+1234567890",
		Role:      "ADMIN",
	}

	tests := []struct {
		name        string
		setupMocks  func(mockUR *mocks.UserRepository)
		expectedErr error
		expectedRes Response
	}{
		{
			name: "Success",
			setupMocks: func(mockUR *mocks.UserRepository) {
				mockUR.On("GetByID", mock.Anything, userID).Return(mockUser, nil).Once()
			},
			expectedErr: nil,
			expectedRes: Response{
				ID:        userID,
				FirstName: mockUser.FirstName,
				LastName:  mockUser.LastName,
				Phone:     mockUser.Phone,
				Role:      mockUser.Role,
			},
		},
		{
			name: "User not found (e.g. deleted)",
			setupMocks: func(mockUR *mocks.UserRepository) {
				mockUR.On("GetByID", mock.Anything, userID).Return(nil, errors.New("no rows in result set")).Once()
			},
			expectedErr: errors.New("no rows in result set"),
			expectedRes: Response{},
		},
		{
			name: "Database error",
			setupMocks: func(mockUR *mocks.UserRepository) {
				mockUR.On("GetByID", mock.Anything, userID).Return(nil, errors.New("db connection lost")).Once()
			},
			expectedErr: errors.New("db connection lost"),
			expectedRes: Response{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUR := new(mocks.UserRepository)
			if tt.setupMocks != nil {
				tt.setupMocks(mockUR)
			}

			uc := NewUseCase(mockUR)
			res, err := uc.Execute(context.Background(), userID)

			if tt.expectedErr != nil {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedRes, res)
			}

			mockUR.AssertExpectations(t)
		})
	}
}
