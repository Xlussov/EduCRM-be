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
	branchID := uuid.New()
	callerAllowed := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerDenied := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}}
	req := Request{
		BranchID:    branchID,
		Name:        "Mathematics",
		Description: "Advanced Math",
	}
	expectedID := uuid.New()
	errDB := errors.New("db err")
	errBranch := errors.New("branch error")

	tests := []struct {
		name        string
		caller      domain.Caller
		mockSetup   func(repo *mocks.SubjectRepository, branchRepo *mocks.BranchRepository)
		expectedErr error
		expectedID  string
		assertRepo  func(t *testing.T, repo *mocks.SubjectRepository, branchRepo *mocks.BranchRepository)
	}{
		{
			name:   "success",
			caller: callerAllowed,
			mockSetup: func(repo *mocks.SubjectRepository, branchRepo *mocks.BranchRepository) {
				branchRepo.On("IsActive", mock.Anything, branchID).Return(true, nil)
				repo.On("Create", mock.Anything, mock.MatchedBy(func(s *domain.Subject) bool {
					if s.Name == req.Name && s.BranchID == branchID {
						s.ID = expectedID
						return true
					}
					return false
				})).Return(nil)
			},
			expectedID: expectedID.String(),
		},
		{
			name:        "access denied",
			caller:      callerDenied,
			expectedErr: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, repo *mocks.SubjectRepository, branchRepo *mocks.BranchRepository) {
				branchRepo.AssertNotCalled(t, "IsActive", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "branch archived",
			caller: callerAllowed,
			mockSetup: func(repo *mocks.SubjectRepository, branchRepo *mocks.BranchRepository) {
				branchRepo.On("IsActive", mock.Anything, branchID).Return(false, nil)
			},
			expectedErr: domain.ErrArchivedReference,
		},
		{
			name:   "branch lookup error",
			caller: callerAllowed,
			mockSetup: func(repo *mocks.SubjectRepository, branchRepo *mocks.BranchRepository) {
				branchRepo.On("IsActive", mock.Anything, branchID).Return(false, errBranch)
			},
			expectedErr: errBranch,
		},
		{
			name:   "db error",
			caller: callerAllowed,
			mockSetup: func(repo *mocks.SubjectRepository, branchRepo *mocks.BranchRepository) {
				branchRepo.On("IsActive", mock.Anything, branchID).Return(true, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(errDB)
			},
			expectedErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.SubjectRepository)
			branchRepo := new(mocks.BranchRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo, branchRepo)
			}

			uc := NewUseCase(repo, branchRepo)
			res, err := uc.Execute(context.Background(), tt.caller, req)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.expectedID, res.ID)
				assert.Equal(t, branchID.String(), res.BranchID)
			}

			if tt.assertRepo != nil {
				tt.assertRepo(t, repo, branchRepo)
				return
			}

			repo.AssertExpectations(t)
			branchRepo.AssertExpectations(t)
		})
	}
}
