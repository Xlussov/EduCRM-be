package postgres

import (
	"context"

	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubjectRepository struct {
	q *sqlc.Queries
}

func NewSubjectRepository(pool *pgxpool.Pool) *SubjectRepository {
	return &SubjectRepository{
		q: sqlc.New(pool),
	}
}

func (r *SubjectRepository) Create(ctx context.Context, subject *domain.Subject) error {
	id, err := r.q.CreateSubject(ctx, sqlc.CreateSubjectParams{
		Name:        subject.Name,
		Description: pgtype.Text{String: subject.Description, Valid: subject.Description != ""},
	})
	if err != nil {
		return err
	}
	subject.ID = id.Bytes
	return nil
}
