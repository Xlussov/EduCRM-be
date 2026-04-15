package create

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	branchRepo domain.BranchRepository
	userRepo   domain.UserRepository
	txManager  domain.TxManager
}

func NewUseCase(br domain.BranchRepository, ur domain.UserRepository, tm domain.TxManager) *UseCase {
	return &UseCase{
		branchRepo: br,
		userRepo:   ur,
		txManager:  tm,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (Response, error) {
	branch := &domain.Branch{
		Name:    req.Name,
		Address: req.Address,
		City:    req.City,
		Status:  domain.StatusActive,
	}

	err := uc.txManager.Transaction(ctx, func(txCtx context.Context) error {
		if err := uc.branchRepo.Create(txCtx, branch); err != nil {
			return err
		}

		if err := uc.userRepo.AssignToBranches(txCtx, caller.UserID, []uuid.UUID{branch.ID}); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return Response{}, err
	}

	return Response{ID: branch.ID}, nil
}
