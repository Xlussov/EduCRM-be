package mocks

import (
	"context"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type AuthRepository struct {
	mock.Mock
}

func (m *AuthRepository) SaveRefreshToken(ctx context.Context, tokenID, userID uuid.UUID, hash string, expiresAt time.Time) error {
	args := m.Called(ctx, tokenID, userID, hash, expiresAt)
	return args.Error(0)
}

func (m *AuthRepository) GetRefreshToken(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, hash)
	var r0 *domain.RefreshToken
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.RefreshToken)
	}
	return r0, args.Error(1)
}

func (m *AuthRepository) RevokeRefreshToken(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *AuthRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
