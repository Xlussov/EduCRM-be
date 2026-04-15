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
	studentID := uuid.New()
	branchID := uuid.New()
	otherBranchID := uuid.New()
	userID := uuid.New()
	errDB := errors.New("db error")

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}}

	tests := []struct {
		name          string
		caller        domain.Caller
		mockSetup     func(repo *mocks.StudentRepository)
		expectedError error
		expectedMsg   string
		assertRepo    func(t *testing.T, repo *mocks.StudentRepository)
	}{
		{
			name:   "success as SUPERADMIN",
			caller: callerSuper,
			mockSetup: func(repo *mocks.StudentRepository) {
				repo.On("GetByID", mock.Anything, studentID).Return(&domain.Student{
					ID:       studentID,
					BranchID: branchID,
					Status:   domain.StatusActive,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, studentID, domain.StatusArchived).Return(nil).Once()
			},
			expectedError: nil,
			expectedMsg:   "success",
		},
		{
			name:   "success as ADMIN",
			caller: callerAdmin,
			mockSetup: func(repo *mocks.StudentRepository) {
				repo.On("GetByID", mock.Anything, studentID).Return(&domain.Student{
					ID:       studentID,
					BranchID: branchID,
					Status:   domain.StatusActive,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, studentID, domain.StatusArchived).Return(nil).Once()
			},
			expectedError: nil,
			expectedMsg:   "success",
		},
		{
			name:   "access denied",
			caller: callerAdminNoAccess,
			mockSetup: func(repo *mocks.StudentRepository) {
				repo.On("GetByID", mock.Anything, studentID).Return(&domain.Student{
					ID:       studentID,
					BranchID: branchID,
					Status:   domain.StatusActive,
				}, nil).Once()
			},
			expectedError: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, repo *mocks.StudentRepository) {
				repo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name:   "error_already_archived",
			caller: callerSuper,
			mockSetup: func(repo *mocks.StudentRepository) {
				repo.On("GetByID", mock.Anything, studentID).Return(&domain.Student{
					ID:       studentID,
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
			mockSetup: func(repo *mocks.StudentRepository) {
				repo.On("GetByID", mock.Anything, studentID).Return((*domain.Student)(nil), errDB).Once()
			},
			expectedError: errDB,
			expectedMsg:   "",
		},
		{
			name:   "error_db_on_updatestatus",
			caller: callerSuper,
			mockSetup: func(repo *mocks.StudentRepository) {
				repo.On("GetByID", mock.Anything, studentID).Return(&domain.Student{
					ID:       studentID,
					BranchID: branchID,
					Status:   domain.StatusActive,
				}, nil).Once()
				repo.On("UpdateStatus", mock.Anything, studentID, domain.StatusArchived).Return(errDB).Once()
			},
			expectedError: errDB,
			expectedMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.StudentRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, studentID)

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
