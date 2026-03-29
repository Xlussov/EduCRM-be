package update

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

func (uc *UseCase) Execute(ctx context.Context, branchID uuid.UUID, req Request) error {
	branch := &domain.Branch{
		ID:      branchID,
		Name:    req.Name,
		Address: req.Address,
		City:    req.City,
	}

	return uc.branchRepo.Update(ctx, branch)
}
