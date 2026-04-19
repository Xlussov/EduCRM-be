package create_individual

import (
	"context"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	scheduleRepo domain.ScheduleRepository
	userRepo     domain.UserRepository
}

func NewUseCase(sr domain.ScheduleRepository, ur domain.UserRepository) *UseCase {
	return &UseCase{
		scheduleRepo: sr,
		userRepo:     ur,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (Response, error) {
	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	teacherInBranch, err := uc.userRepo.CheckTeacherInBranch(ctx, req.TeacherID, req.BranchID)
	if err != nil {
		return Response{}, err
	}
	if !teacherInBranch {
		return Response{}, domain.ErrTeacherNotInBranch
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return Response{}, domain.ErrInvalidInput
	}

	start, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return Response{}, domain.ErrInvalidInput
	}

	end, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return Response{}, domain.ErrInvalidInput
	}

	if !start.Before(end) {
		return Response{}, domain.ErrInvalidInput
	}

	if req.StudentID != nil {
		hasStudentConflict, err := uc.scheduleRepo.CheckStudentConflict(ctx, *req.StudentID, date, start, end)
		if err != nil {
			return Response{}, err
		}
		if hasStudentConflict {
			return Response{}, domain.ErrStudentScheduleConflict
		}
	}

	hasTeacherConflict, err := uc.scheduleRepo.CheckTeacherConflict(ctx, req.TeacherID, date, start, end)
	if err != nil {
		return Response{}, err
	}
	if hasTeacherConflict {
		return Response{}, domain.ErrTeacherScheduleConflict
	}

	lesson := &domain.Lesson{
		BranchID:  req.BranchID,
		TeacherID: req.TeacherID,
		SubjectID: req.SubjectID,
		StudentID: req.StudentID,
		Date:      date,
		StartTime: start,
		EndTime:   end,
		Status:    domain.LessonStatusScheduled,
	}

	if err := uc.scheduleRepo.CreateLesson(ctx, lesson); err != nil {
		return Response{}, err
	}

	return Response{
		ID:     lesson.ID,
		Status: string(lesson.Status),
	}, nil
}
