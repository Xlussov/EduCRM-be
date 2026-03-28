package create

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	branchRepo domain.BranchRepository
	userRepo   domain.UserRepository
}

func NewUseCase(br domain.BranchRepository, ur domain.UserRepository) *UseCase {
	return &UseCase{
		branchRepo: br,
		userRepo:   ur,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, req Request) (Response, error) {
	branch := &domain.Branch{
		Name:    req.Name,
		Address: req.Address,
		City:    req.City,
		Status:  domain.StatusActive,
	}

	if err := uc.branchRepo.Create(ctx, branch); err != nil {
		return Response{}, err
	}

	// Привязываем создавшего ADMIN/SUPERADMIN к новому филиалу
	if err := uc.userRepo.AssignToBranches(ctx, userID, []uuid.UUID{branch.ID}); err != nil {
		return Response{}, err
	}

	return Response{ID: branch.ID}, nil
}
