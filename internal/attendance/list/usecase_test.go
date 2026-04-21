package list

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/internal/domain/mocks"
)

func TestUseCase_Execute(t *testing.T) {
	lessonID := uuid.New()
	branchID := uuid.New()
	teacherID := uuid.New()
	studentID := uuid.New()

	lesson := &domain.Lesson{ID: lessonID, BranchID: branchID, TeacherID: teacherID}
	present := true
	note := "ok"

	attendance := []domain.LessonAttendanceStudent{
		{
			StudentID: studentID,
			FirstName: "Ada",
			LastName:  "Lovelace",
			Status:    domain.StatusActive,
			IsPresent: &present,
			Notes:     &note,
		},
	}

	errDB := errors.New("db error")

	tests := []struct {
		name        string
		caller      domain.Caller
		setupMocks  func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository)
		expectedErr error
		attCalled   bool
	}{
		{
			name:   "success admin",
			caller: domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(lesson, nil).Once()
				ar.On("GetLessonAttendance", mock.Anything, lessonID).Return(attendance, nil).Once()
			},
			attCalled: true,
		},
		{
			name:   "admin branch access denied",
			caller: domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}},
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(lesson, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
			attCalled:   false,
		},
		{
			name:   "teacher mismatch",
			caller: domain.Caller{UserID: uuid.New(), Role: domain.RoleTeacher, BranchIDs: []uuid.UUID{branchID}},
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(lesson, nil).Once()
			},
			expectedErr: domain.ErrBranchAccessDenied,
			attCalled:   false,
		},
		{
			name:   "lesson not found",
			caller: domain.Caller{UserID: uuid.New(), Role: domain.RoleSuperadmin},
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(nil, domain.ErrNotFound).Once()
			},
			expectedErr: domain.ErrNotFound,
			attCalled:   false,
		},
		{
			name:   "attendance repo error",
			caller: domain.Caller{UserID: uuid.New(), Role: domain.RoleSuperadmin},
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(lesson, nil).Once()
				ar.On("GetLessonAttendance", mock.Anything, lessonID).Return(nil, errDB).Once()
			},
			expectedErr: errDB,
			attCalled:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := new(mocks.ScheduleRepository)
			ar := new(mocks.AttendanceRepository)

			if tt.setupMocks != nil {
				tt.setupMocks(sr, ar)
			}

			uc := NewUseCase(sr, ar)
			res, err := uc.Execute(context.Background(), tt.caller, Request{ID: lessonID})

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				assert.Empty(t, res.Attendance)
			} else {
				require.NoError(t, err)
				require.Len(t, res.Attendance, 1)
				assert.Equal(t, studentID, res.Attendance[0].StudentID)
				assert.Equal(t, "Ada", res.Attendance[0].FirstName)
				assert.Equal(t, "Lovelace", res.Attendance[0].LastName)
				assert.Equal(t, string(domain.StatusActive), res.Attendance[0].Status)
				assert.NotNil(t, res.Attendance[0].IsPresent)
			}

			if !tt.attCalled {
				ar.AssertNotCalled(t, "GetLessonAttendance", mock.Anything, lessonID)
			}

			sr.AssertExpectations(t)
			ar.AssertExpectations(t)
		})
	}
}
