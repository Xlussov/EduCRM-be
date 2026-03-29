package postgres

import (
	"context"

	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BranchRepository struct {
	q *sqlc.Queries
}

func NewBranchRepository(pool *pgxpool.Pool) *BranchRepository {
	return &BranchRepository{
		q: sqlc.New(pool),
	}
}

func (r *BranchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	id, err := r.q.CreateBranch(ctx, sqlc.CreateBranchParams{
		Name:    branch.Name,
		Address: branch.Address,
		City:    branch.City,
	})
	if err != nil {
		return err
	}
	branch.ID = id.Bytes
	return nil
}

func (r *BranchRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.EntityStatus) error {
	return r.q.UpdateBranchStatus(ctx, sqlc.UpdateBranchStatusParams{
		Status: sqlc.NullEntityStatus{EntityStatus: sqlc.EntityStatus(status), Valid: true},
		ID:     pgtype.UUID{Bytes: id, Valid: true},
	})
}
