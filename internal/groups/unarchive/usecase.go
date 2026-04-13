package unarchive

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var ()

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

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, groupID uuid.UUID) (Response, error) {
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return Response{}, err
	}

	if group.Status == domain.StatusActive {
		return Response{}, domain.ErrAlreadyActive
	}

	if role == "ADMIN" {
		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return Response{}, err
		}
		hasAccess := false
		for _, bid := range branchIDs {
			if bid == group.BranchID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return Response{}, domain.ErrBranchAccessDenied
		}
	}

	err = uc.groupRepo.UpdateStatus(ctx, groupID, domain.StatusActive)
	if err != nil {
		return Response{}, err
	}

	return Response{Message: "success"}, nil
}
