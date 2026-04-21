package logout

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/hash"
)

type UseCase struct {
	authRepo domain.AuthRepository
}

func NewUseCase(ar domain.AuthRepository) *UseCase {
	return &UseCase{
		authRepo: ar,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) (Response, error) {
	tokenHash := hash.SHA256Token(req.RefreshToken)

	token, err := uc.authRepo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return Response{Message: "Successfully logged out"}, nil
		}
		return Response{}, err
	}

	if token.IsRevoked {
		return Response{Message: "Successfully logged out"}, nil
	}

	err = uc.authRepo.RevokeRefreshToken(ctx, token.ID)
	if err != nil {
		return Response{}, err
	}

	return Response{Message: "Successfully logged out"}, nil
}
