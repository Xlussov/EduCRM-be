package update

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
	req := Request{
		Name:        "New Math",
		Description: "New Desc",
	}
	subjectID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.SubjectRepository)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(s *domain.Subject) bool {
			return s.Name == req.Name && s.Description == req.Description && s.ID == subjectID
		})).Return(nil)

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), subjectID, req)

		assert.Equal(t, "success", res.Message)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		repo := new(mocks.SubjectRepository)
		repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))

		uc := NewUseCase(repo)
		_, err := uc.Execute(context.Background(), subjectID, req)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
