package repos

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) db(ctx context.Context) sqlc.DBTX {
	if tx := postgres.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	q := sqlc.New(r.db(ctx))
	row, err := q.CreateUser(ctx, sqlc.CreateUserParams{
		Phone:        user.Phone,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         sqlc.UserRole(user.Role),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAlreadyExists
		}
		return err
	}
	user.ID = row.ID.Bytes
	return nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	q := sqlc.New(r.db(ctx))
	row, err := q.GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:           row.ID.Bytes,
		Phone:        row.Phone,
		PasswordHash: row.PasswordHash,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
		Role:         domain.Role(row.Role),
		IsActive:     row.IsActive.Bool,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	q := sqlc.New(r.db(ctx))
	row, err := q.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:           row.ID.Bytes,
		Phone:        row.Phone,
		PasswordHash: row.PasswordHash,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
		Role:         domain.Role(row.Role),
		IsActive:     row.IsActive.Bool,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}, nil
}

func (r *UserRepository) AssignToBranches(ctx context.Context, userID uuid.UUID, branchIDs []uuid.UUID) error {
	q := sqlc.New(r.db(ctx))
	for _, bID := range branchIDs {
		err := q.AssignUserToBranch(ctx, sqlc.AssignUserToBranchParams{
			UserID:   pgtype.UUID{Bytes: userID, Valid: true},
			BranchID: pgtype.UUID{Bytes: bID, Valid: true},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) GetUserBranchIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	q := sqlc.New(r.db(ctx))
	rows, err := q.GetUserBranchIDs(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	res := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		res[i] = r.Bytes
	}
	return res, nil
}

func (r *UserRepository) IsBranchActive(ctx context.Context, branchID uuid.UUID) (bool, error) {
	q := sqlc.New(r.db(ctx))
	return q.IsBranchActive(ctx, pgtype.UUID{Bytes: branchID, Valid: true})
}

func (r *UserRepository) CountActiveBranchesByIDs(ctx context.Context, branchIDs []uuid.UUID) (int, error) {
	q := sqlc.New(r.db(ctx))

	idList := make([]pgtype.UUID, 0, len(branchIDs))
	for _, id := range branchIDs {
		idList = append(idList, pgtype.UUID{Bytes: id, Valid: true})
	}

	count, err := q.CountActiveBranchesByIDs(ctx, idList)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
