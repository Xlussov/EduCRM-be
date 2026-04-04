package list

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchAccessDenied = errors.New("branch access denied")
	ErrBranchIDRequired   = errors.New("branch_id is required")
)

type UseCase struct {
	planRepo domain.SubscriptionRepository
	userRepo domain.UserRepository
}

func NewUseCase(pr domain.SubscriptionRepository, ur domain.UserRepository) *UseCase {
	return &UseCase{
		planRepo: pr,
		userRepo: ur,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, req Request) ([]PlanResponse, error) {
	if req.BranchID == uuid.Nil {
		return nil, ErrBranchIDRequired
	}

	if role == "ADMIN" {
		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return nil, err
		}
		hasAccess := false
		for _, bid := range branchIDs {
			if bid == req.BranchID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return nil, ErrBranchAccessDenied
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
		})
	}

	if res == nil {
		res = []PlanResponse{}
	}

	return res, nil
}
