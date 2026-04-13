package update

import (
	"context"
	"testing"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Execute(t *testing.T) {
	adminID := uuid.New()
	branchID := uuid.New()
	newPhone := "+998901234567"

	req := Request{
		FirstName: "Updated",
		LastName:  "Admin",
		Phone:     newPhone,
		BranchIDs: []uuid.UUID{branchID},
	}

	tests := []struct {
		name          string
		mockSetup     func(repo *mocks.UserRepository)
		expectedError error
	}{
		{
			name: "success",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetByID", mock.Anything, adminID).Return(&domain.User{ID: adminID, Role: domain.RoleAdmin, IsActive: true}, nil).Once()
				repo.On("CountActiveBranchesByIDs", mock.Anything, req.BranchIDs).Return(1, nil).Once()
				repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.ID == adminID && u.Phone == newPhone && u.FirstName == req.FirstName && u.LastName == req.LastName
				})).Return(nil).Once()
				repo.On("DeleteUserBranches", mock.Anything, adminID).Return(nil).Once()
				repo.On("AssignToBranches", mock.Anything, adminID, req.BranchIDs).Return(nil).Once()
				repo.On("GetWithBranchesByID", mock.Anything, adminID).Return(&domain.UserWithBranches{
					ID:        adminID,
					Phone:     newPhone,
					FirstName: req.FirstName,
					LastName:  req.LastName,
					Role:      domain.RoleAdmin,
					IsActive:  true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Branches:  []domain.UserBranch{{ID: branchID, Name: "Main"}},
				}, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "error_already_archived",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetByID", mock.Anything, adminID).Return(&domain.User{ID: adminID, Role: domain.RoleAdmin, IsActive: false}, nil).Once()
			},
			expectedError: domain.ErrAlreadyArchived,
		},
		{
			name: "phone_conflict",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetByID", mock.Anything, adminID).Return(&domain.User{ID: adminID, Role: domain.RoleAdmin, IsActive: true}, nil).Once()
				repo.On("CountActiveBranchesByIDs", mock.Anything, req.BranchIDs).Return(1, nil).Once()
				repo.On("UpdateUser", mock.Anything, mock.Anything).Return(domain.ErrAlreadyExists).Once()
			},
			expectedError: domain.ErrPhoneAlreadyExists,
		},
		{
			name: "archived_branch_in_request",
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetByID", mock.Anything, adminID).Return(&domain.User{ID: adminID, Role: domain.RoleAdmin, IsActive: true}, nil).Once()
				repo.On("CountActiveBranchesByIDs", mock.Anything, req.BranchIDs).Return(0, nil).Once()
			},
			expectedError: domain.ErrArchivedReference,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepository)
			tt.mockSetup(repo)

			uc := NewUseCase(repo, &mocks.MockTxManager{})
			res, err := uc.Execute(context.Background(), adminID, req)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, adminID, res.ID)
				assert.Equal(t, "ACTIVE", res.Status)
				assert.Len(t, res.Branches, 1)
			}

			repo.AssertExpectations(t)
		})
	}
}
