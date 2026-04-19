package create_group

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
)

func TestUseCase_Execute(t *testing.T) {
	userID := uuid.New()
	branchID := uuid.New()
	teacherID := uuid.New()
	subjectID := uuid.New()
	groupID := uuid.New()
	student1 := uuid.New()
	student2 := uuid.New()

	validReq := Request{
		BranchID:  branchID,
		TeacherID: teacherID,
		SubjectID: subjectID,
		GroupID:   groupID,
		Date:      "2026-05-15",
		StartTime: "10:00",
		EndTime:   "11:00",
	}

	tests := []struct {
		name        string
		req         Request
		caller      domain.Caller
		mockSetup   func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository)
		expectedErr error
	}{
		{
			name:   "success",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil)
				gRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{student1, student2}, nil)
				sRepo.On("CheckStudentConflict", mock.Anything, student1, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil)
				sRepo.On("CheckStudentConflict", mock.Anything, student2, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil)

				sRepo.On("CreateLesson", mock.Anything, mock.MatchedBy(func(l *domain.Lesson) bool {
					return l.BranchID == branchID && l.TeacherID == teacherID && l.SubjectID == subjectID && *l.GroupID == groupID && l.StudentID == nil && l.Status == domain.LessonStatusScheduled
				})).Return(nil).Run(func(args mock.Arguments) {
					l := args.Get(1).(*domain.Lesson)
					l.ID = uuid.New()
				})
			},
			expectedErr: nil,
		},
		{
			name:   "teacher conflict",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(true, nil)
			},
			expectedErr: domain.ErrTeacherScheduleConflict,
		},
		{
			name:   "second student conflict aborts",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil)
				gRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{student1, student2}, nil)
				sRepo.On("CheckStudentConflict", mock.Anything, student1, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Once()
				sRepo.On("CheckStudentConflict", mock.Anything, student2, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(true, nil).Once()
			},
			expectedErr: domain.ErrStudentScheduleConflict,
		},
		{
			name:        "branch access denied",
			req:         validReq,
			caller:      domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}},
			mockSetup:   nil,
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "teacher not in branch",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(false, nil)
			},
			expectedErr: domain.ErrTeacherNotInBranch,
		},
		{
			name: "invalid time input",
			req: Request{
				BranchID:  branchID,
				TeacherID: teacherID,
				SubjectID: subjectID,
				GroupID:   groupID,
				Date:      "2026-05-15",
				StartTime: "25:00",
				EndTime:   "11:00",
			},
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
			},
			expectedErr: domain.ErrInvalidInput,
		},
		{
			name:   "group repo error",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil)
				gRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sRepo := new(mocks.ScheduleRepository)
			gRepo := new(mocks.GroupRepository)
			uRepo := new(mocks.UserRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(sRepo, gRepo, uRepo)
			}

			uc := NewUseCase(sRepo, gRepo, uRepo)
			res, err := uc.Execute(context.Background(), tt.caller, tt.req)

			if tt.expectedErr != nil {
				if tt.expectedErr.Error() == "db error" {
					assert.Error(t, err)
				} else {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
				assert.Empty(t, res.ID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, res.ID)
				assert.Equal(t, string(domain.LessonStatusScheduled), res.Status)
			}

			sRepo.AssertExpectations(t)
			gRepo.AssertExpectations(t)
			uRepo.AssertExpectations(t)
		})
	}
}
