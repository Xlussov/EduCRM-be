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
	userID := uuid.New()
	branchID := uuid.New()
	otherBranchID := uuid.New()
	studentID := uuid.New()
	dob := "2010-05-15"

	validReq := Request{
		BranchID:    branchID,
		FirstName:   "John",
		LastName:    "Doe",
		Dob:         &dob,
		ParentName:  "Jane Doe",
		ParentPhone: "+1234567890",
	}

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}}

	errDB := errors.New("db error")

	tests := []struct {
		name       string
		caller     domain.Caller
		req        Request
		setupMock  func(sr *mocks.StudentRepository, ur *mocks.UserRepository)
		wantErr    error
		wantID     uuid.UUID
		assertRepo func(t *testing.T, sr *mocks.StudentRepository, ur *mocks.UserRepository)
	}{
		{
			name:   "success as SUPERADMIN",
			caller: callerSuper,
			req:    validReq,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				ur.On("IsBranchActive", mock.Anything, branchID).Return(true, nil)
				sr.On("Create", mock.Anything, mock.AnythingOfType("*domain.Student")).Return(nil).Run(func(args mock.Arguments) {
					s := args.Get(1).(*domain.Student)
					s.ID = studentID
				})
			},
			wantID: studentID,
		},
		{
			name:   "success as ADMIN with branch access",
			caller: callerAdmin,
			req:    validReq,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				ur.On("IsBranchActive", mock.Anything, branchID).Return(true, nil)
				sr.On("Create", mock.Anything, mock.AnythingOfType("*domain.Student")).Return(nil).Run(func(args mock.Arguments) {
					s := args.Get(1).(*domain.Student)
					s.ID = studentID
				})
			},
			wantID: studentID,
		},
		{
			name:    "ADMIN without branch access",
			caller:  callerAdminNoAccess,
			req:     validReq,
			wantErr: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
				ur.AssertNotCalled(t, "IsBranchActive", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "branch is archived",
			caller: callerSuper,
			req:    validReq,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				ur.On("IsBranchActive", mock.Anything, branchID).Return(false, nil)
			},
			wantErr: domain.ErrArchivedReference,
		},
		{
			name:   "db error on IsBranchActive",
			caller: callerSuper,
			req:    validReq,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				ur.On("IsBranchActive", mock.Anything, branchID).Return(false, errDB)
			},
			wantErr: errDB,
		},
		{
			name:   "db error on create",
			caller: callerSuper,
			req:    validReq,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				ur.On("IsBranchActive", mock.Anything, branchID).Return(true, nil)
				sr.On("Create", mock.Anything, mock.Anything).Return(errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := new(mocks.StudentRepository)
			ur := new(mocks.UserRepository)
			if tt.setupMock != nil {
				tt.setupMock(sr, ur)
			}

			uc := NewUseCase(sr, ur)
			res, err := uc.Execute(context.Background(), tt.caller, tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, res.ID)
			}

			if tt.assertRepo != nil {
				tt.assertRepo(t, sr, ur)
				return
			}
			sr.AssertExpectations(t)
			ur.AssertExpectations(t)
		})
	}
}
