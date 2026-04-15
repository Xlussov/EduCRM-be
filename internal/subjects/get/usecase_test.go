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
	subjectID := uuid.New()
	branchID := uuid.New()
	now := time.Now()
	callerSuper := domain.Caller{UserID: uuid.New(), Role: domain.RoleSuperadmin}
	callerAllowed := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerDenied := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}}
	errDB := errors.New("db error")

	expectedSubject := &domain.Subject{
		ID:          subjectID,
		BranchID:    branchID,
		Name:        "Test Subject",
		Description: "A description",
		Status:      domain.StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name        string
		caller      domain.Caller
		req         Request
		mockSetup   func(repo *mocks.SubjectRepository)
		wantRes     *Response
		expectedErr error
		assertRepo  func(t *testing.T, repo *mocks.SubjectRepository)
	}{
		{
			name:   "success",
			caller: callerSuper,
			req:    Request{ID: subjectID},
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(expectedSubject, nil)
			},
			wantRes: &Response{
				Subject: SubjectResponse{
					ID:          subjectID.String(),
					BranchID:    branchID.String(),
					Name:        "Test Subject",
					Description: "A description",
					Status:      domain.StatusActive,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:   "admin access denied",
			caller: callerDenied,
			req:    Request{ID: subjectID},
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(expectedSubject, nil)
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "admin access allowed",
			caller: callerAllowed,
			req:    Request{ID: subjectID},
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(expectedSubject, nil)
			},
			wantRes: &Response{
				Subject: SubjectResponse{
					ID:          subjectID.String(),
					BranchID:    branchID.String(),
					Name:        "Test Subject",
					Description: "A description",
					Status:      domain.StatusActive,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:   "repo error",
			caller: callerSuper,
			req:    Request{ID: subjectID},
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(nil, errDB)
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
			res, err := uc.Execute(context.Background(), tt.caller, tt.req)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRes, res)
			}

			if tt.assertRepo != nil {
				tt.assertRepo(t, repo)
				return
			}

			repo.AssertExpectations(t)
		})
	}
}
