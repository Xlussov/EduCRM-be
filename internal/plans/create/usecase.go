package create

import (
	"context"
	"errors"
	"fmt"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

var (
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

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (Response, error) {
	if domain.RequiresBranchAccess(caller.Role) {
		if !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
			return Response{}, domain.ErrBranchAccessDenied
		}
	}

	isActive, err := uc.userRepo.IsBranchActive(ctx, req.BranchID)
	if err != nil {
		return Response{}, err
	}
	if !isActive {
		return Response{}, domain.ErrArchivedReference
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
