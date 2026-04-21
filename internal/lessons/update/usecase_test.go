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
	lessonID := uuid.New()
	branchID := uuid.New()
	teacherID := uuid.New()
	newTeacherID := uuid.New()
	subjectID := uuid.New()
	newSubjectID := uuid.New()
	studentID := uuid.New()
	groupID := uuid.New()

	validReq := Request{
		Date:      "2026-05-15",
		StartTime: "10:00",
		EndTime:   "11:00",
		TeacherID: newTeacherID,
		SubjectID: newSubjectID,
	}

	validLesson := &domain.Lesson{
		ID:        lessonID,
		BranchID:  branchID,
		TeacherID: teacherID,
		SubjectID: subjectID,
		StudentID: &studentID,
		Status:    domain.LessonStatusScheduled,
	}

	groupLesson := &domain.Lesson{
		ID:        lessonID,
		BranchID:  branchID,
		TeacherID: teacherID,
		SubjectID: subjectID,
		GroupID:   &groupID,
		Status:    domain.LessonStatusScheduled,
	}

	errDB := errors.New("db err")

	tests := []struct {
		name          string
		caller        domain.Caller
		req           Request
		lesson        *domain.Lesson
		mockSetup     func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository)
		expectedError error
	}{
		{
			name:   "success_individual",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			req:    validReq,
			lesson: validLesson,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(validLesson, nil).Once()
				userRepo.On("CheckTeacherInBranch", mock.Anything, newTeacherID, branchID).Return(true, nil).Once()
				repo.On("CheckStudentConflictExcludingLesson", mock.Anything, studentID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), lessonID).Return(false, nil).Once()
				repo.On("CheckTeacherConflictExcludingLesson", mock.Anything, newTeacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), lessonID).Return(false, nil).Once()
				repo.On("UpdateLesson", mock.Anything, mock.MatchedBy(func(l *domain.Lesson) bool {
					return l.ID == lessonID && l.TeacherID == newTeacherID && l.SubjectID == newSubjectID
				})).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "success_group",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			req:    validReq,
			lesson: groupLesson,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(groupLesson, nil).Once()
				userRepo.On("CheckTeacherInBranch", mock.Anything, newTeacherID, branchID).Return(true, nil).Once()
				groupRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return([]uuid.UUID{studentID}, nil).Once()
				repo.On("CheckStudentConflictExcludingLesson", mock.Anything, studentID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), lessonID).Return(false, nil).Once()
				repo.On("CheckTeacherConflictExcludingLesson", mock.Anything, newTeacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), lessonID).Return(false, nil).Once()
				repo.On("UpdateLesson", mock.Anything, mock.MatchedBy(func(l *domain.Lesson) bool {
					return l.ID == lessonID && l.TeacherID == newTeacherID && l.SubjectID == newSubjectID
				})).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "branch_access_denied",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}},
			req:    validReq,
			lesson: validLesson,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(validLesson, nil).Once()
			},
			expectedError: domain.ErrBranchAccessDenied,
		},
		{
			name:   "teacher_not_in_branch",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			req:    validReq,
			lesson: validLesson,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(validLesson, nil).Once()
				userRepo.On("CheckTeacherInBranch", mock.Anything, newTeacherID, branchID).Return(false, nil).Once()
			},
			expectedError: domain.ErrTeacherNotInBranch,
		},
		{
			name:   "lesson_not_scheduled",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req:    validReq,
			lesson: &domain.Lesson{ID: lessonID, BranchID: branchID, Status: domain.LessonStatusCompleted},
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(&domain.Lesson{ID: lessonID, BranchID: branchID, Status: domain.LessonStatusCompleted}, nil).Once()
			},
			expectedError: domain.ErrLessonNotScheduled,
		},
		{
			name:   "invalid_time",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req: Request{
				Date:      "2026-05-15",
				StartTime: "12:00",
				EndTime:   "11:00",
				TeacherID: newTeacherID,
				SubjectID: newSubjectID,
			},
			lesson: validLesson,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(validLesson, nil).Once()
				userRepo.On("CheckTeacherInBranch", mock.Anything, newTeacherID, branchID).Return(true, nil).Once()
			},
			expectedError: domain.ErrInvalidInput,
		},
		{
			name:   "teacher_conflict",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req:    validReq,
			lesson: validLesson,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(validLesson, nil).Once()
				userRepo.On("CheckTeacherInBranch", mock.Anything, newTeacherID, branchID).Return(true, nil).Once()
				repo.On("CheckStudentConflictExcludingLesson", mock.Anything, studentID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), lessonID).Return(false, nil).Once()
				repo.On("CheckTeacherConflictExcludingLesson", mock.Anything, newTeacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), lessonID).Return(true, nil).Once()
			},
			expectedError: domain.ErrTeacherScheduleConflict,
		},
		{
			name:   "student_conflict",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req:    validReq,
			lesson: validLesson,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(validLesson, nil).Once()
				userRepo.On("CheckTeacherInBranch", mock.Anything, newTeacherID, branchID).Return(true, nil).Once()
				repo.On("CheckStudentConflictExcludingLesson", mock.Anything, studentID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), lessonID).Return(true, nil).Once()
			},
			expectedError: domain.ErrStudentScheduleConflict,
		},
		{
			name:   "partial overlap teacher conflict",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req: Request{
				Date:      "2026-05-15",
				StartTime: "10:30",
				EndTime:   "11:30",
				TeacherID: newTeacherID,
				SubjectID: newSubjectID,
			},
			lesson: validLesson,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(validLesson, nil).Once()
				userRepo.On("CheckTeacherInBranch", mock.Anything, newTeacherID, branchID).Return(true, nil).Once()
				repo.On("CheckStudentConflictExcludingLesson", mock.Anything, studentID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), lessonID).Return(false, nil).Once()
				repo.On("CheckTeacherConflictExcludingLesson", mock.Anything, newTeacherID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), lessonID).Return(true, nil).Once()
			},
			expectedError: domain.ErrTeacherScheduleConflict,
		},
		{
			name:   "partial overlap student conflict",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req: Request{
				Date:      "2026-05-15",
				StartTime: "09:30",
				EndTime:   "10:30",
				TeacherID: newTeacherID,
				SubjectID: newSubjectID,
			},
			lesson: validLesson,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(validLesson, nil).Once()
				userRepo.On("CheckTeacherInBranch", mock.Anything, newTeacherID, branchID).Return(true, nil).Once()
				repo.On("CheckStudentConflictExcludingLesson", mock.Anything, studentID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"), lessonID).Return(true, nil).Once()
			},
			expectedError: domain.ErrStudentScheduleConflict,
		},
		{
			name:   "group_students_error",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req:    validReq,
			lesson: groupLesson,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(groupLesson, nil).Once()
				userRepo.On("CheckTeacherInBranch", mock.Anything, newTeacherID, branchID).Return(true, nil).Once()
				groupRepo.On("GetActiveStudentIDs", mock.Anything, groupID).Return(nil, errDB).Once()
			},
			expectedError: errDB,
		},
		{
			name:   "lesson_not_found",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			req:    validReq,
			lesson: nil,
			mockSetup: func(repo *mocks.ScheduleRepository, groupRepo *mocks.GroupRepository, userRepo *mocks.UserRepository) {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(nil, domain.ErrNotFound).Once()
			},
			expectedError: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.ScheduleRepository)
			groupRepo := new(mocks.GroupRepository)
			userRepo := new(mocks.UserRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo, groupRepo, userRepo)
			}

			uc := NewUseCase(repo, groupRepo, userRepo)
			res, err := uc.Execute(context.Background(), tt.caller, lessonID, tt.req)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Empty(t, res.ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, lessonID, res.ID)
				assert.Equal(t, string(domain.LessonStatusScheduled), res.Status)
			}

			repo.AssertExpectations(t)
			groupRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}
