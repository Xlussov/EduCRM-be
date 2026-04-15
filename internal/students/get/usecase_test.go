package get

import (
	"context"
	"errors"
	"testing"
	"time"

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
	dob := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	student := &domain.Student{
		ID:          studentID,
		BranchID:    branchID,
		FirstName:   "John",
		LastName:    "Doe",
		Dob:         &dob,
		ParentName:  "Jane Doe",
		ParentPhone: "+1234567890",
		Status:      domain.StatusActive,
		CreatedAt:   time.Now(),
	}

	callerSuper := domain.Caller{UserID: userID, Role: domain.RoleSuperadmin}
	callerAdmin := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}}
	callerTeacher := domain.Caller{UserID: userID, Role: domain.RoleTeacher, BranchIDs: []uuid.UUID{branchID}}
	callerAdminNoAccess := domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}}

	tests := []struct {
		name      string
		caller    domain.Caller
		setupMock func(sr *mocks.StudentRepository)
		wantErr   error
	}{
		{
			name:   "success as SUPERADMIN",
			caller: callerSuper,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(student, nil)
			},
		},
		{
			name:   "success as ADMIN with branch access",
			caller: callerAdmin,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(student, nil)
			},
		},
		{
			name:   "success as TEACHER with branch access",
			caller: callerTeacher,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(student, nil)
				sr.On("IsTeacherStudent", mock.Anything, userID, studentID).Return(true, nil)
			},
		},
		{
			name:   "teacher not assigned to student",
			caller: callerTeacher,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(student, nil)
				sr.On("IsTeacherStudent", mock.Anything, userID, studentID).Return(false, nil)
			},
			wantErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "access denied",
			caller: callerAdminNoAccess,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(student, nil)
			},
			wantErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "not found",
			caller: callerSuper,
			setupMock: func(sr *mocks.StudentRepository) {
				sr.On("GetByID", mock.Anything, studentID).Return(nil, errors.New("no rows"))
			},
			wantErr: ErrStudentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := new(mocks.StudentRepository)
			tt.setupMock(sr)

			uc := NewUseCase(sr)
			res, err := uc.Execute(context.Background(), tt.caller, studentID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, studentID, res.ID)
				assert.Equal(t, "2010-05-15", *res.Dob)
			}

			sr.AssertExpectations(t)
		})
	}
}
