package get

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
)

func TestUseCase_Execute(t *testing.T) {
	lessonID := uuid.New()
	branchID := uuid.New()
	otherBranchID := uuid.New()
	teacherID := uuid.New()
	otherTeacherID := uuid.New()
	subjectID := uuid.New()
	studentID := uuid.New()
	groupID := uuid.New()

	date, _ := time.Parse(dateLayout, "2026-05-01")
	start, _ := time.Parse(timeLayout, "10:00")
	end, _ := time.Parse(timeLayout, "11:00")

	baseLesson := &domain.Lesson{
		ID:        lessonID,
		BranchID:  branchID,
		TeacherID: teacherID,
		SubjectID: subjectID,
		StudentID: &studentID,
		GroupID:   &groupID,
		Date:      date,
		StartTime: start,
		EndTime:   end,
		Status:    domain.LessonStatusScheduled,
	}

	tests := []struct {
		name        string
		caller      domain.Caller
		lesson      *domain.Lesson
		repoErr     error
		expectedErr error
	}{
		{
			name:   "superadmin_success",
			caller: domain.Caller{Role: domain.RoleSuperadmin},
			lesson: baseLesson,
		},
		{
			name:   "admin_success",
			caller: domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			lesson: baseLesson,
		},
		{
			name:        "admin_access_denied",
			caller:      domain.Caller{Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{otherBranchID}},
			lesson:      baseLesson,
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "teacher_success",
			caller: domain.Caller{Role: domain.RoleTeacher, UserID: teacherID, BranchIDs: []uuid.UUID{branchID}},
			lesson: baseLesson,
		},
		{
			name:        "teacher_branch_denied",
			caller:      domain.Caller{Role: domain.RoleTeacher, UserID: teacherID, BranchIDs: []uuid.UUID{otherBranchID}},
			lesson:      baseLesson,
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:   "teacher_other_lesson",
			caller: domain.Caller{Role: domain.RoleTeacher, UserID: otherTeacherID, BranchIDs: []uuid.UUID{branchID}},
			lesson: &domain.Lesson{
				ID:        lessonID,
				BranchID:  branchID,
				TeacherID: teacherID,
				SubjectID: subjectID,
				StudentID: &studentID,
				GroupID:   &groupID,
				Date:      date,
				StartTime: start,
				EndTime:   end,
				Status:    domain.LessonStatusScheduled,
			},
			expectedErr: domain.ErrBranchAccessDenied,
		},
		{
			name:        "lesson_not_found",
			caller:      domain.Caller{Role: domain.RoleSuperadmin},
			repoErr:     domain.ErrNotFound,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.ScheduleRepository)
			if tt.repoErr != nil {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(nil, tt.repoErr).Once()
			} else {
				repo.On("GetLessonByID", mock.Anything, lessonID).Return(tt.lesson, nil).Once()
			}

			uc := NewUseCase(repo)
			res, err := uc.Execute(context.Background(), tt.caller, Request{ID: lessonID})

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, lessonID, res.ID)
			assert.Equal(t, branchID, res.BranchID)
			assert.Equal(t, teacherID, res.TeacherID)
			assert.Equal(t, subjectID, res.SubjectID)
			assert.Equal(t, date.Format(dateLayout), res.Date)
			assert.Equal(t, start.Format(timeLayout), res.StartTime)
			assert.Equal(t, end.Format(timeLayout), res.EndTime)
			assert.Equal(t, string(domain.LessonStatusScheduled), res.Status)

			repo.AssertExpectations(t)
		})
	}
}
