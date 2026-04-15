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
	branchID := uuid.New()
	callerAllowed := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerDenied := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}}
	req := Request{
		BranchID:    branchID,
		Name:        "New Math",
		Description: "New Desc",
	}
	subjectID := uuid.New()
	updatedDomain := &domain.Subject{
		ID:          subjectID,
		BranchID:    branchID,
		Name:        req.Name,
		Description: req.Description,
		Status:      domain.StatusActive,
	}
	errDB := errors.New("db error")

	tests := []struct {
		name        string
		caller      domain.Caller
		mockSetup   func(repo *mocks.SubjectRepository)
		expectedErr error
		expectedRes *Response
		assertRepo  func(t *testing.T, repo *mocks.SubjectRepository)
	}{
		{
			name:   "success",
			caller: callerAllowed,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{ID: subjectID, BranchID: branchID, Status: domain.StatusActive}, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(s *domain.Subject) bool {
					return s.Name == req.Name && s.Description == req.Description && s.ID == subjectID && s.BranchID == req.BranchID
				})).Return(updatedDomain, nil)
			},
			expectedRes: &Response{
				ID:          subjectID.String(),
				BranchID:    branchID.String(),
				Name:        req.Name,
				Description: req.Description,
				Status:      string(domain.StatusActive),
			},
		},
		{
			name:   "access denied",
			caller: callerDenied,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{ID: subjectID, BranchID: branchID, Status: domain.StatusActive}, nil)
			},
			expectedErr: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, repo *mocks.SubjectRepository) {
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "db error",
			caller: callerAllowed,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{ID: subjectID, BranchID: branchID, Status: domain.StatusActive}, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil, errDB)
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
			res, err := uc.Execute(context.Background(), tt.caller, subjectID, req)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				if tt.expectedRes != nil {
					assert.Equal(t, tt.expectedRes.ID, res.ID)
					assert.Equal(t, tt.expectedRes.BranchID, res.BranchID)
					assert.Equal(t, tt.expectedRes.Name, res.Name)
					assert.Equal(t, tt.expectedRes.Description, res.Description)
					assert.Equal(t, tt.expectedRes.Status, res.Status)
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
