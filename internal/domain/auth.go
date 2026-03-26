package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	IsRevoked bool
	CreatedAt time.Time
}

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, tokenID, userID uuid.UUID, hash string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, hash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id uuid.UUID) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
}
