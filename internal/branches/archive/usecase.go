package archive

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	branchRepo domain.BranchRepository
}

func NewUseCase(br domain.BranchRepository) *UseCase {
	return &UseCase{
		branchRepo: br,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, branchID uuid.UUID) (Response, error) {
	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, branchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	branch, err := uc.branchRepo.GetByID(ctx, branchID)
	if err != nil {
		return Response{}, err
	}

	if branch.Status == domain.StatusArchived {
		return Response{}, domain.ErrAlreadyArchived
	}

	err = uc.branchRepo.UpdateStatus(ctx, branchID, domain.StatusArchived)
	if err != nil {
		return Response{}, err
	}
	return Response{Message: "success"}, nil
}
