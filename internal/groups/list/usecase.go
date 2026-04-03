package list

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchIDRequired   = errors.New("branch_id is required")
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
	if req.BranchID == uuid.Nil {
		return Response{}, ErrBranchIDRequired
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

	groups, err := uc.groupRepo.GetByBranchID(ctx, req.BranchID)
	if err != nil {
		return Response{}, err
	}

	var res []GroupResponse
	for _, g := range groups {
		res = append(res, GroupResponse{
			ID:            g.ID,
			Name:          g.Name,
			StudentsCount: g.StudentsCount,
			Status:        g.Status,
		})
	}

	if res == nil {
		res = []GroupResponse{}
	}

	return Response{Groups: res}, nil
}
