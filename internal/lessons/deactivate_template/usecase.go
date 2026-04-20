package deactivate_template

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	scheduleRepo domain.ScheduleRepository
	txManager    domain.TxManager
}

func NewUseCase(sr domain.ScheduleRepository, tm domain.TxManager) *UseCase {
	return &UseCase{scheduleRepo: sr, txManager: tm}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, templateID uuid.UUID) (Response, error) {
	template, err := uc.scheduleRepo.GetTemplateByID(ctx, templateID)
	if err != nil {
		return Response{}, fmt.Errorf("get template: %w", err)
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, template.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if !template.IsActive {
		return Response{}, domain.ErrTemplateNotActive
	}

	err = uc.txManager.Transaction(ctx, func(txCtx context.Context) error {
		if err := uc.scheduleRepo.DeactivateTemplate(txCtx, templateID); err != nil {
			return fmt.Errorf("deactivate template: %w", err)
		}
		if err := uc.scheduleRepo.CancelFutureLessonsByTemplate(txCtx, templateID); err != nil {
			return fmt.Errorf("cancel future lessons: %w", err)
		}
		return nil
	})
	if err != nil {
		return Response{}, err
	}

	return Response{Message: "success"}, nil
}
