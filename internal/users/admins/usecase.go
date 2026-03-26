package admins

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	userRepo domain.UserRepository
}

func NewUseCase(ur domain.UserRepository) *UseCase {
	return &UseCase{userRepo: ur}
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

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return Response{}, err
	}

	if len(req.BranchIDs) > 0 {
		if err := uc.userRepo.AssignToBranches(ctx, user.ID, req.BranchIDs); err != nil {
			return Response{}, err
		}
	}

	return Response{ID: user.ID}, nil
}
