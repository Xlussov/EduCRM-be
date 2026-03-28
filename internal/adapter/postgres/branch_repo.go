package postgres

import (
	"context"

	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
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
