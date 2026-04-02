package repos

import (
	"context"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		pool: pool,
	}
}

func (r *AuthRepository) db(ctx context.Context) sqlc.DBTX {
	if tx := postgres.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, tokenID, userID uuid.UUID, hash string, expiresAt time.Time) error {
	q := sqlc.New(r.db(ctx))
	return q.SaveRefreshToken(ctx, sqlc.SaveRefreshTokenParams{
		ID:        pgtype.UUID{Bytes: tokenID, Valid: true},
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		TokenHash: hash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	q := sqlc.New(r.db(ctx))
	row, err := q.GetRefreshToken(ctx, hash)
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
	q := sqlc.New(r.db(ctx))
	return q.RevokeRefreshToken(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *AuthRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	q := sqlc.New(r.db(ctx))
	return q.RevokeAllUserTokens(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}
