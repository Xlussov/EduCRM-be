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
	userID := uuid.New()
	studentID := uuid.New()
	branchID := uuid.New()
	dob := "2015-01-01"

	req := Request{
		FirstName:   "Jane",
		LastName:    "Smith",
		Dob:         &dob,
		ParentName:  "John Smith",
		ParentPhone: "12345678",
	}

	tests := []struct {
		name      string
		role      string
		req       Request
		setupMock func(sr *mocks.StudentRepository, ur *mocks.UserRepository)
		wantErr   error
	}{
		{
			name: "success as SUPERADMIN",
			role: "SUPERADMIN",
			req:  req,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("Update", mock.Anything, mock.MatchedBy(func(s *domain.Student) bool {
					return s.ID == studentID && s.FirstName == "Jane" && s.LastName == "Smith"
				})).Return(nil)
			},
		},
		{
			name: "success as ADMIN",
			role: "ADMIN",
			req:  req,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("GetBranchID", mock.Anything, studentID).Return(branchID, nil)
				ur.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branchID}, nil)
				sr.On("Update", mock.Anything, mock.MatchedBy(func(s *domain.Student) bool {
					return s.ID == studentID
				})).Return(nil)
			},
		},
		{
			name: "ADMIN access denied",
			role: "ADMIN",
			req:  req,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				otherBranch := uuid.New()
				sr.On("GetBranchID", mock.Anything, studentID).Return(branchID, nil)
				ur.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{otherBranch}, nil)
			},
			wantErr: ErrBranchAccessDenied,
		},
		{
			name: "error getting branch id",
			role: "ADMIN",
			req:  req,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("GetBranchID", mock.Anything, studentID).Return(uuid.Nil, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
		{
			name: "error getting user branch ids",
			role: "ADMIN",
			req:  req,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("GetBranchID", mock.Anything, studentID).Return(branchID, nil)
				ur.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{}, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
		{
			name: "update error",
			role: "SUPERADMIN",
			req:  req,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: errors.New("db error"),
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
			res, err := uc.Execute(context.Background(), userID, tt.role, studentID, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "success", res.Message)
			}

			sr.AssertExpectations(t)
			ur.AssertExpectations(t)
		})
	}
}
