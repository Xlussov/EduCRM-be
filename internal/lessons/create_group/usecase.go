package create_group

import (
	"context"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	scheduleRepo domain.ScheduleRepository
	groupRepo    domain.GroupRepository
}

func NewUseCase(sr domain.ScheduleRepository, gr domain.GroupRepository) *UseCase {
	return &UseCase{
		scheduleRepo: sr,
		groupRepo:    gr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (Response, error) {
	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
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

	hasTeacherConflict, err := uc.scheduleRepo.CheckTeacherConflict(ctx, req.TeacherID, date, start, end)
	if err != nil {
		return Response{}, err
	}
	if hasTeacherConflict {
		return Response{}, domain.ErrTeacherScheduleConflict
	}

	studentIDs, err := uc.groupRepo.GetActiveStudentIDs(ctx, req.GroupID)
	if err != nil {
		return Response{}, err
	}

	for _, studentID := range studentIDs {
		hasConflict, err := uc.scheduleRepo.CheckStudentConflict(ctx, studentID, date, start, end)
		if err != nil {
			return Response{}, err
		}
		if hasConflict {
			return Response{}, domain.ErrStudentScheduleConflict
		}
	}

	lesson := &domain.Lesson{
		BranchID:  req.BranchID,
		TeacherID: req.TeacherID,
		SubjectID: req.SubjectID,
		GroupID:   &req.GroupID,
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
