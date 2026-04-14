package login

import (
	"context"
	"errors"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/hash"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	userRepo   domain.UserRepository
	authRepo   domain.AuthRepository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewUseCase(ur domain.UserRepository, ar domain.AuthRepository, secret string, accessTTL, refreshTTL time.Duration) *UseCase {
	return &UseCase{
		userRepo:   ur,
		authRepo:   ar,
		jwtSecret:  secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) (Response, error) {
	user, err := uc.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return Response{}, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return Response{}, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return Response{}, errors.New("user is not active")
	}

	branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, user.ID)
	if err != nil {
		return Response{}, err
	}

	bStr := make([]string, len(branchIDs))
	for i, b := range branchIDs {
		bStr[i] = b.String()
	}

	accessToken, err := uc.generateToken(user.ID.String(), string(user.Role), bStr, uc.accessTTL)
	if err != nil {
		return Response{}, err
	}

	refreshTokenID := uuid.New()
	refreshTokenStr, err := uc.generateToken(user.ID.String(), string(user.Role), bStr, uc.refreshTTL)
	if err != nil {
		return Response{}, err
	}

	hashStr := hash.SHA256Token(refreshTokenStr)

	err = uc.authRepo.SaveRefreshToken(ctx, refreshTokenID, user.ID, hashStr, time.Now().Add(uc.refreshTTL))
	if err != nil {
		return Response{}, err
	}

	return Response{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User: UserDTO{
			ID:       user.ID,
			Role:     string(user.Role),
			Branches: branchIDs,
		},
	}, nil
}

func (uc *UseCase) generateToken(userID, role string, branches []string, d time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"role":       role,
		"branch_ids": branches,
		"exp":        jwt.NewNumericDate(time.Now().Add(d)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
