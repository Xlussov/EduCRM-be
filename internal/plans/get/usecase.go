package get

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	planRepo domain.SubscriptionRepository
}

func NewUseCase(pr domain.SubscriptionRepository) *UseCase {
	return &UseCase{planRepo: pr}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, planID uuid.UUID) (Response, error) {
	plan, err := uc.planRepo.GetPlanDetailsByID(ctx, planID)
	if err != nil {
		return Response{}, err
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, plan.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	subjects := make([]Subject, 0, len(plan.Subjects))
	for _, s := range plan.Subjects {
		subjects = append(subjects, Subject{ID: s.ID, Name: s.Name})
	}

	grids := make([]PricingGrid, 0, len(plan.PricingGrid))
	for _, g := range plan.PricingGrid {
		grids = append(grids, PricingGrid{Lessons: g.LessonsPerMonth, Price: g.PricePerLesson})
	}

	return Response{
		ID:          plan.ID,
		Name:        plan.Name,
		Type:        string(plan.Type),
		Subjects:    subjects,
		PricingGrid: grids,
	}, nil
}
