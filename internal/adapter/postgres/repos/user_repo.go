package repos

import (
	"context"
	"errors"
	"fmt"

	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
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

func (r *UserRepository) GetWithBranchesByID(ctx context.Context, id uuid.UUID) (*domain.UserWithBranches, error) {
	q := sqlc.New(r.db(ctx))
	rows, err := q.GetUserWithBranchesByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("get user with branches by id: %w", err)
	}
	if len(rows) == 0 {
		return nil, domain.ErrNotFound
	}

	first := rows[0]
	res := &domain.UserWithBranches{
		ID:        first.ID.Bytes,
		Phone:     first.Phone,
		FirstName: first.FirstName,
		LastName:  first.LastName,
		Role:      domain.Role(first.Role),
		IsActive:  first.IsActive.Bool,
		CreatedAt: first.CreatedAt.Time,
		UpdatedAt: first.UpdatedAt.Time,
		Branches:  make([]domain.UserBranch, 0, len(rows)),
	}

	for _, row := range rows {
		if row.BranchID.Valid {
			res.Branches = append(res.Branches, domain.UserBranch{
				ID:   row.BranchID.Bytes,
				Name: row.BranchName.String,
			})
		}
	}

	return res, nil
}

func (r *UserRepository) GetAdmins(ctx context.Context) ([]*domain.UserWithBranches, error) {
	q := sqlc.New(r.db(ctx))
	rows, err := q.GetAdmins(ctx)
	if err != nil {
		return nil, fmt.Errorf("get admins: %w", err)
	}

	return foldUsersWithBranchesFromAdmins(rows), nil
}

func (r *UserRepository) GetTeachers(ctx context.Context, branchIDs []uuid.UUID) ([]*domain.UserWithBranches, error) {
	q := sqlc.New(r.db(ctx))

	var filter []pgtype.UUID
	if len(branchIDs) > 0 {
		filter = make([]pgtype.UUID, 0, len(branchIDs))
		for _, id := range branchIDs {
			filter = append(filter, pgtype.UUID{Bytes: id, Valid: true})
		}
	}

	rows, err := q.GetTeachers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get teachers: %w", err)
	}

	return foldUsersWithBranchesFromTeachers(rows), nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	q := sqlc.New(r.db(ctx))
	rowsAffected, err := q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:        pgtype.UUID{Bytes: user.ID, Valid: true},
		Phone:     user.Phone,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("update user: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUserStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	q := sqlc.New(r.db(ctx))
	rowsAffected, err := q.UpdateUserStatus(ctx, sqlc.UpdateUserStatusParams{
		ID:       pgtype.UUID{Bytes: id, Valid: true},
		IsActive: pgtype.Bool{Bool: isActive, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("update user status: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *UserRepository) DeleteUserBranches(ctx context.Context, userID uuid.UUID) error {
	q := sqlc.New(r.db(ctx))
	if err := q.DeleteUserBranches(ctx, pgtype.UUID{Bytes: userID, Valid: true}); err != nil {
		return fmt.Errorf("delete user branches: %w", err)
	}

	return nil
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

func (r *UserRepository) CheckTeacherInBranch(ctx context.Context, teacherID, branchID uuid.UUID) (bool, error) {
	q := sqlc.New(r.db(ctx))
	return q.CheckTeacherInBranch(ctx, sqlc.CheckTeacherInBranchParams{
		UserID:   pgtype.UUID{Bytes: teacherID, Valid: true},
		BranchID: pgtype.UUID{Bytes: branchID, Valid: true},
	})
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

func foldUsersWithBranchesFromAdmins(rows []sqlc.GetAdminsRow) []*domain.UserWithBranches {
	users := make(map[uuid.UUID]*domain.UserWithBranches, len(rows))
	order := make([]uuid.UUID, 0, len(rows))

	for _, row := range rows {
		id := row.ID.Bytes
		u, ok := users[id]
		if !ok {
			u = &domain.UserWithBranches{
				ID:        id,
				Phone:     row.Phone,
				FirstName: row.FirstName,
				LastName:  row.LastName,
				Role:      domain.Role(row.Role),
				IsActive:  row.IsActive.Bool,
				CreatedAt: row.CreatedAt.Time,
				UpdatedAt: row.UpdatedAt.Time,
				Branches:  make([]domain.UserBranch, 0, 1),
			}
			users[id] = u
			order = append(order, id)
		}

		if row.BranchID.Valid {
			u.Branches = append(u.Branches, domain.UserBranch{ID: row.BranchID.Bytes, Name: row.BranchName.String})
		}
	}

	res := make([]*domain.UserWithBranches, 0, len(users))
	for _, id := range order {
		res = append(res, users[id])
	}

	return res
}

func foldUsersWithBranchesFromTeachers(rows []sqlc.GetTeachersRow) []*domain.UserWithBranches {
	users := make(map[uuid.UUID]*domain.UserWithBranches, len(rows))
	order := make([]uuid.UUID, 0, len(rows))

	for _, row := range rows {
		id := row.ID.Bytes
		u, ok := users[id]
		if !ok {
			u = &domain.UserWithBranches{
				ID:        id,
				Phone:     row.Phone,
				FirstName: row.FirstName,
				LastName:  row.LastName,
				Role:      domain.Role(row.Role),
				IsActive:  row.IsActive.Bool,
				CreatedAt: row.CreatedAt.Time,
				UpdatedAt: row.UpdatedAt.Time,
				Branches:  make([]domain.UserBranch, 0, 1),
			}
			users[id] = u
			order = append(order, id)
		}

		if row.BranchID.Valid {
			u.Branches = append(u.Branches, domain.UserBranch{ID: row.BranchID.Bytes, Name: row.BranchName.String})
		}
	}

	res := make([]*domain.UserWithBranches, 0, len(users))
	for _, id := range order {
		res = append(res, users[id])
	}

	return res
}
