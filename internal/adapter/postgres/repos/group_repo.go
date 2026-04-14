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

type GroupRepository struct {
	pool *pgxpool.Pool
}

func NewGroupRepository(pool *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{
		pool: pool,
	}
}

func (r *GroupRepository) db(ctx context.Context) sqlc.DBTX {
	if tx := postgres.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func (r *GroupRepository) Create(ctx context.Context, group *domain.Group) error {
	q := sqlc.New(r.db(ctx))
	id, err := q.CreateGroup(ctx, sqlc.CreateGroupParams{
		BranchID: pgtype.UUID{Bytes: group.BranchID, Valid: true},
		Name:     group.Name,
	})
	if err != nil {
		return err
	}
	group.ID = id.Bytes
	return nil
}

func (r *GroupRepository) GetByBranchID(ctx context.Context, branchID uuid.UUID, status *domain.EntityStatus) ([]*domain.GroupWithCount, error) {
	q := sqlc.New(r.db(ctx))
	rows, err := q.GetGroupsByBranchID(ctx, sqlc.GetGroupsByBranchIDParams{
		BranchID: pgtype.UUID{Bytes: branchID, Valid: true},
		Status:   toGroupNullEntityStatus(status),
	})
	if err != nil {
		return nil, err
	}
	var res []*domain.GroupWithCount
	for _, row := range rows {
		res = append(res, &domain.GroupWithCount{
			Group: domain.Group{
				ID:     row.ID.Bytes,
				Name:   row.Name,
				Status: domain.EntityStatus(row.Status.EntityStatus),
			},
			StudentsCount: int(row.StudentsCount),
		})
	}
	return res, nil
}

func toGroupNullEntityStatus(status *domain.EntityStatus) sqlc.NullEntityStatus {
	if status == nil {
		return sqlc.NullEntityStatus{}
	}

	return sqlc.NullEntityStatus{
		EntityStatus: sqlc.EntityStatus(*status),
		Valid:        true,
	}
}

func (r *GroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	q := sqlc.New(r.db(ctx))
	row, err := q.GetGroupByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &domain.Group{
		ID:        row.ID.Bytes,
		BranchID:  row.BranchID.Bytes,
		Name:      row.Name,
		Status:    domain.EntityStatus(row.Status.EntityStatus),
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *GroupRepository) UpdateName(ctx context.Context, id uuid.UUID, name string) (*domain.Group, error) {
	q := sqlc.New(r.db(ctx))
	row, err := q.UpdateGroupName(ctx, sqlc.UpdateGroupNameParams{
		Name: name,
		ID:   pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &domain.Group{
		ID:        row.ID.Bytes,
		BranchID:  row.BranchID.Bytes,
		Name:      row.Name,
		Status:    domain.EntityStatus(row.Status.EntityStatus),
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *GroupRepository) AddStudents(ctx context.Context, groupID uuid.UUID, studentIDs []uuid.UUID, joinedAt time.Time) error {
	if len(studentIDs) == 0 {
		return nil
	}

	q := sqlc.New(r.db(ctx))
	return q.AddStudentsToGroupBulk(ctx, sqlc.AddStudentsToGroupBulkParams{
		GroupID:  pgtype.UUID{Bytes: groupID, Valid: true},
		Column2:  toPGUUIDArray(studentIDs),
		JoinedAt: pgtype.Timestamptz{Time: joinedAt, Valid: true},
	})
}

func (r *GroupRepository) RemoveStudents(ctx context.Context, groupID uuid.UUID, studentIDs []uuid.UUID, leftAt time.Time) error {
	if len(studentIDs) == 0 {
		return nil
	}

	q := sqlc.New(r.db(ctx))
	return q.RemoveStudentsFromGroupBulk(ctx, sqlc.RemoveStudentsFromGroupBulkParams{
		Column1: toPGUUIDArray(studentIDs),
		GroupID: pgtype.UUID{Bytes: groupID, Valid: true},
		LeftAt:  pgtype.Timestamptz{Time: leftAt, Valid: true},
	})
}

func (r *GroupRepository) GetActiveStudentIDs(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	q := sqlc.New(r.db(ctx))
	rows, err := q.GetGroupActiveStudentIDs(ctx, pgtype.UUID{Bytes: groupID, Valid: true})
	if err != nil {
		return nil, err
	}
	var res []uuid.UUID
	for _, row := range rows {
		res = append(res, row.Bytes)
	}
	return res, nil
}

func (r *GroupRepository) GetStudents(ctx context.Context, groupID uuid.UUID) ([]*domain.GroupStudent, error) {
	q := sqlc.New(r.db(ctx))
	rows, err := q.GetGroupStudents(ctx, pgtype.UUID{Bytes: groupID, Valid: true})
	if err != nil {
		return nil, err
	}
	var res []*domain.GroupStudent
	for _, row := range rows {
		res = append(res, &domain.GroupStudent{
			ID:        row.ID.Bytes,
			FirstName: row.FirstName,
			LastName:  row.LastName,
			Status:    domain.EntityStatus(row.Status.EntityStatus),
			Phone:     pgTextToPtr(row.Phone),
			Email:     pgTextToPtr(row.Email),
		})
	}
	return res, nil
}

func toPGUUIDArray(ids []uuid.UUID) []pgtype.UUID {
	res := make([]pgtype.UUID, 0, len(ids))
	for _, id := range ids {
		res = append(res, pgtype.UUID{Bytes: id, Valid: true})
	}

	return res
}

func pgTextToPtr(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}

	value := v.String
	return &value
}

func (r *GroupRepository) GetBranchID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	q := sqlc.New(r.db(ctx))
	row, err := q.GetGroupBranchID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return uuid.UUID{}, err
	}
	return row.Bytes, nil
}

func (r *GroupRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.EntityStatus) error {
	q := sqlc.New(r.db(ctx))
	return q.UpdateGroupStatus(ctx, sqlc.UpdateGroupStatusParams{
		Status: sqlc.NullEntityStatus{
			EntityStatus: sqlc.EntityStatus(status),
			Valid:        true,
		},
		ID: pgtype.UUID{Bytes: id, Valid: true},
	})
}
