package me

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	userRepo domain.UserRepository
}

func NewUseCase(ur domain.UserRepository) *UseCase {
	return &UseCase{
		userRepo: ur,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID) (Response, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return Response{}, err
	}

	return Response{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      user.Role,
	}, nil
}
