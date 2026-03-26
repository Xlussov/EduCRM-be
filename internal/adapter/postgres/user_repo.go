package postgres

import (
	"context"

	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		q: sqlc.New(pool),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	row, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Phone:        user.Phone,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         sqlc.UserRole(user.Role),
	})
	if err != nil {
		return err
	}
	user.ID = row.ID.Bytes
	// Assign generated values back to the user object
	// Normally we'd do a returning * to populate dates, assuming sqlc returned it:
	return nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	row, err := r.q.GetUserByPhone(ctx, phone)
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
	row, err := r.q.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
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
	for _, bID := range branchIDs {
		err := r.q.AssignUserToBranch(ctx, sqlc.AssignUserToBranchParams{
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
	rows, err := r.q.GetUserBranchIDs(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	res := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		res[i] = r.Bytes
	}
	return res, nil
}
