package list

import (
	"context"
	"fmt"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	scheduleRepo   domain.ScheduleRepository
	attendanceRepo domain.AttendanceRepository
}

func NewUseCase(sr domain.ScheduleRepository, ar domain.AttendanceRepository) *UseCase {
	return &UseCase{scheduleRepo: sr, attendanceRepo: ar}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (Response, error) {
	lesson, err := uc.scheduleRepo.GetLessonByID(ctx, req.ID)
	if err != nil {
		return Response{}, fmt.Errorf("get lesson: %w", err)
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, lesson.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if caller.Role == domain.RoleTeacher && lesson.TeacherID != caller.UserID {
		return Response{}, domain.ErrBranchAccessDenied
	}

	attendance, err := uc.attendanceRepo.GetLessonAttendance(ctx, lesson.ID)
	if err != nil {
		return Response{}, fmt.Errorf("get lesson attendance: %w", err)
	}

	return Response{Attendance: toResponse(attendance)}, nil
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
