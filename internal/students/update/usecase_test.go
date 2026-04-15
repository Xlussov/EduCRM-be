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
	otherBranchID := uuid.New()
	dob := "2015-01-01"

	req := Request{
		FirstName:   "Jane",
		LastName:    "Smith",
		Dob:         &dob,
		ParentName:  "John Smith",
		ParentPhone: "12345678",
	}

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}}

	errDB := errors.New("db error")

	tests := []struct {
		name       string
		caller     domain.Caller
		req        Request
		setupMock  func(sr *mocks.StudentRepository)
		wantErr    error
		assertRepo func(t *testing.T, sr *mocks.StudentRepository)
	}{
		{
			name:   "success as SUPERADMIN",
			caller: callerSuper,
			req:    req,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(&domain.Student{ID: studentID, BranchID: branchID, Status: domain.StatusActive}, nil)
				sr.On("Update", mock.Anything, mock.MatchedBy(func(s *domain.Student) bool {
					return s.ID == studentID && s.FirstName == "Jane" && s.LastName == "Smith"
				})).Return(&domain.Student{
					ID:        studentID,
					BranchID:  branchID,
					FirstName: "Jane",
					LastName:  "Smith",
				}, nil)
			},
		},
		{
			name:   "success as ADMIN",
			caller: callerAdmin,
			req:    req,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(&domain.Student{ID: studentID, BranchID: branchID, Status: domain.StatusActive}, nil)
				sr.On("Update", mock.Anything, mock.MatchedBy(func(s *domain.Student) bool {
					return s.ID == studentID
				})).Return(&domain.Student{
					ID:        studentID,
					BranchID:  branchID,
					FirstName: "Jane",
					LastName:  "Smith",
				}, nil)
			},
		},
		{
			name:   "ADMIN access denied",
			caller: callerAdminNoAccess,
			req:    req,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(&domain.Student{ID: studentID, BranchID: branchID, Status: domain.StatusActive}, nil)
			},
			wantErr: domain.ErrBranchAccessDenied,
			assertRepo: func(t *testing.T, sr *mocks.StudentRepository) {
				sr.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "cannot edit archived",
			caller: callerSuper,
			req:    req,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(&domain.Student{ID: studentID, BranchID: branchID, Status: domain.StatusArchived}, nil)
			},
			wantErr: domain.ErrCannotEditArchived,
		},
		{
			name:   "error getting student",
			caller: callerSuper,
			req:    req,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return((*domain.Student)(nil), errDB)
			},
			wantErr: errDB,
		},
		{
			name:   "update error",
			caller: callerSuper,
			req:    req,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(&domain.Student{ID: studentID, BranchID: branchID, Status: domain.StatusActive}, nil)
				sr.On("Update", mock.Anything, mock.Anything).Return(nil, errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := new(mocks.StudentRepository)
			if tt.setupMock != nil {
				tt.setupMock(sr)
			}

			uc := NewUseCase(sr)
			res, err := uc.Execute(context.Background(), tt.caller, studentID, tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, studentID.String(), res.ID)
				assert.Equal(t, "Jane", res.FirstName)
				assert.Equal(t, "Smith", res.LastName)
			}

			if tt.assertRepo != nil {
				tt.assertRepo(t, sr)
				return
			}
			sr.AssertExpectations(t)
		})
	}
}
