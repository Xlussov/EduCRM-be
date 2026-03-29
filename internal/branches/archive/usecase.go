package archive

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
	if err := uc.branchRepo.UpdateStatus(ctx, branchID, domain.StatusArchived); err != nil {
		return Response{}, err
	}
	return Response{Message: "success"}, nil
}
