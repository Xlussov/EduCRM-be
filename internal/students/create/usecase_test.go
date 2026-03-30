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

	tests := []struct {
		name      string
		role      string
		req       Request
		setupMock func(sr *mocks.StudentRepository, ur *mocks.UserRepository)
		wantErr   error
		wantID    uuid.UUID
	}{
		{
			name: "success as SUPERADMIN",
			role: "SUPERADMIN",
			req:  validReq,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("Create", mock.Anything, mock.MatchedBy(func(s *domain.Student) bool {
					s.ID = studentID
					return s.FirstName == "John" && s.LastName == "Doe" && s.BranchID == branchID
				})).Return(nil)
			},
			wantID: studentID,
		},
		{
			name: "success as ADMIN with branch access",
			role: "ADMIN",
			req:  validReq,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				ur.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{branchID}, nil)
				sr.On("Create", mock.Anything, mock.MatchedBy(func(s *domain.Student) bool {
					s.ID = studentID
					return s.BranchID == branchID
				})).Return(nil)
			},
			wantID: studentID,
		},
		{
			name: "ADMIN without branch access",
			role: "ADMIN",
			req:  validReq,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				otherBranch := uuid.New()
				ur.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID{otherBranch}, nil)
			},
			wantErr: ErrBranchAccessDenied,
		},
		{
			name: "empty first name",
			role: "SUPERADMIN",
			req: Request{
				BranchID:    branchID,
				LastName:    "Doe",
				ParentName:  "Jane",
				ParentPhone: "+1234567890",
			},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {},
			wantErr:   ErrFirstNameRequired,
		},
		{
			name: "empty last name",
			role: "SUPERADMIN",
			req: Request{
				BranchID:    branchID,
				FirstName:   "John",
				ParentName:  "Jane",
				ParentPhone: "+1234567890",
			},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {},
			wantErr:   ErrLastNameRequired,
		},
		{
			name: "empty parent name",
			role: "SUPERADMIN",
			req: Request{
				BranchID:    branchID,
				FirstName:   "John",
				LastName:    "Doe",
				ParentPhone: "+1234567890",
			},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {},
			wantErr:   ErrParentNameRequired,
		},
		{
			name: "empty parent phone",
			role: "SUPERADMIN",
			req: Request{
				BranchID:   branchID,
				FirstName:  "John",
				LastName:   "Doe",
				ParentName: "Jane",
			},
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {},
			wantErr:   ErrParentPhoneRequired,
		},
		{
			name: "invalid dob format",
			role: "SUPERADMIN",
			req: func() Request {
				bad := "not-a-date"
				r := validReq
				r.Dob = &bad
				return r
			}(),
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {},
			wantErr:   errors.New("invalid date of birth format, expected YYYY-MM-DD"),
		},
		{
			name: "db error on create",
			role: "SUPERADMIN",
			req:  validReq,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				sr.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
		{
			name: "db error fetching branch IDs",
			role: "ADMIN",
			req:  validReq,
			setupMock: func(sr *mocks.StudentRepository, ur *mocks.UserRepository) {
				ur.On("GetUserBranchIDs", mock.Anything, userID).Return([]uuid.UUID(nil), errors.New("db error"))
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
			res, err := uc.Execute(context.Background(), userID, tt.role, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, res.ID)
			}

			sr.AssertExpectations(t)
			ur.AssertExpectations(t)
		})
	}
}
