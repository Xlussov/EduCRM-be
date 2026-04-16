package list

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchIDRequired = errors.New("branch_id is required")
)

type UseCase struct {
	groupRepo domain.GroupRepository
}

func NewUseCase(gr domain.GroupRepository) *UseCase {
	return &UseCase{
		groupRepo: gr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (Response, error) {
	if req.BranchID == uuid.Nil {
		return Response{}, ErrBranchIDRequired
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	status, err := domain.ParseEntityStatus(req.Status)
	if err != nil {
		return Response{}, err
	}

	var groups []*domain.GroupWithCount
	if caller.Role == domain.RoleTeacher {
		groups, err = uc.groupRepo.GetByBranchIDAndTeacherID(ctx, req.BranchID, caller.UserID, status)
	} else {
		groups, err = uc.groupRepo.GetByBranchID(ctx, req.BranchID, status)
	}
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
