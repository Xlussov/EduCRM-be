package create

import (
	"context"
	"errors"
	"fmt"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchAccessDenied    = errors.New("branch access denied")
	ErrSubjectBranchMismatch = errors.New("all subjects must belong to the plan branch")
)

type UseCase struct {
	txManager domain.TxManager
	planRepo  domain.SubscriptionRepository
	userRepo  domain.UserRepository
}

func NewUseCase(tx domain.TxManager, pr domain.SubscriptionRepository, ur domain.UserRepository) *UseCase {
	return &UseCase{
		txManager: tx,
		planRepo:  pr,
		userRepo:  ur,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, req Request) (Response, error) {
	isActive, err := uc.userRepo.IsBranchActive(ctx, req.BranchID)
	if err != nil {
		return Response{}, err
	}
	if !isActive {
		return Response{}, domain.ErrArchivedReference
	}

	if role == "ADMIN" {
		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return Response{}, err
		}
		hasAccess := false
		for _, bid := range branchIDs {
			if bid == req.BranchID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return Response{}, ErrBranchAccessDenied
		}
	}

	matchedSubjectsCount, err := uc.planRepo.CountSubjectsInBranch(ctx, req.BranchID, req.SubjectIDs)
	if err != nil {
		return Response{}, err
	}
	if matchedSubjectsCount != len(req.SubjectIDs) {
		return Response{}, ErrSubjectBranchMismatch
	}

	plan := &domain.Plan{
		BranchID: req.BranchID,
		Name:     req.Name,
		Type:     domain.PlanType(req.Type),
		Status:   domain.StatusActive,
	}

	var grids []*domain.PricingGrid
	for _, g := range req.PricingGrid {
		grids = append(grids, &domain.PricingGrid{
			LessonsPerMonth: g.Lessons,
			PricePerLesson:  g.Price,
		})
	}

	err = uc.txManager.Transaction(ctx, func(txCtx context.Context) error {
		if err := uc.planRepo.CreatePlan(txCtx, plan, req.SubjectIDs, grids); err != nil {
			return fmt.Errorf("failed to create plan: %w", err)
		}
		return nil
	})

	if err != nil {
		return Response{}, err
	}

	return Response{ID: plan.ID}, nil
}
