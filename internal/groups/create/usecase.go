package create

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
	groupRepo domain.GroupRepository
	userRepo  domain.UserRepository
}

func NewUseCase(gr domain.GroupRepository, ur domain.UserRepository) *UseCase {
	return &UseCase{
		groupRepo: gr,
		userRepo:  ur,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, req Request) (Response, error) {
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

	group := &domain.Group{
		BranchID: req.BranchID,
		Name:     req.Name,
	}

	if err := uc.groupRepo.Create(ctx, group); err != nil {
		return Response{}, err
	}

	return Response{ID: group.ID}, nil
}
