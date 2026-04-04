package archive

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchAccessDenied = errors.New("branch access denied")
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

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, planID uuid.UUID, req Request) (Response, error) {
	plan, err := uc.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return Response{}, err
	}

	if role == "ADMIN" {
		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return Response{}, err
		}
		hasAccess := false
		for _, bid := range branchIDs {
			if bid == plan.BranchID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return Response{}, ErrBranchAccessDenied
		}
	}

	if err := uc.planRepo.UpdatePlanStatus(ctx, planID, domain.EntityStatus(req.Status)); err != nil {
		return Response{}, err
	}

	return Response{Message: "success"}, nil
}
