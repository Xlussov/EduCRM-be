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
	studentID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(sr *mocks.StudentRepository)
		wantMsg   string
		wantErr   string
	}{
		{
			name: "success",
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("UpdateStatus", mock.Anything, studentID, domain.StatusArchived).Return(nil)
			},
			wantMsg: "success",
		},
		{
			name: "db error",
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("UpdateStatus", mock.Anything, studentID, domain.StatusArchived).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := new(mocks.StudentRepository)
			tt.setupMock(sr)

			uc := NewUseCase(sr)
			res, err := uc.Execute(context.Background(), studentID)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMsg, res.Message)
			}

			sr.AssertExpectations(t)
		})
	}
}
