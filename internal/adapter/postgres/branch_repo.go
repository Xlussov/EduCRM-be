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

func (r *BranchRepository) GetAll(ctx context.Context) ([]*domain.Branch, error) {
	branches, err := r.q.GetAllBranches(ctx)
	if err != nil {
		return nil, err
	}
	var res []*domain.Branch
	for _, b := range branches {
		res = append(res, &domain.Branch{
			ID:        b.ID.Bytes,
			Name:      b.Name,
			Address:   b.Address,
			City:      b.City,
			Status:    domain.EntityStatus(b.Status.EntityStatus),
			CreatedAt: b.CreatedAt.Time,
			UpdatedAt: b.UpdatedAt.Time,
		})
	}
	return res, nil
}

func (r *BranchRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Branch, error) {
	branches, err := r.q.GetBranchesByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	var res []*domain.Branch
	for _, b := range branches {
		res = append(res, &domain.Branch{
			ID:        b.ID.Bytes,
			Name:      b.Name,
			Address:   b.Address,
			City:      b.City,
			Status:    domain.EntityStatus(b.Status.EntityStatus),
			CreatedAt: b.CreatedAt.Time,
			UpdatedAt: b.UpdatedAt.Time,
		})
	}
	return res, nil
}

func (r *BranchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error) {
	b, err := r.q.GetBranchByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &domain.Branch{
		ID:        b.ID.Bytes,
		Name:      b.Name,
		Address:   b.Address,
		City:      b.City,
		Status:    domain.EntityStatus(b.Status.EntityStatus),
		CreatedAt: b.CreatedAt.Time,
		UpdatedAt: b.UpdatedAt.Time,
	}, nil
}

func (r *BranchRepository) Update(ctx context.Context, branch *domain.Branch) error {
	return r.q.UpdateBranch(ctx, sqlc.UpdateBranchParams{
		Name:    branch.Name,
		Address: branch.Address,
		City:    branch.City,
		ID:      pgtype.UUID{Bytes: branch.ID, Valid: true},
	})
}
