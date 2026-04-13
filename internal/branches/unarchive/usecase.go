package unarchive

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

func (uc *UseCase) Execute(ctx context.Context, branchID uuid.UUID) (Response, error) {
	branch, err := uc.branchRepo.GetByID(ctx, branchID)
	if err != nil {
		return Response{}, err
	}

	if branch.Status == domain.StatusActive {
		return Response{}, domain.ErrAlreadyActive
	}
	err = uc.branchRepo.UpdateStatus(ctx, branchID, domain.StatusActive)
	if err != nil {
		return Response{}, err
	}
	return Response{Message: "success"}, nil
}
