package create

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
		Name:        "Mathematics",
		Description: "Advanced Math",
	}

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.SubjectRepository)
		expectedID := uuid.New()

		// Using mock.MatchedBy to inject ID into the passed reference
		repo.On("Create", mock.Anything, mock.MatchedBy(func(s *domain.Subject) bool {
			if s.Name == req.Name {
				s.ID = expectedID
				return true
			}
			return false
		})).Return(nil)

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, expectedID.String(), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("empty name", func(t *testing.T) {
		repo := new(mocks.SubjectRepository)
		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), Request{Description: "No Name"})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "subject name is required", err.Error())
	})

	t.Run("db error", func(t *testing.T) {
		repo := new(mocks.SubjectRepository)
		repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db err"))

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
