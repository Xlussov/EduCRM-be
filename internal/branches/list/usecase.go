package list

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	branchRepo domain.BranchRepository
}

func NewUseCase(br domain.BranchRepository) *UseCase {
	return &UseCase{
		branchRepo: br,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) ([]BranchResponse, error) {
	status, err := domain.ParseEntityStatus(req.Status)
	if err != nil {
		return nil, err
	}

	var branches []*domain.Branch

	if caller.Role == domain.RoleSuperadmin {
		branches, err = uc.branchRepo.GetAll(ctx, status)
	} else {
		branches, err = uc.branchRepo.GetByUserID(ctx, caller.UserID, status)
	}

	if err != nil {
		return nil, err
	}

	res := make([]BranchResponse, 0, len(branches))
	for _, b := range branches {
		res = append(res, BranchResponse{
			ID:      b.ID,
			Name:    b.Name,
			Address: b.Address,
			City:    b.City,
			Status:  string(b.Status),
		})
	}

	return res, nil
}
