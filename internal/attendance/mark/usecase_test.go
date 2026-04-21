package mark

import (
	"context"
	"errors"
	"testing"
	"time"

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
	studentA := uuid.New()
	studentB := uuid.New()
	otherStudent := uuid.New()

	now := time.Now()
	pastDate := now.AddDate(0, 0, -1)
	futureDate := now.AddDate(0, 0, 1)

	makeLesson := func(date time.Time, status domain.LessonStatus) *domain.Lesson {
		return &domain.Lesson{
			ID:        lessonID,
			BranchID:  branchID,
			TeacherID: teacherID,
			Date:      date,
			Status:    status,
		}
	}

	present := true
	absent := false
	note := "note"

	expectedAttendance := []domain.LessonAttendanceStudent{
		{StudentID: studentA, FirstName: "Ada", LastName: "Lovelace", Status: domain.StatusActive},
		{StudentID: studentB, FirstName: "Alan", LastName: "Turing", Status: domain.StatusActive},
	}

	req := Request{
		Attendance: []AttendanceItem{
			{StudentID: studentA, IsPresent: &present, Notes: &note},
			{StudentID: studentB, IsPresent: &absent},
		},
	}

	errDB := errors.New("db error")

	tests := []struct {
		name        string
		caller      domain.Caller
		lesson      *domain.Lesson
		req         Request
		setupMocks  func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository)
		expectedErr error
		attCalls    int
		upsertCall  bool
		updateCall  bool
	}{
		{
			name:   "success admin",
			caller: domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			lesson: makeLesson(pastDate, domain.LessonStatusScheduled),
			req:    req,
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(makeLesson(pastDate, domain.LessonStatusScheduled), nil).Once()
				ar.On("GetLessonAttendance", mock.Anything, lessonID).Return(expectedAttendance, nil).Twice()
				ar.On("UpsertAttendance", mock.Anything, mock.MatchedBy(func(items []domain.Attendance) bool {
					if len(items) != 2 {
						return false
					}
					return items[0].LessonID == lessonID && items[1].LessonID == lessonID
				})).Return(nil).Once()
				sr.On("UpdateLessonStatus", mock.Anything, lessonID, domain.LessonStatusCompleted).Return(nil).Once()
			},
			attCalls:   2,
			upsertCall: true,
			updateCall: true,
		},
		{
			name:        "attendance required",
			caller:      domain.Caller{UserID: uuid.New(), Role: domain.RoleSuperadmin},
			lesson:      makeLesson(pastDate, domain.LessonStatusScheduled),
			req:         Request{},
			expectedErr: ErrAttendanceRequired,
		},
		{
			name:        "lesson cancelled",
			caller:      domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			lesson:      makeLesson(pastDate, domain.LessonStatusCancelled),
			req:         req,
			expectedErr: ErrLessonCancelled,
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(makeLesson(pastDate, domain.LessonStatusCancelled), nil).Once()
			},
		},
		{
			name:        "lesson in future",
			caller:      domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			lesson:      makeLesson(futureDate, domain.LessonStatusScheduled),
			req:         req,
			expectedErr: ErrLessonInFuture,
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(makeLesson(futureDate, domain.LessonStatusScheduled), nil).Once()
			},
		},
		{
			name:   "student not in lesson",
			caller: domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			lesson: makeLesson(pastDate, domain.LessonStatusScheduled),
			req: Request{Attendance: []AttendanceItem{
				{StudentID: otherStudent, IsPresent: &present},
			}},
			expectedErr: ErrStudentNotInLesson,
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(makeLesson(pastDate, domain.LessonStatusScheduled), nil).Once()
				ar.On("GetLessonAttendance", mock.Anything, lessonID).Return(expectedAttendance, nil).Once()
			},
			attCalls: 1,
		},
		{
			name:   "duplicate student",
			caller: domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			lesson: makeLesson(pastDate, domain.LessonStatusScheduled),
			req: Request{Attendance: []AttendanceItem{
				{StudentID: studentA, IsPresent: &present},
				{StudentID: studentA, IsPresent: &present},
			}},
			expectedErr: ErrDuplicateStudentID,
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(makeLesson(pastDate, domain.LessonStatusScheduled), nil).Once()
				ar.On("GetLessonAttendance", mock.Anything, lessonID).Return(expectedAttendance, nil).Once()
			},
			attCalls: 1,
		},
		{
			name:   "is_present required",
			caller: domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			lesson: makeLesson(pastDate, domain.LessonStatusScheduled),
			req: Request{Attendance: []AttendanceItem{
				{StudentID: studentA},
			}},
			expectedErr: ErrIsPresentRequired,
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(makeLesson(pastDate, domain.LessonStatusScheduled), nil).Once()
				ar.On("GetLessonAttendance", mock.Anything, lessonID).Return(expectedAttendance, nil).Once()
			},
			attCalls: 1,
		},
		{
			name:        "branch access denied",
			caller:      domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{uuid.New()}},
			lesson:      makeLesson(pastDate, domain.LessonStatusScheduled),
			req:         req,
			expectedErr: domain.ErrBranchAccessDenied,
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(makeLesson(pastDate, domain.LessonStatusScheduled), nil).Once()
			},
		},
		{
			name:        "upsert error",
			caller:      domain.Caller{UserID: uuid.New(), Role: domain.RoleAdmin, BranchIDs: []uuid.UUID{branchID}},
			lesson:      makeLesson(pastDate, domain.LessonStatusScheduled),
			req:         req,
			expectedErr: errDB,
			setupMocks: func(sr *mocks.ScheduleRepository, ar *mocks.AttendanceRepository) {
				sr.On("GetLessonByID", mock.Anything, lessonID).Return(makeLesson(pastDate, domain.LessonStatusScheduled), nil).Once()
				ar.On("GetLessonAttendance", mock.Anything, lessonID).Return(expectedAttendance, nil).Once()
				ar.On("UpsertAttendance", mock.Anything, mock.Anything).Return(errDB).Once()
			},
			attCalls:   1,
			upsertCall: true,
			updateCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := new(mocks.ScheduleRepository)
			ar := new(mocks.AttendanceRepository)
			tx := new(mocks.MockTxManager)

			if tt.setupMocks != nil {
				tt.setupMocks(sr, ar)
			}

			uc := NewUseCase(sr, ar, tx)
			res, err := uc.Execute(context.Background(), tt.caller, lessonID, tt.req)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				assert.Empty(t, res.Attendance)
			} else {
				require.NoError(t, err)
				require.Len(t, res.Attendance, len(expectedAttendance))
			}

			if tt.attCalls == 0 {
				ar.AssertNotCalled(t, "GetLessonAttendance", mock.Anything, lessonID)
			}

			if !tt.updateCall {
				sr.AssertNotCalled(t, "UpdateLessonStatus", mock.Anything, mock.Anything, mock.Anything)
			}

			sr.AssertExpectations(t)
			ar.AssertExpectations(t)
		})
	}
}
