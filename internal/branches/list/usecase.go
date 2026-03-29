package list

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

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string) ([]BranchResponse, error) {
	var branches []*domain.Branch
	var err error

	if role == "SUPERADMIN" {
		branches, err = uc.branchRepo.GetAll(ctx)
	} else {
		// ADMIN
		branches, err = uc.branchRepo.GetByUserID(ctx, userID)
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
