package repos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type ScheduleRepository struct {
	pool *pgxpool.Pool
}

func NewScheduleRepository(pool *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{
		pool: pool,
	}
}

func (r *ScheduleRepository) db(ctx context.Context) sqlc.DBTX {
	if tx := postgres.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func (r *ScheduleRepository) CreateLesson(ctx context.Context, lesson *domain.Lesson) error {
	q := sqlc.New(r.db(ctx))

	studentID := pgtype.UUID{Valid: false}
	if lesson.StudentID != nil {
		studentID = pgtype.UUID{Bytes: *lesson.StudentID, Valid: true}
	}
	groupID := pgtype.UUID{Valid: false}
	if lesson.GroupID != nil {
		groupID = pgtype.UUID{Bytes: *lesson.GroupID, Valid: true}
	}
	templateID := pgtype.UUID{Valid: false}
	if lesson.TemplateID != nil {
		templateID = pgtype.UUID{Bytes: *lesson.TemplateID, Valid: true}
	}

	res, err := q.CreateLesson(ctx, sqlc.CreateLessonParams{
		BranchID:   pgtype.UUID{Bytes: lesson.BranchID, Valid: true},
		TemplateID: templateID,
		TeacherID:  pgtype.UUID{Bytes: lesson.TeacherID, Valid: true},
		SubjectID:  pgtype.UUID{Bytes: lesson.SubjectID, Valid: true},
		StudentID:  studentID,
		GroupID:    groupID,
		Date:       lesson.Date,
		StartTime:  lesson.StartTime,
		EndTime:    lesson.EndTime,
		Status:     sqlc.NullLessonStatus{LessonStatus: sqlc.LessonStatus(lesson.Status), Valid: true},
	})

	if err != nil {
		return err
	}

	lesson.ID = res.ID.Bytes
	lesson.CreatedAt = res.CreatedAt.Time
	return nil
}

func (r *ScheduleRepository) CreateTemplate(ctx context.Context, template *domain.Template) error {
	q := sqlc.New(r.db(ctx))

	studentID := pgtype.UUID{Valid: false}
	if template.StudentID != nil {
		studentID = pgtype.UUID{Bytes: *template.StudentID, Valid: true}
	}
	groupID := pgtype.UUID{Valid: false}
	if template.GroupID != nil {
		groupID = pgtype.UUID{Bytes: *template.GroupID, Valid: true}
	}

	res, err := q.CreateTemplate(ctx, sqlc.CreateTemplateParams{
		BranchID:  pgtype.UUID{Bytes: template.BranchID, Valid: true},
		TeacherID: pgtype.UUID{Bytes: template.TeacherID, Valid: true},
		SubjectID: pgtype.UUID{Bytes: template.SubjectID, Valid: true},
		StudentID: studentID,
		GroupID:   groupID,
		DayOfWeek: int32(template.DayOfWeek),
		StartTime: template.StartTime,
		EndTime:   template.EndTime,
		StartDate: template.StartDate,
		EndDate:   template.EndDate,
		IsActive:  pgtype.Bool{Bool: template.IsActive, Valid: true},
	})
	if err != nil {
		return err
	}

	template.ID = res.ID.Bytes
	return nil
}

func (r *ScheduleRepository) BulkCreateLessons(ctx context.Context, lessons []domain.Lesson) error {
	if len(lessons) == 0 {
		return nil
	}

	params := make([]sqlc.BulkCreateLessonsParams, len(lessons))
	for i, l := range lessons {
		studentID := pgtype.UUID{Valid: false}
		if l.StudentID != nil {
			studentID = pgtype.UUID{Bytes: *l.StudentID, Valid: true}
		}
		groupID := pgtype.UUID{Valid: false}
		if l.GroupID != nil {
			groupID = pgtype.UUID{Bytes: *l.GroupID, Valid: true}
		}
		templateID := pgtype.UUID{Valid: false}
		if l.TemplateID != nil {
			templateID = pgtype.UUID{Bytes: *l.TemplateID, Valid: true}
		}

		params[i] = sqlc.BulkCreateLessonsParams{
			BranchID:   pgtype.UUID{Bytes: l.BranchID, Valid: true},
			TemplateID: templateID,
			TeacherID:  pgtype.UUID{Bytes: l.TeacherID, Valid: true},
			SubjectID:  pgtype.UUID{Bytes: l.SubjectID, Valid: true},
			StudentID:  studentID,
			GroupID:    groupID,
			Date:       l.Date,
			StartTime:  l.StartTime,
			EndTime:    l.EndTime,
			Status:     sqlc.NullLessonStatus{LessonStatus: sqlc.LessonStatus(l.Status), Valid: true},
		}
	}

	q := sqlc.New(r.db(ctx))
	_, err := q.BulkCreateLessons(ctx, params)
	return err
}

func (r *ScheduleRepository) UpdateLessonStatus(ctx context.Context, id uuid.UUID, status domain.LessonStatus) error {
	q := sqlc.New(r.db(ctx))
	return q.UpdateLessonStatus(ctx, sqlc.UpdateLessonStatusParams{
		ID:     pgtype.UUID{Bytes: id, Valid: true},
		Status: sqlc.NullLessonStatus{LessonStatus: sqlc.LessonStatus(status), Valid: true},
	})
}

func (r *ScheduleRepository) GetLessonByID(ctx context.Context, id uuid.UUID) (*domain.Lesson, error) {
	q := sqlc.New(r.db(ctx))
	res, err := q.GetLessonByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	var studentID *uuid.UUID
	if res.StudentID.Valid {
		id := uuid.UUID(res.StudentID.Bytes)
		studentID = &id
	}

	var groupID *uuid.UUID
	if res.GroupID.Valid {
		id := uuid.UUID(res.GroupID.Bytes)
		groupID = &id
	}

	var templateID *uuid.UUID
	if res.TemplateID.Valid {
		id := uuid.UUID(res.TemplateID.Bytes)
		templateID = &id
	}

	return &domain.Lesson{
		ID:         res.ID.Bytes,
		BranchID:   res.BranchID.Bytes,
		TemplateID: templateID,
		TeacherID:  res.TeacherID.Bytes,
		SubjectID:  res.SubjectID.Bytes,
		StudentID:  studentID,
		GroupID:    groupID,
		Date:       res.Date,
		StartTime:  res.StartTime,
		EndTime:    res.EndTime,
		Status:     domain.LessonStatus(res.Status.LessonStatus),
		CreatedAt:  res.CreatedAt.Time,
	}, nil
}

func (r *ScheduleRepository) CheckTeacherConflict(ctx context.Context, teacherID uuid.UUID, date time.Time, start, end time.Time) (bool, error) {
	q := sqlc.New(r.db(ctx))
	return q.CheckTeacherConflict(ctx, sqlc.CheckTeacherConflictParams{
		TeacherID: pgtype.UUID{Bytes: teacherID, Valid: true},
		Date:      date,
		StartTime: start,
		EndTime:   end,
	})
}

func (r *ScheduleRepository) CheckStudentConflict(ctx context.Context, studentID uuid.UUID, date time.Time, start, end time.Time) (bool, error) {
	q := sqlc.New(r.db(ctx))
	return q.CheckStudentConflict(ctx, sqlc.CheckStudentConflictParams{
		StudentID: pgtype.UUID{Bytes: studentID, Valid: true},
		Date:      date,
		StartTime: start,
		EndTime:   end,
	})
}

func (r *ScheduleRepository) GetTeacherSchedule(ctx context.Context, teacherID uuid.UUID, from, to time.Time) ([]domain.Lesson, error) {
	q := sqlc.New(r.db(ctx))

	res, err := q.GetTeacherSchedule(ctx, sqlc.GetTeacherScheduleParams{
		TeacherID: pgtype.UUID{Bytes: teacherID, Valid: true},
		Date:      from,
		Date_2:    to,
	})
	if err != nil {
		return nil, err
	}

	lessons := make([]domain.Lesson, len(res))
	for i, l := range res {
		var studentID *uuid.UUID
		if l.StudentID.Valid {
			id := uuid.UUID(l.StudentID.Bytes)
			studentID = &id
		}
		var groupID *uuid.UUID
		if l.GroupID.Valid {
			id := uuid.UUID(l.GroupID.Bytes)
			groupID = &id
		}
		var templateID *uuid.UUID
		if l.TemplateID.Valid {
			id := uuid.UUID(l.TemplateID.Bytes)
			templateID = &id
		}
		lessons[i] = domain.Lesson{
			ID:         l.ID.Bytes,
			BranchID:   l.BranchID.Bytes,
			TemplateID: templateID,
			TeacherID:  l.TeacherID.Bytes,
			SubjectID:  l.SubjectID.Bytes,
			StudentID:  studentID,
			GroupID:    groupID,
			Date:       l.Date,
			StartTime:  l.StartTime,
			EndTime:    l.EndTime,
			Status:     domain.LessonStatus(l.Status.LessonStatus),
			CreatedAt:  l.CreatedAt.Time,
		}
	}
	return lessons, nil
}
