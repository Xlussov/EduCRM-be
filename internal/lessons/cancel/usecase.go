package cancel

import (
	"context"

	"github.com/google/uuid"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	scheduleRepo domain.ScheduleRepository
}

func NewUseCase(sr domain.ScheduleRepository) *UseCase {
	return &UseCase{
		scheduleRepo: sr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, id uuid.UUID) (Response, error) {
	lesson, err := uc.scheduleRepo.GetLessonByID(ctx, id)
	if err != nil {
		return Response{}, err
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, lesson.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if lesson.Status == domain.LessonStatusCancelled || lesson.Status == domain.LessonStatusCompleted {
		return Response{}, domain.ErrInvalidInput
	}

	if err := uc.scheduleRepo.UpdateLessonStatus(ctx, id, domain.LessonStatusCancelled); err != nil {
		return Response{}, err
	}

	return Response{
		Message: "success",
	}, nil
}
