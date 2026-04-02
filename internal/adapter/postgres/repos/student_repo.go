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

type StudentRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{
		pool: pool,
	}
}

func (r *StudentRepository) db(ctx context.Context) sqlc.DBTX {
	if tx := postgres.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func (r *StudentRepository) Create(ctx context.Context, student *domain.Student) error {
	params := sqlc.CreateStudentParams{
		BranchID:    pgtype.UUID{Bytes: student.BranchID, Valid: true},
		FirstName:   student.FirstName,
		LastName:    student.LastName,
		ParentName:  student.ParentName,
		ParentPhone: student.ParentPhone,
	}

	if student.Dob != nil {
		params.Dob = pgtype.Date{Time: *student.Dob, Valid: true}
	}
	if student.Phone != nil {
		params.Phone = pgtype.Text{String: *student.Phone, Valid: true}
	}
	if student.Email != nil {
		params.Email = pgtype.Text{String: *student.Email, Valid: true}
	}
	if student.Address != nil {
		params.Address = pgtype.Text{String: *student.Address, Valid: true}
	}
	if student.ParentEmail != nil {
		params.ParentEmail = pgtype.Text{String: *student.ParentEmail, Valid: true}
	}
	if student.ParentRelationship != nil {
		params.ParentRelationship = pgtype.Text{String: *student.ParentRelationship, Valid: true}
	}

	q := sqlc.New(r.db(ctx))
	id, err := q.CreateStudent(ctx, params)
	if err != nil {
		return err
	}
	student.ID = id.Bytes
	return nil
}

func (r *StudentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.EntityStatus) error {
	q := sqlc.New(r.db(ctx))
	err := q.UpdateStudentStatus(ctx, sqlc.UpdateStudentStatusParams{
		Status: sqlc.NullEntityStatus{EntityStatus: sqlc.EntityStatus(status), Valid: true},
		ID:     pgtype.UUID{Bytes: id, Valid: true},
	})
	return err
}

func (r *StudentRepository) GetBranchID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	q := sqlc.New(r.db(ctx))
	branchID, err := q.GetStudentBranchID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return uuid.UUID{}, err
	}
	return branchID.Bytes, nil
}

func (r *StudentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Student, error) {
	q := sqlc.New(r.db(ctx))
	s, err := q.GetStudentByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return r.toDomain(s), nil
}

func (r *StudentRepository) Update(ctx context.Context, student *domain.Student) (*domain.Student, error) {
	params := sqlc.UpdateStudentParams{
		FirstName:   student.FirstName,
		LastName:    student.LastName,
		ParentName:  student.ParentName,
		ParentPhone: student.ParentPhone,
		ID:          pgtype.UUID{Bytes: student.ID, Valid: true},
	}
	if student.Dob != nil {
		params.Dob = pgtype.Date{Time: *student.Dob, Valid: true}
	}
	if student.Phone != nil {
		params.Phone = pgtype.Text{String: *student.Phone, Valid: true}
	}
	if student.Email != nil {
		params.Email = pgtype.Text{String: *student.Email, Valid: true}
	}
	if student.Address != nil {
		params.Address = pgtype.Text{String: *student.Address, Valid: true}
	}
	if student.ParentEmail != nil {
		params.ParentEmail = pgtype.Text{String: *student.ParentEmail, Valid: true}
	}
	if student.ParentRelationship != nil {
		params.ParentRelationship = pgtype.Text{String: *student.ParentRelationship, Valid: true}
	}
	q := sqlc.New(r.db(ctx))
	s, err := q.UpdateStudent(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.toDomain(s), nil
}

func (r *StudentRepository) GetByBranchID(ctx context.Context, branchID uuid.UUID) ([]*domain.Student, error) {
	q := sqlc.New(r.db(ctx))
	students, err := q.GetStudentsByBranchID(ctx, pgtype.UUID{Bytes: branchID, Valid: true})
	if err != nil {
		return nil, err
	}
	var res []*domain.Student
	for _, s := range students {
		res = append(res, r.toDomain(s))
	}
	return res, nil
}

func (r *StudentRepository) toDomain(s sqlc.Student) *domain.Student {
	student := &domain.Student{
		ID:          s.ID.Bytes,
		BranchID:    s.BranchID.Bytes,
		FirstName:   s.FirstName,
		LastName:    s.LastName,
		ParentName:  s.ParentName,
		ParentPhone: s.ParentPhone,
		Status:      domain.EntityStatus(s.Status.EntityStatus),
		CreatedAt:   s.CreatedAt.Time,
	}
	if s.Dob.Valid {
		t := s.Dob.Time
		student.Dob = &t
	}
	if s.Phone.Valid {
		student.Phone = &s.Phone.String
	}
	if s.Email.Valid {
		student.Email = &s.Email.String
	}
	if s.Address.Valid {
		student.Address = &s.Address.String
	}
	if s.ParentEmail.Valid {
		student.ParentEmail = &s.ParentEmail.String
	}
	if s.ParentRelationship.Valid {
		student.ParentRelationship = &s.ParentRelationship.String
	}
	return student
}
