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
	subjectID := uuid.New()
	errDB := errors.New("db error")

	tests := []struct {
		name          string
		mockSetup     func(repo *mocks.SubjectRepository)
		expectedError error
		expectedMsg   string
	}{
		{
			name: "success",
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{
					ID:     subjectID,
					Status: domain.StatusArchived,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, subjectID, domain.StatusActive).Return(nil).Once()
			},
			expectedError: nil,
			expectedMsg:   "success",
		},
		{
			name: "error_already_active",
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{
					ID:     subjectID,
					Status: domain.StatusActive,
				}, nil).Once()
			},
			expectedError: domain.ErrAlreadyActive,
			expectedMsg:   "",
		},
		{
			name: "error_db_on_getbyid",
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return((*domain.Subject)(nil), errDB).Once()
			},
			expectedError: errDB,
			expectedMsg:   "",
		},
		{
			name: "error_db_on_updatestatus",
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{
					ID:     subjectID,
					Status: domain.StatusArchived,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, subjectID, domain.StatusActive).Return(errDB).Once()
			},
			expectedError: errDB,
			expectedMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.SubjectRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), subjectID)

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
