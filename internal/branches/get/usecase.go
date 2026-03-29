package get

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
	b, err := uc.branchRepo.GetByID(ctx, branchID)
	if err != nil {
		return Response{}, err
	}

	return Response{
		ID:      b.ID,
		Name:    b.Name,
		Address: b.Address,
		City:    b.City,
		Status:  string(b.Status),
	}, nil
}
