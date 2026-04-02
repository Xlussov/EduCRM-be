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

type SubjectRepository struct {
	pool *pgxpool.Pool
}

func NewSubjectRepository(pool *pgxpool.Pool) *SubjectRepository {
	return &SubjectRepository{
		pool: pool,
	}
}

func (r *SubjectRepository) db(ctx context.Context) sqlc.DBTX {
	if tx := postgres.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func (r *SubjectRepository) Create(ctx context.Context, subject *domain.Subject) error {
	q := sqlc.New(r.db(ctx))
	id, err := q.CreateSubject(ctx, sqlc.CreateSubjectParams{
		Name:        subject.Name,
		Description: pgtype.Text{String: subject.Description, Valid: subject.Description != ""},
	})
	if err != nil {
		return err
	}
	subject.ID = id.Bytes
	return nil
}

func (r *SubjectRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.EntityStatus) error {
	q := sqlc.New(r.db(ctx))
	return q.UpdateSubjectStatus(ctx, sqlc.UpdateSubjectStatusParams{
		Status: sqlc.NullEntityStatus{EntityStatus: sqlc.EntityStatus(status), Valid: true},
		ID:     pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *SubjectRepository) GetAll(ctx context.Context) ([]*domain.Subject, error) {
	q := sqlc.New(r.db(ctx))
	rows, err := q.GetAllSubjects(ctx)
	if err != nil {
		return nil, err
	}

	subjects := make([]*domain.Subject, 0, len(rows))
	for _, row := range rows {
		subjects = append(subjects, &domain.Subject{
			ID:          row.ID.Bytes,
			Name:        row.Name,
			Description: row.Description.String,
			Status:      domain.EntityStatus(row.Status.EntityStatus),
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		})
	}
	return subjects, nil
}

func (r *SubjectRepository) Update(ctx context.Context, subject *domain.Subject) error {
	q := sqlc.New(r.db(ctx))
	return q.UpdateSubject(ctx, sqlc.UpdateSubjectParams{
		Name:        subject.Name,
		Description: pgtype.Text{String: subject.Description, Valid: subject.Description != ""},
		ID:          pgtype.UUID{Bytes: subject.ID, Valid: true},
	})
}
