package mark

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

var (
	ErrAttendanceRequired = errors.New("attendance is required")
	ErrLessonCancelled    = errors.New("lesson is cancelled")
	ErrLessonInFuture     = errors.New("lesson is in the future")
	ErrStudentNotInLesson = errors.New("student is not in lesson")
	ErrDuplicateStudentID = errors.New("duplicate student_id")
	ErrIsPresentRequired  = errors.New("is_present is required")
)

type UseCase struct {
	scheduleRepo   domain.ScheduleRepository
	attendanceRepo domain.AttendanceRepository
	txManager      domain.TxManager
}

func NewUseCase(sr domain.ScheduleRepository, ar domain.AttendanceRepository, tm domain.TxManager) *UseCase {
	return &UseCase{scheduleRepo: sr, attendanceRepo: ar, txManager: tm}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, lessonID uuid.UUID, req Request) (Response, error) {
	if len(req.Attendance) == 0 {
		return Response{}, ErrAttendanceRequired
	}

	lesson, err := uc.scheduleRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return Response{}, fmt.Errorf("get lesson: %w", err)
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, lesson.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if caller.Role == domain.RoleTeacher && lesson.TeacherID != caller.UserID {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if lesson.Status == domain.LessonStatusCancelled {
		return Response{}, ErrLessonCancelled
	}

	if isFutureLesson(lesson.Date) {
		return Response{}, ErrLessonInFuture
	}

	expected, err := uc.attendanceRepo.GetLessonAttendance(ctx, lesson.ID)
	if err != nil {
		return Response{}, fmt.Errorf("get lesson attendance: %w", err)
	}

	expectedSet := make(map[uuid.UUID]struct{}, len(expected))
	for _, item := range expected {
		expectedSet[item.StudentID] = struct{}{}
	}

	seen := make(map[uuid.UUID]struct{}, len(req.Attendance))
	records := make([]domain.Attendance, 0, len(req.Attendance))

	for _, item := range req.Attendance {
		if item.IsPresent == nil {
			return Response{}, ErrIsPresentRequired
		}
		if _, ok := expectedSet[item.StudentID]; !ok {
			return Response{}, ErrStudentNotInLesson
		}
		if _, ok := seen[item.StudentID]; ok {
			return Response{}, ErrDuplicateStudentID
		}
		seen[item.StudentID] = struct{}{}

		records = append(records, domain.Attendance{
			LessonID:  lesson.ID,
			StudentID: item.StudentID,
			IsPresent: *item.IsPresent,
			Notes:     item.Notes,
		})
	}

	err = uc.txManager.Transaction(ctx, func(txCtx context.Context) error {
		if err := uc.attendanceRepo.UpsertAttendance(txCtx, records); err != nil {
			return fmt.Errorf("upsert attendance: %w", err)
		}
		if err := uc.scheduleRepo.UpdateLessonStatus(txCtx, lesson.ID, domain.LessonStatusCompleted); err != nil {
			return fmt.Errorf("update lesson status: %w", err)
		}
		return nil
	})
	if err != nil {
		return Response{}, err
	}

	updated, err := uc.attendanceRepo.GetLessonAttendance(ctx, lesson.ID)
	if err != nil {
		return Response{}, fmt.Errorf("get lesson attendance: %w", err)
	}

	return Response{Attendance: toResponse(updated)}, nil
}

func isFutureLesson(date time.Time) bool {
	location := date.Location()
	now := time.Now().In(location)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	return date.After(today)
}

func toResponse(items []domain.LessonAttendanceStudent) []StudentAttendance {
	res := make([]StudentAttendance, 0, len(items))
	for _, item := range items {
		res = append(res, StudentAttendance{
			StudentID: item.StudentID,
			FirstName: item.FirstName,
			LastName:  item.LastName,
			Status:    string(item.Status),
			IsPresent: item.IsPresent,
			Notes:     item.Notes,
		})
	}
	return res
}
