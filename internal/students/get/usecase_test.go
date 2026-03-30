package get

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
	studentID := uuid.New()
	branchID := uuid.New()
	dob := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	student := &domain.Student{
		ID:          studentID,
		BranchID:    branchID,
		FirstName:   "John",
		LastName:    "Doe",
		Dob:         &dob,
		ParentName:  "Jane Doe",
		ParentPhone: "+1234567890",
		Status:      domain.StatusActive,
		CreatedAt:   time.Now(),
	}

	tests := []struct {
		name      string
		setupMock func(sr *mocks.StudentRepository)
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(student, nil)
			},
		},
		{
			name: "not found",
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(nil, errors.New("no rows"))
			},
			wantErr: ErrStudentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := new(mocks.StudentRepository)
			tt.setupMock(sr)

			uc := NewUseCase(sr)
			res, err := uc.Execute(context.Background(), studentID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, studentID, res.ID)
				assert.Equal(t, "2010-05-15", *res.Dob)
			}

			sr.AssertExpectations(t)
		})
	}
}
