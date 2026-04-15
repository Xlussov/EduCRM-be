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
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Execute(t *testing.T) {
	subjectID := uuid.New()
	branchID := uuid.New()
	now := time.Now()
	callerAllowed := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerDenied := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}}

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
	errDB := errors.New("db error")

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.SubjectRepository)
		expectedErr   error
		expectedCount int
		assertRepo    func(t *testing.T, repo *mocks.SubjectRepository)
	}{
		{
			name:   "success",
			caller: callerAllowed,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetAll", mock.Anything, branchID, (*domain.EntityStatus)(nil)).Return(expected, nil)
			},
			expectedCount: 1,
		},
		{
			name:        "access denied",
			caller:      callerDenied,
			expectedErr: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, repo *mocks.SubjectRepository) {
				repo.AssertNotCalled(t, "GetAll", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name:   "db error",
			caller: callerAllowed,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetAll", mock.Anything, branchID, (*domain.EntityStatus)(nil)).Return(nil, errDB)
			},
			expectedErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.SubjectRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, Request{BranchID: branchID})

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Len(t, res.Subjects, tt.expectedCount)
				if tt.expectedCount > 0 {
					assert.Equal(t, subjectID.String(), res.Subjects[0].ID)
					assert.Equal(t, branchID.String(), res.Subjects[0].BranchID)
				}
			}

			if tt.assertRepo != nil {
				tt.assertRepo(t, repo)
				return
			}

			repo.AssertExpectations(t)
		})
	}
}
