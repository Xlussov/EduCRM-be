package refresh

import (
	"context"
	"errors"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/Xlussov/EduCRM-be/pkg/hash"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UseCase struct {
	userRepo  domain.UserRepository
	authRepo  domain.AuthRepository
	jwtSecret string
}

func NewUseCase(ur domain.UserRepository, ar domain.AuthRepository, secret string) *UseCase {
	return &UseCase{
		userRepo:  ur,
		authRepo:  ar,
		jwtSecret: secret,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) (Response, error) {
	hashStr := hash.SHA256Token(req.RefreshToken)
	tokenRecord, err := uc.authRepo.GetRefreshToken(ctx, hashStr)
	if err != nil {
		return Response{}, errors.New("invalid refresh token")
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return Response{}, errors.New("refresh token expired")
	}

	if tokenRecord.IsRevoked {
		// Token Reuse Detection
		_ = uc.authRepo.RevokeAllUserTokens(ctx, tokenRecord.UserID)
		return Response{}, errors.New("token reused")
	}

	// Revoke current token (rotation)
	err = uc.authRepo.RevokeRefreshToken(ctx, tokenRecord.ID)
	if err != nil {
		return Response{}, err
	}

	user, err := uc.userRepo.GetByID(ctx, tokenRecord.UserID)
	if err != nil {
		return Response{}, errors.New("invalid user")
	}

	userBranchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, tokenRecord.UserID)
	if err != nil {
		return Response{}, err
	}

	bStr := make([]string, len(userBranchIDs))
	for i, b := range userBranchIDs {
		bStr[i] = b.String()
	}

	accessToken, err := uc.generateToken(user.ID.String(), string(user.Role), bStr, time.Minute*15)
	if err != nil {
		return Response{}, err
	}

	refreshTokenID := uuid.New()
	refreshTokenStr, err := uc.generateToken(user.ID.String(), string(user.Role), bStr, time.Hour*24*7)
	if err != nil {
		return Response{}, err
	}

	newHashStr := hash.SHA256Token(refreshTokenStr)
	err = uc.authRepo.SaveRefreshToken(ctx, refreshTokenID, user.ID, newHashStr, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return Response{}, err
	}

	return Response{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
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
