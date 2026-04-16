package archive

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var ()

type UseCase struct {
	groupRepo domain.GroupRepository
}

func NewUseCase(gr domain.GroupRepository) *UseCase {
	return &UseCase{
		groupRepo: gr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, groupID uuid.UUID) (Response, error) {
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return Response{}, err
	}

	if group.Status == domain.StatusArchived {
		return Response{}, domain.ErrAlreadyArchived
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, group.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	err = uc.groupRepo.UpdateStatus(ctx, groupID, domain.StatusArchived)
	if err != nil {
		return Response{}, err
	}

	return Response{Message: "success"}, nil
}
