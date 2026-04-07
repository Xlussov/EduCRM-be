package admins

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	userRepo  domain.UserRepository
	txManager domain.TxManager
}

func NewUseCase(ur domain.UserRepository, tm domain.TxManager) *UseCase {
	return &UseCase{userRepo: ur, txManager: tm}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) (Response, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return Response{}, err
	}

	user := &domain.User{
		Phone:        req.Phone,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         domain.RoleAdmin,
	}
	err = uc.txManager.Transaction(ctx, func(txCtx context.Context) error {

		if err := uc.userRepo.Create(txCtx, user); err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				return domain.ErrPhoneAlreadyExists
			}
			return err
		}

		if len(req.BranchIDs) > 0 {
			activeCount, err := uc.userRepo.CountActiveBranchesByIDs(txCtx, req.BranchIDs)
			if err != nil {
				return err
			}
			if activeCount != len(req.BranchIDs) {
				return domain.ErrArchivedReference
			}

			if err := uc.userRepo.AssignToBranches(txCtx, user.ID, req.BranchIDs); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return Response{}, err
	}

	return Response{ID: user.ID}, nil
}
