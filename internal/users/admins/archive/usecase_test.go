package archive

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
	adminID := uuid.New()
	errDB := errors.New("db error")

	tests := []struct {
		name          string
		mockSetup     func(repo *mocks.UserRepository)
		expectedError error
		expectedMsg   string
	}{
		{
			name: "success",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetByID", mock.Anything, adminID).Return(&domain.User{ID: adminID, Role: domain.RoleAdmin, IsActive: true}, nil).Once()
				repo.On("UpdateUserStatus", mock.Anything, adminID, false).Return(nil).Once()
			},
			expectedError: nil,
			expectedMsg:   "success",
		},
		{
			name: "already_archived",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetByID", mock.Anything, adminID).Return(&domain.User{ID: adminID, Role: domain.RoleAdmin, IsActive: false}, nil).Once()
			},
			expectedError: domain.ErrAlreadyArchived,
			expectedMsg:   "",
		},
		{
			name: "not_an_admin",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetByID", mock.Anything, adminID).Return(&domain.User{ID: adminID, Role: domain.RoleTeacher, IsActive: true}, nil).Once()
			},
			expectedError: domain.ErrNotFound,
			expectedMsg:   "",
		},
		{
			name: "repo_get_error",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetByID", mock.Anything, adminID).Return((*domain.User)(nil), errDB).Once()
			},
			expectedError: errDB,
			expectedMsg:   "",
		},
		{
			name: "repo_update_error",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetByID", mock.Anything, adminID).Return(&domain.User{ID: adminID, Role: domain.RoleAdmin, IsActive: true}, nil).Once()
				repo.On("UpdateUserStatus", mock.Anything, adminID, false).Return(errDB).Once()
			},
			expectedError: errDB,
			expectedMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), adminID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, res.Message)
			}

			repo.AssertExpectations(t)
		})
	}
}
