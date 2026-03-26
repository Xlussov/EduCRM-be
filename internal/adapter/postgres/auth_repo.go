package postgres

import (
	"context"
	"time"

	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	q *sqlc.Queries
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		q: sqlc.New(pool),
	}
}

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, tokenID, userID uuid.UUID, hash string, expiresAt time.Time) error {
	return r.q.SaveRefreshToken(ctx, sqlc.SaveRefreshTokenParams{
		ID:        pgtype.UUID{Bytes: tokenID, Valid: true},
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		TokenHash: hash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	row, err := r.q.GetRefreshToken(ctx, hash)
	if err != nil {
		return nil, err
	}
	return &domain.RefreshToken{
		ID:        row.ID.Bytes,
		UserID:    row.UserID.Bytes,
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt.Time,
		IsRevoked: row.IsRevoked.Bool,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, id uuid.UUID) error {
	return r.q.RevokeRefreshToken(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *AuthRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return r.q.RevokeAllUserTokens(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}
