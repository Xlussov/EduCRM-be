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
	callerSuper := domain.Caller{UserID: uuid.New(), Role: domain.RoleSuperadmin}
	callerAdminAllowed := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerAdminDenied := domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}}
	errDB := errors.New("db err")
	req := Request{
		Name:    "Updated Name",
		Address: "Updated Address",
		City:    "Updated City",
	}
	updatedDomain := &domain.Branch{
		ID:      branchID,
		Name:    req.Name,
		Address: req.Address,
		City:    req.City,
		Status:  domain.StatusActive,
	}

	tests := []struct {
		name        string
		caller      domain.Caller
		mockSetup   func(repo *mocks.BranchRepository)
		expectedErr error
		expectedRes *Response
		assertRepo  func(t *testing.T, repo *mocks.BranchRepository)
	}{
		{
			name:   "success",
			caller: callerSuper,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return(&domain.Branch{ID: branchID, Status: domain.StatusActive}, nil)
				repo.On("Update", mock.Anything, &domain.Branch{
					ID:      branchID,
					Name:    req.Name,
					Address: req.Address,
					City:    req.City,
				}).Return(updatedDomain, nil)
			},
			expectedRes: &Response{
				ID:      branchID.String(),
				Name:    req.Name,
				Address: req.Address,
				City:    req.City,
				Status:  string(domain.StatusActive),
			},
		},
		{
			name:        "admin access denied",
			caller:      callerAdminDenied,
			expectedErr: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, repo *mocks.BranchRepository) {
				repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "admin access allowed",
			caller: callerAdminAllowed,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return(&domain.Branch{ID: branchID, Status: domain.StatusActive}, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(&domain.Branch{ID: branchID, Status: domain.StatusActive}, nil)
			},
		},
		{
			name:   "db error",
			caller: callerSuper,
			mockSetup: func(repo *mocks.BranchRepository) {
				repo.On("GetByID", mock.Anything, branchID).Return(&domain.Branch{ID: branchID, Status: domain.StatusActive}, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil, errDB)
			},
			expectedErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.BranchRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, branchID, req)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				if tt.expectedRes != nil {
					assert.Equal(t, tt.expectedRes.ID, res.ID)
					assert.Equal(t, tt.expectedRes.Name, res.Name)
					assert.Equal(t, tt.expectedRes.Address, res.Address)
					assert.Equal(t, tt.expectedRes.City, res.City)
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
