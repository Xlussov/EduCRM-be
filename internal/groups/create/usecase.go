package create

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	groupRepo domain.GroupRepository
	userRepo  domain.UserRepository
}

func NewUseCase(gr domain.GroupRepository, ur domain.UserRepository) *UseCase {
	return &UseCase{
		groupRepo: gr,
		userRepo:  ur,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (Response, error) {
	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	isActive, err := uc.userRepo.IsBranchActive(ctx, req.BranchID)
	if err != nil {
		return Response{}, err
	}
	if !isActive {
		return Response{}, domain.ErrArchivedReference
	}

	group := &domain.Group{
		BranchID: req.BranchID,
		Name:     req.Name,
	}

	if err := uc.groupRepo.Create(ctx, group); err != nil {
		return Response{}, err
	}

	return Response{ID: group.ID}, nil
}
