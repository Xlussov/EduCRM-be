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
	subjectID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.SubjectRepository)
		repo.On("UpdateStatus", mock.Anything, subjectID, domain.StatusArchived).Return(nil)

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), subjectID)

		assert.Equal(t, "success", res.Message)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		repo := new(mocks.SubjectRepository)
		repo.On("UpdateStatus", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db err"))

		uc := NewUseCase(repo)
		_, err := uc.Execute(context.Background(), subjectID)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
