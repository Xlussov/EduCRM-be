package list

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchIDRequired = errors.New("branch_id is required")
)

type UseCase struct {
	planRepo domain.SubscriptionRepository
}

func NewUseCase(pr domain.SubscriptionRepository) *UseCase {
	return &UseCase{
		planRepo: pr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) ([]PlanResponse, error) {
	if req.BranchID == uuid.Nil {
		return nil, ErrBranchIDRequired
	}

	if domain.RequiresBranchAccess(caller.Role) {
		if !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
			return nil, domain.ErrBranchAccessDenied
		}
	}

	plans, err := uc.planRepo.GetPlansByBranchID(ctx, req.BranchID)
	if err != nil {
		return nil, err
	}

	var res []PlanResponse
	for _, p := range plans {
		subjects := make([]Subject, 0, len(p.Subjects))
		for _, s := range p.Subjects {
			subjects = append(subjects, Subject{
				ID:   s.ID,
				Name: s.Name,
			})
		}

		grids := make([]PricingGrid, 0, len(p.PricingGrid))
		for _, g := range p.PricingGrid {
			grids = append(grids, PricingGrid{
				Lessons: g.LessonsPerMonth,
				Price:   g.PricePerLesson,
			})
		}

		res = append(res, PlanResponse{
			ID:          p.ID,
			Name:        p.Name,
			Type:        string(p.Type),
			Subjects:    subjects,
			PricingGrid: grids,
			Status:      string(p.Status),
		})
	}

	if res == nil {
		res = []PlanResponse{}
	}

	return res, nil
}
