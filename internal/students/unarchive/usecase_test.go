package unarchive

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
	studentID := uuid.New()
	errDB := errors.New("db error")

	tests := []struct {
		name          string
		mockSetup     func(repo *mocks.StudentRepository)
		expectedError error
		expectedMsg   string
	}{
		{
			name: "success",
			mockSetup: func(repo *mocks.StudentRepository) {
				repo.On("GetByID", mock.Anything, studentID).Return(&domain.Student{
					ID:     studentID,
					Status: domain.StatusArchived,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, studentID, domain.StatusActive).Return(nil).Once()
			},
			expectedError: nil,
			expectedMsg:   "success",
		},
		{
			name: "error_already_active",
			mockSetup: func(repo *mocks.StudentRepository) {
				repo.On("GetByID", mock.Anything, studentID).Return(&domain.Student{
					ID:     studentID,
					Status: domain.StatusActive,
				}, nil).Once()
			},
			expectedError: domain.ErrAlreadyActive,
			expectedMsg:   "",
		},
		{
			name: "error_db_on_getbyid",
			mockSetup: func(repo *mocks.StudentRepository) {
				repo.On("GetByID", mock.Anything, studentID).Return((*domain.Student)(nil), errDB).Once()
			},
			expectedError: errDB,
			expectedMsg:   "",
		},
		{
			name: "error_db_on_updatestatus",
			mockSetup: func(repo *mocks.StudentRepository) {
				repo.On("GetByID", mock.Anything, studentID).Return(&domain.Student{
					ID:     studentID,
					Status: domain.StatusArchived,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, studentID, domain.StatusActive).Return(errDB).Once()
			},
			expectedError: errDB,
			expectedMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.StudentRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), studentID)

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
