package unarchive

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	userRepo domain.UserRepository
}

func NewUseCase(ur domain.UserRepository) *UseCase {
	return &UseCase{userRepo: ur}
}

func (uc *UseCase) Execute(ctx context.Context, _ domain.Caller, adminID uuid.UUID) (Response, error) {
	user, err := uc.userRepo.GetByID(ctx, adminID)
	if err != nil {
		return Response{}, err
	}

	if user.Role != domain.RoleAdmin {
		return Response{}, domain.ErrNotFound
	}

	if user.IsActive {
		return Response{}, domain.ErrAlreadyActive
	}

	if err := uc.userRepo.UpdateUserStatus(ctx, adminID, true); err != nil {
		return Response{}, err
	}

	return Response{Message: "success"}, nil
}
