package archive

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	planRepo domain.SubscriptionRepository
}

func NewUseCase(pr domain.SubscriptionRepository) *UseCase {
	return &UseCase{
		planRepo: pr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, planID uuid.UUID) (Response, error) {
	plan, err := uc.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return Response{}, err
	}

	if domain.RequiresBranchAccess(caller.Role) {
		if !domain.HasBranchAccess(caller.BranchIDs, plan.BranchID) {
			return Response{}, domain.ErrBranchAccessDenied
		}
	}

	if err := uc.planRepo.UpdatePlanStatus(ctx, planID, domain.StatusArchived); err != nil {
		return Response{}, err
	}

	return Response{Message: "success"}, nil
}
