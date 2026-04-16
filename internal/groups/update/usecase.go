package update

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	groupRepo domain.GroupRepository
}

func NewUseCase(gr domain.GroupRepository) *UseCase {
	return &UseCase{
		groupRepo: gr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, groupID uuid.UUID, req Request) (Response, error) {
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return Response{}, err
	}
	if group.Status == domain.StatusArchived {
		return Response{}, domain.ErrCannotEditArchived
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, group.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	updatedGroup, err := uc.groupRepo.UpdateName(ctx, groupID, req.Name)
	if err != nil {
		return Response{}, err
	}

	return Response{
		ID:       updatedGroup.ID,
		BranchID: updatedGroup.BranchID,
		Name:     updatedGroup.Name,
		Status:   updatedGroup.Status,
	}, nil
}
