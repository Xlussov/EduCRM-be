package list

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUseCase_Execute(t *testing.T) {
	subjectID := uuid.New()
	branchID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.SubjectRepository)
		expected := []*domain.Subject{
			{
				ID:          subjectID,
				BranchID:    branchID,
				Name:        "Math",
				Description: "Math desc",
				Status:      domain.StatusActive,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}
		repo.On("GetAll", context.Background(), branchID).Return(expected, nil)

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), Request{BranchID: branchID})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Subjects, 1)
		assert.Equal(t, subjectID.String(), res.Subjects[0].ID)
		assert.Equal(t, branchID.String(), res.Subjects[0].BranchID)
		repo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		repo := new(mocks.SubjectRepository)
		repo.On("GetAll", context.Background(), branchID).Return(nil, errors.New("db error"))

		uc := NewUseCase(repo)
		res, err := uc.Execute(context.Background(), Request{BranchID: branchID})

		assert.Error(t, err)
		assert.Nil(t, res)
		repo.AssertExpectations(t)
	})
}
