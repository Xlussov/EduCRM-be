package repos

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BranchRepository struct {
	pool *pgxpool.Pool
}

func NewBranchRepository(pool *pgxpool.Pool) *BranchRepository {
	return &BranchRepository{
		pool: pool,
	}
}

func (r *BranchRepository) db(ctx context.Context) sqlc.DBTX {
	if tx := postgres.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func (r *BranchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	q := sqlc.New(r.db(ctx))
	id, err := q.CreateBranch(ctx, sqlc.CreateBranchParams{
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
	q := sqlc.New(r.db(ctx))
	err := q.UpdateBranchStatus(ctx, sqlc.UpdateBranchStatusParams{
		Status: sqlc.NullEntityStatus{EntityStatus: sqlc.EntityStatus(status), Valid: true},
		ID:     pgtype.UUID{Bytes: id, Valid: true},
	})
	return err
}

func (r *BranchRepository) GetAll(ctx context.Context, status *domain.EntityStatus) ([]*domain.Branch, error) {
	q := sqlc.New(r.db(ctx))
	branches, err := q.GetAllBranches(ctx, toNullEntityStatus(status))
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

func (r *BranchRepository) GetByUserID(ctx context.Context, userID uuid.UUID, status *domain.EntityStatus) ([]*domain.Branch, error) {
	q := sqlc.New(r.db(ctx))
	branches, err := q.GetBranchesByUserID(ctx, sqlc.GetBranchesByUserIDParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Status: toNullEntityStatus(status),
	})
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

func (r *BranchRepository) IsActive(ctx context.Context, id uuid.UUID) (bool, error) {
	q := sqlc.New(r.db(ctx))
	return q.IsBranchActive(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *BranchRepository) CountActiveByIDs(ctx context.Context, ids []uuid.UUID) (int, error) {
	q := sqlc.New(r.db(ctx))

	idList := make([]pgtype.UUID, 0, len(ids))
	for _, id := range ids {
		idList = append(idList, pgtype.UUID{Bytes: id, Valid: true})
	}

	count, err := q.CountActiveBranchesByIDs(ctx, idList)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *BranchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error) {
	q := sqlc.New(r.db(ctx))
	b, err := q.GetBranchByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
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

func (r *BranchRepository) Update(ctx context.Context, branch *domain.Branch) (*domain.Branch, error) {
	q := sqlc.New(r.db(ctx))
	b, err := q.UpdateBranch(ctx, sqlc.UpdateBranchParams{
		Name:    branch.Name,
		Address: branch.Address,
		City:    branch.City,
		ID:      pgtype.UUID{Bytes: branch.ID, Valid: true},
	})
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

func toNullEntityStatus(status *domain.EntityStatus) sqlc.NullEntityStatus {
	if status == nil {
		return sqlc.NullEntityStatus{}
	}

	return sqlc.NullEntityStatus{
		EntityStatus: sqlc.EntityStatus(*status),
		Valid:        true,
	}
}
