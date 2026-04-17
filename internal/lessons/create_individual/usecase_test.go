package create_individual

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
)

func TestUseCase_Execute(t *testing.T) {
	userID := uuid.New()
	branchID := uuid.New()
	teacherID := uuid.New()
	subjectID := uuid.New()
	studentID := uuid.New()

	validReq := Request{
		BranchID:  branchID,
		TeacherID: teacherID,
		SubjectID: subjectID,
		StudentID: &studentID,
		Date:      "2026-05-15",
		StartTime: "10:00",
		EndTime:   "11:00",
	}

	errDB := errors.New("db err")

	tests := []struct {
		name          string
		req           Request
		caller        domain.Caller
		mockSetup     func(repo *mocks.ScheduleRepository)
		expectedErr   error
		expectedCount int
	}{
		{
			name:   "admin - success",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("CheckStudentConflict", mock.Anything, studentID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil)
				repo.On("CheckTeacherConflict", mock.Anything, teacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil)
				repo.On("CreateLesson", mock.Anything, mock.MatchedBy(func(l *domain.Lesson) bool {
					return l.BranchID == branchID && l.TeacherID == teacherID && l.SubjectID == subjectID && *l.StudentID == studentID && l.Status == domain.LessonStatusScheduled
				})).Return(nil).Run(func(args mock.Arguments) {
					l := args.Get(1).(*domain.Lesson)
					l.ID = uuid.New()
				})
			},
			expectedErr:   nil,
			expectedCount: 1,
		},
		{
			name:   "admin - teacher conflict",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("CheckStudentConflict", mock.Anything, studentID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil)
				repo.On("CheckTeacherConflict", mock.Anything, teacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(true, nil)
			},
			expectedErr: domain.ErrTeacherScheduleConflict,
		},
		{
			name:   "admin - student conflict",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("CheckStudentConflict", mock.Anything, studentID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(true, nil)
			},
			expectedErr: domain.ErrStudentScheduleConflict,
		},
		{
			name:        "admin - branch access denied",
			req:         validReq,
			caller:      domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}}, // Diff branch
			mockSetup:   nil,
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name: "superadmin - success without branch access check",
			req: Request{
				BranchID:  branchID,
				TeacherID: teacherID,
				SubjectID: subjectID,
				StudentID: nil,
				Date:      "2026-05-15",
				StartTime: "10:00",
				EndTime:   "11:00",
			},
			caller: domain.Caller{UserID: userID, Role: domain.RoleSuperadmin, BranchIDs: []uuid.UUID{}},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("CheckTeacherConflict", mock.Anything, teacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil)
				repo.On("CreateLesson", mock.Anything, mock.MatchedBy(func(l *domain.Lesson) bool {
					return l.BranchID == branchID && l.StudentID == nil && l.Status == domain.LessonStatusScheduled
				})).Return(nil).Run(func(args mock.Arguments) {
					l := args.Get(1).(*domain.Lesson)
					l.ID = uuid.New()
				})
			},
			expectedErr:   nil,
			expectedCount: 1,
		},
		{
			name:   "db error",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleSuperadmin},
			mockSetup: func(repo *mocks.ScheduleRepository) {
				repo.On("CheckStudentConflict", mock.Anything, studentID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, errDB)
			},
			expectedErr: errDB,
		},
		{
			name: "invalid time input",
			req: Request{
				BranchID:  branchID,
				TeacherID: teacherID,
				SubjectID: subjectID,
				Date:      "2026-05-15",
				StartTime: "25:00",
				EndTime:   "11:00",
			},
			caller:      domain.Caller{UserID: userID, Role: domain.RoleSuperadmin},
			mockSetup:   nil,
			expectedErr: domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.ScheduleRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, tt.req)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Empty(t, res.ID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, res.ID)
				assert.Equal(t, string(domain.LessonStatusScheduled), res.Status)
			}

			repo.AssertExpectations(t)
		})
	}
}
