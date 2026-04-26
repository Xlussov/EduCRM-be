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

type ReportRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewReportRepository(pool *pgxpool.Pool) domain.ReportRepository {
	return &ReportRepositoryImpl{
		pool: pool,
	}
}

func (r *ReportRepositoryImpl) db(ctx context.Context) sqlc.DBTX {
	if tx := postgres.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func (r *ReportRepositoryImpl) GetStudentAttendanceHistory(ctx context.Context, studentID uuid.UUID, startDate, endDate *time.Time, subjectID *uuid.UUID) ([]domain.StudentAttendanceReportItem, error) {
	var pgStartDate pgtype.Date
	if startDate != nil {
		pgStartDate = pgtype.Date{Time: *startDate, Valid: true}
	}

	var pgEndDate pgtype.Date
	if endDate != nil {
		pgEndDate = pgtype.Date{Time: *endDate, Valid: true}
	}

	var pgSubjectID pgtype.UUID
	if subjectID != nil {
		pgSubjectID = pgtype.UUID{Bytes: *subjectID, Valid: true}
	}

	params := sqlc.GetStudentAttendanceHistoryParams{
		StudentID: pgtype.UUID{Bytes: studentID, Valid: true},
		StartDate: pgStartDate,
		EndDate:   pgEndDate,
		SubjectID: pgSubjectID,
	}

	q := sqlc.New(r.db(ctx))
	rows, err := q.GetStudentAttendanceHistory(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]domain.StudentAttendanceReportItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.StudentAttendanceReportItem{
			Date:        row.Date,
			Time:        row.StartTime,
			SubjectName: row.SubjectName,
			IsPresent:   row.IsPresent,
			Notes:       row.Notes,
		})
	}

	return result, nil
}

func (r *ReportRepositoryImpl) GetBranchStatistics(ctx context.Context, branchID uuid.UUID, startDate, endDate *time.Time) (*domain.BranchStatisticsReport, error) {
	var pgStartDate pgtype.Date
	if startDate != nil {
		pgStartDate = pgtype.Date{Time: *startDate, Valid: true}
	}

	var pgEndDate pgtype.Date
	if endDate != nil {
		pgEndDate = pgtype.Date{Time: *endDate, Valid: true}
	}

	pgBranchID := pgtype.UUID{Bytes: branchID, Valid: true}

	q := sqlc.New(r.db(ctx))

	activeStudents, err := q.CountActiveStudentsByBranch(ctx, pgBranchID)
	if err != nil {
		return nil, err
	}

	completedLessons, err := q.CountCompletedLessonsByBranch(ctx, sqlc.CountCompletedLessonsByBranchParams{
		BranchID:  pgBranchID,
		StartDate: pgStartDate,
		EndDate:   pgEndDate,
	})
	if err != nil {
		return nil, err
	}

	cancelledLessons, err := q.CountCancelledLessonsByBranch(ctx, sqlc.CountCancelledLessonsByBranchParams{
		BranchID:  pgBranchID,
		StartDate: pgStartDate,
		EndDate:   pgEndDate,
	})
	if err != nil {
		return nil, err
	}

	attendanceStats, err := q.GetBranchAttendanceStats(ctx, sqlc.GetBranchAttendanceStatsParams{
		BranchID:  pgBranchID,
		StartDate: pgStartDate,
		EndDate:   pgEndDate,
	})
	if err != nil {
		return nil, err
	}

	var attendancePercentage float64
	if attendanceStats.TotalAttendanceRecords > 0 {
		attendancePercentage = float64(attendanceStats.TotalPresentRecords) / float64(attendanceStats.TotalAttendanceRecords) * 100
	}

	return &domain.BranchStatisticsReport{
		ActiveStudents:       int(activeStudents),
		CompletedLessons:     int(completedLessons),
		CancelledLessons:     int(cancelledLessons),
		AttendancePercentage: attendancePercentage,
	}, nil
}

func (r *ReportRepositoryImpl) GetTeacherStatistics(ctx context.Context, teacherID uuid.UUID, startDate, endDate *time.Time) (*domain.TeacherStatisticsReport, error) {
	var pgStartDate pgtype.Date
	if startDate != nil {
		pgStartDate = pgtype.Date{Time: *startDate, Valid: true}
	}

	var pgEndDate pgtype.Date
	if endDate != nil {
		pgEndDate = pgtype.Date{Time: *endDate, Valid: true}
	}

	q := sqlc.New(r.db(ctx))

	params := sqlc.GetTeacherStatisticsParams{
		TeacherID: pgtype.UUID{Bytes: teacherID, Valid: true},
		StartDate: pgStartDate,
		EndDate:   pgEndDate,
	}

	stats, err := q.GetTeacherStatistics(ctx, params)
	if err != nil {
		return nil, err
	}

	return &domain.TeacherStatisticsReport{
		ScheduledLessons: int(stats.ScheduledLessons),
		CompletedLessons: int(stats.CompletedLessons),
		CancelledLessons: int(stats.CancelledLessons),
	}, nil
}
