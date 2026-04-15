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
	errDB := errors.New("db error")
	branchID := uuid.New()
	callerSuper := domain.Caller{UserID: uuid.New(), Role: domain.RoleSuperadmin}
	callerAllowed := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerDenied := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}}

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.SubjectRepository)
		expectedError error
		expectedMsg   string
		assertRepo    func(t *testing.T, repo *mocks.SubjectRepository)
	}{
		{
			name:   "superadmin success",
			caller: callerSuper,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{
					ID:       subjectID,
					BranchID: branchID,
					Status:   domain.StatusActive,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, subjectID, domain.StatusArchived).Return(nil).Once()
			},
			expectedError: nil,
			expectedMsg:   "success",
		},
		{
			name:   "admin access denied",
			caller: callerDenied,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{
					ID:       subjectID,
					BranchID: branchID,
					Status:   domain.StatusActive,
				}, nil).Once()
			},
			expectedError: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, repo *mocks.SubjectRepository) {
				repo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name:   "admin access allowed",
			caller: callerAllowed,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{
					ID:       subjectID,
					BranchID: branchID,
					Status:   domain.StatusActive,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, subjectID, domain.StatusArchived).Return(nil).Once()
			},
			expectedError: nil,
			expectedMsg:   "success",
		},
		{
			name:   "error_already_archived",
			caller: callerSuper,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{
					ID:       subjectID,
					BranchID: branchID,
					Status:   domain.StatusArchived,
				}, nil).Once()
			},
			expectedError: domain.ErrAlreadyArchived,
			expectedMsg:   "",
		},
		{
			name:   "error_db_on_getbyid",
			caller: callerSuper,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return((*domain.Subject)(nil), errDB).Once()
			},
			expectedError: errDB,
			expectedMsg:   "",
		},
		{
			name:   "error_db_on_updatestatus",
			caller: callerSuper,
			mockSetup: func(repo *mocks.SubjectRepository) {
				repo.On("GetByID", mock.Anything, subjectID).Return(&domain.Subject{
					ID:       subjectID,
					BranchID: branchID,
					Status:   domain.StatusActive,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, subjectID, domain.StatusArchived).Return(errDB).Once()
			},
			expectedError: errDB,
			expectedMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.SubjectRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, subjectID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, res.Message)
			}

			if tt.assertRepo != nil {
				tt.assertRepo(t, repo)
				return
			}

			repo.AssertExpectations(t)
		})
	}
}
