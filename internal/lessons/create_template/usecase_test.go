package create_template

import (
	"context"
	"errors"
	"testing"
	"time"

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
		BranchID:   branchID,
		TeacherID:  teacherID,
		SubjectID:  subjectID,
		GroupID:    &groupID,
		DaysOfWeek: []int32{5}, // Friday
		StartTime:  "10:00",
		EndTime:    "11:00",
		StartDate:  "2026-05-01", // Friday
		EndDate:    "2026-05-15", // Friday (total 3 matches: 1st, 8th, 15th)
	}

	multiDayReq := Request{
		BranchID:   branchID,
		TeacherID:  teacherID,
		SubjectID:  subjectID,
		GroupID:    &groupID,
		DaysOfWeek: []int32{1, 5}, // Monday, Friday
		StartTime:  "10:00",
		EndTime:    "11:00",
		StartDate:  "2026-05-01", // Friday
		EndDate:    "2026-05-11", // Monday (total 4 matches: 1st, 4th, 8th, 11th)
	}

	tests := []struct {
		name        string
		req         Request
		caller      domain.Caller
		mockSetup   func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository)
		expectedErr error
		expectedCnt int
		expectedCnf int
	}{
		{
			name:   "full success",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
				gRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{student1, student2}, nil)

				sRepo.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(tmp *domain.Template) bool {
					return tmp.BranchID == branchID && tmp.TeacherID == teacherID
				})).Return(nil).Run(func(args mock.Arguments) {
					l := args.Get(1).(*domain.Template)
					l.ID = uuid.New()
				})

				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Times(3)
				sRepo.On("CheckStudentConflict", mock.Anything, student1, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Times(3)
				sRepo.On("CheckStudentConflict", mock.Anything, student2, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Times(3)

				sRepo.On("BulkCreateLessons", mock.Anything, mock.MatchedBy(func(ls []domain.Lesson) bool {
					return len(ls) == 3
				})).Return(nil)
			},
			expectedErr: nil,
			expectedCnt: 3,
			expectedCnf: 0,
		},
		{
			name:   "multi-day success",
			req:    multiDayReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
				gRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{student1, student2}, nil)

				sRepo.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(tmp *domain.Template) bool {
					return tmp.BranchID == branchID && tmp.TeacherID == teacherID
				})).Return(nil).Run(func(args mock.Arguments) {
					l := args.Get(1).(*domain.Template)
					l.ID = uuid.New()
				})

				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Times(4)
				sRepo.On("CheckStudentConflict", mock.Anything, student1, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Times(4)
				sRepo.On("CheckStudentConflict", mock.Anything, student2, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Times(4)

				sRepo.On("BulkCreateLessons", mock.Anything, mock.MatchedBy(func(ls []domain.Lesson) bool {
					return len(ls) == 4
				})).Return(nil)
			},
			expectedErr: nil,
			expectedCnt: 4,
			expectedCnf: 0,
		},
		{
			name: "partial overlap teacher conflict",
			req: Request{
				BranchID:   branchID,
				TeacherID:  teacherID,
				SubjectID:  subjectID,
				GroupID:    &groupID,
				DaysOfWeek: []int32{5}, // Friday
				StartTime:  "10:30",
				EndTime:    "11:30", // overlaps with 10:00 - 11:00
				StartDate:  "2026-05-01",
				EndDate:    "2026-05-01", // Friday (just 1 match)
			},
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
				gRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{student1, student2}, nil)

				sRepo.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(tmp *domain.Template) bool {
					return true
				})).Return(nil).Run(func(args mock.Arguments) {
					l := args.Get(1).(*domain.Template)
					l.ID = uuid.New()
				})

				date1, _ := time.Parse("2006-01-02", "2026-05-01")
				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, date1, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(true, nil).Once() // partial overlap conflict
			},
			expectedErr: nil,
			expectedCnt: 0,
			expectedCnf: 1,
		},
		{
			name: "partial overlap student conflict",
			req: Request{
				BranchID:   branchID,
				TeacherID:  teacherID,
				SubjectID:  subjectID,
				GroupID:    &groupID,
				DaysOfWeek: []int32{5}, // Friday
				StartTime:  "09:30",
				EndTime:    "10:30", // overlaps with 10:00 - 11:00
				StartDate:  "2026-05-01",
				EndDate:    "2026-05-01", // Friday (just 1 match)
			},
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
				gRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{student1, student2}, nil)

				sRepo.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(tmp *domain.Template) bool {
					return true
				})).Return(nil).Run(func(args mock.Arguments) {
					l := args.Get(1).(*domain.Template)
					l.ID = uuid.New()
				})

				date1, _ := time.Parse("2006-01-02", "2026-05-01")
				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, date1, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Once()
				sRepo.On("CheckStudentConflict", mock.Anything, student1, date1, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(true, nil).Once() // partial overlap conflict
			},
			expectedErr: nil,
			expectedCnt: 0,
			expectedCnf: 1,
		},
		{
			name:   "partial success (1 date conflicted)",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
				gRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{student1, student2}, nil)

				sRepo.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(tmp *domain.Template) bool {
					return true
				})).Return(nil).Run(func(args mock.Arguments) {
					l := args.Get(1).(*domain.Template)
					l.ID = uuid.New()
				})

				date1, _ := time.Parse("2006-01-02", "2026-05-01")
				date2, _ := time.Parse("2006-01-02", "2026-05-08")
				date3, _ := time.Parse("2006-01-02", "2026-05-15")

				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, date1, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Once()
				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, date2, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(true, nil).Once() // Conflict here
				sRepo.On("CheckTeacherConflict", mock.Anything, teacherID, date3, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Once()

				sRepo.On("CheckStudentConflict", mock.Anything, student1, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Times(2)
				sRepo.On("CheckStudentConflict", mock.Anything, student2, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(false, nil).Times(2)

				sRepo.On("BulkCreateLessons", mock.Anything, mock.MatchedBy(func(ls []domain.Lesson) bool {
					return len(ls) == 2
				})).Return(nil)
			},
			expectedErr: nil,
			expectedCnt: 2,
			expectedCnf: 1,
		},
		{
			name:        "access denied",
			req:         validReq,
			caller:      domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}},
			mockSetup:   nil,
			expectedErr: domain.ErrBranchAccessDenied,
			expectedCnt: 0,
			expectedCnf: 0,
		},
		{
			name:   "teacher not in branch",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(false, nil)
			},
			expectedErr: domain.ErrTeacherNotInBranch,
			expectedCnt: 0,
			expectedCnf: 0,
		},
		{
			name:   "db error during get students",
			req:    validReq,
			caller: domain.Caller{UserID: userID, Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			mockSetup: func(sRepo *mocks.ScheduleRepository, gRepo *mocks.GroupRepository, uRepo *mocks.UserRepository) {
				uRepo.On("CheckTeacherInBranch", mock.Anything, teacherID, branchID).Return(true, nil)
				gRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
			expectedCnt: 0,
			expectedCnf: 0,
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
				assert.Empty(t, res.TemplateID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, res.TemplateID)
				assert.Equal(t, tt.expectedCnt, res.CreatedLessonsCount)
				assert.Len(t, res.Conflicts, tt.expectedCnf)
			}

			sRepo.AssertExpectations(t)
			gRepo.AssertExpectations(t)
			uRepo.AssertExpectations(t)
		})
	}
}
