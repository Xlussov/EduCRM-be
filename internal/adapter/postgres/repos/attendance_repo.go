package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type AttendanceRepository struct {
	pool *pgxpool.Pool
}

func NewAttendanceRepository(pool *pgxpool.Pool) *AttendanceRepository {
	return &AttendanceRepository{pool: pool}
}

func (r *AttendanceRepository) db(ctx context.Context) sqlc.DBTX {
	if tx := postgres.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func (r *AttendanceRepository) UpsertAttendance(ctx context.Context, attendance []domain.Attendance) error {
	if len(attendance) == 0 {
		return nil
	}

	lessonID := attendance[0].LessonID
	studentIDs := make([]pgtype.UUID, 0, len(attendance))
	isPresent := make([]bool, 0, len(attendance))
	notes := make([]string, 0, len(attendance))

	for _, record := range attendance {
		studentIDs = append(studentIDs, pgtype.UUID{Bytes: record.StudentID, Valid: true})
		isPresent = append(isPresent, record.IsPresent)
		if record.Notes != nil {
			notes = append(notes, *record.Notes)
		} else {
			notes = append(notes, "")
		}
	}

	q := sqlc.New(r.db(ctx))
	return q.UpsertAttendance(ctx, sqlc.UpsertAttendanceParams{
		Column1: pgtype.UUID{Bytes: lessonID, Valid: true},
		Column2: studentIDs,
		Column3: isPresent,
		Column4: notes,
	})
}

func (r *AttendanceRepository) GetLessonAttendance(ctx context.Context, lessonID uuid.UUID) ([]domain.LessonAttendanceStudent, error) {
	q := sqlc.New(r.db(ctx))
	rows, err := q.GetLessonAttendance(ctx, pgtype.UUID{Bytes: lessonID, Valid: true})
	if err != nil {
		return nil, err
	}

	res := make([]domain.LessonAttendanceStudent, 0, len(rows))
	for _, row := range rows {
		var isPresent *bool
		if row.IsPresent.Valid {
			value := row.IsPresent.Bool
			isPresent = &value
		}

		var notes *string
		if row.Notes.Valid {
			value := row.Notes.String
			notes = &value
		}

		res = append(res, domain.LessonAttendanceStudent{
			StudentID: row.StudentID.Bytes,
			FirstName: row.FirstName,
			LastName:  row.LastName,
			Status:    domain.EntityStatus(row.Status.EntityStatus),
			IsPresent: isPresent,
			Notes:     notes,
		})
	}

	return res, nil
}

func (r *AttendanceRepository) GetStudentAttendance(ctx context.Context, studentID uuid.UUID, from, to time.Time) ([]domain.Attendance, error) {
	const query = `
		SELECT a.id, a.lesson_id, a.student_id, a.is_present, a.notes, a.created_at
		FROM attendance a
		JOIN lessons l ON l.id = a.lesson_id
		WHERE a.student_id = $1 AND l.date >= $2 AND l.date <= $3
		ORDER BY l.date ASC, l.start_time ASC
	`

	rows, err := r.db(ctx).Query(ctx, query, pgtype.UUID{Bytes: studentID, Valid: true}, from, to)
	if err != nil {
		return nil, fmt.Errorf("query student attendance: %w", err)
	}
	defer rows.Close()

	res := make([]domain.Attendance, 0)
	for rows.Next() {
		var id pgtype.UUID
		var lessonID pgtype.UUID
		var sID pgtype.UUID
		var isPresent bool
		var notes pgtype.Text
		var createdAt pgtype.Timestamptz

		if err := rows.Scan(&id, &lessonID, &sID, &isPresent, &notes, &createdAt); err != nil {
			return nil, fmt.Errorf("scan student attendance: %w", err)
		}

		var notePtr *string
		if notes.Valid {
			value := notes.String
			notePtr = &value
		}

		res = append(res, domain.Attendance{
			ID:        id.Bytes,
			LessonID:  lessonID.Bytes,
			StudentID: sID.Bytes,
			IsPresent: isPresent,
			Notes:     notePtr,
			CreatedAt: createdAt.Time,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan student attendance: %w", err)
	}

	return res, nil
}
