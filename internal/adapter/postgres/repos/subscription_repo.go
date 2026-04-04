package repos

import (
	"context"
	"fmt"

	"github.com/Xlussov/EduCRM-be/internal/adapter/postgres/postgres"
	sqlc "github.com/Xlussov/EduCRM-be/internal/adapter/postgres/sqlc"
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{
		pool: pool,
	}
}

func (r *SubscriptionRepository) db(ctx context.Context) sqlc.DBTX {
	if tx := postgres.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func mapFloatToNumeric(f float64) pgtype.Numeric {
	var num pgtype.Numeric
	_ = num.Scan(fmt.Sprintf("%f", f))
	return num
}

func mapNumericToFloat(num pgtype.Numeric) float64 {
	v, err := num.Float64Value()
	if err == nil && v.Valid {
		return v.Float64
	}
	return 0
}

func (r *SubscriptionRepository) CreatePlan(ctx context.Context, plan *domain.Plan, subjectIDs []uuid.UUID, grids []*domain.PricingGrid) error {
	q := sqlc.New(r.db(ctx))

	// 1. Create plan
	res, err := q.CreateSubscriptionPlan(ctx, sqlc.CreateSubscriptionPlanParams{
		BranchID: pgtype.UUID{Bytes: plan.BranchID, Valid: true},
		Name:     plan.Name,
		Type:     sqlc.PlanType(plan.Type),
		Status:   sqlc.NullEntityStatus{EntityStatus: sqlc.EntityStatus(plan.Status), Valid: true},
	})
	if err != nil {
		return err
	}

	plan.ID = res.ID.Bytes
	plan.CreatedAt = res.CreatedAt.Time

	// 2. Create subjects
	for _, subID := range subjectIDs {
		err := q.CreatePlanSubject(ctx, sqlc.CreatePlanSubjectParams{
			PlanID:    res.ID,
			SubjectID: pgtype.UUID{Bytes: subID, Valid: true},
		})
		if err != nil {
			return err
		}
	}

	// 3. Create grids
	for _, grid := range grids {
		err := q.CreatePlanPricingGrid(ctx, sqlc.CreatePlanPricingGridParams{
			PlanID:          res.ID,
			LessonsPerMonth: int32(grid.LessonsPerMonth),
			PricePerLesson:  mapFloatToNumeric(grid.PricePerLesson),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SubscriptionRepository) GetPlansByBranchID(ctx context.Context, branchID uuid.UUID) ([]*domain.PlanDetails, error) {
	q := sqlc.New(r.db(ctx))
	rows, err := q.GetPlansByBranchID(ctx, pgtype.UUID{Bytes: branchID, Valid: true})
	if err != nil {
		return nil, err
	}

	var plans []*domain.PlanDetails
	for _, row := range rows {
		plan := &domain.PlanDetails{
			Plan: domain.Plan{
				ID:        row.ID.Bytes,
				BranchID:  row.BranchID.Bytes,
				Name:      row.Name,
				Type:      domain.PlanType(row.Type),
				Status:    domain.EntityStatus(row.Status.EntityStatus),
				CreatedAt: row.CreatedAt.Time,
			},
		}

		// Fetch subjects
		subRows, err := q.GetPlanSubjects(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		var subjects []*domain.SubjectBase
		for _, sub := range subRows {
			subjects = append(subjects, &domain.SubjectBase{
				ID:   sub.ID.Bytes,
				Name: sub.Name,
			})
		}
		plan.Subjects = subjects

		// Fetch grids
		gridRows, err := q.GetPlanPricingGrids(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		var grids []*domain.PricingGrid
		for _, gr := range gridRows {
			grids = append(grids, &domain.PricingGrid{
				ID:              gr.ID.Bytes,
				PlanID:          gr.PlanID.Bytes,
				LessonsPerMonth: int(gr.LessonsPerMonth),
				PricePerLesson:  mapNumericToFloat(gr.PricePerLesson),
			})
		}
		plan.PricingGrid = grids

		plans = append(plans, plan)
	}

	return plans, nil
}

func (r *SubscriptionRepository) AssignToStudent(ctx context.Context, sub *domain.StudentSubscription) error {
	q := sqlc.New(r.db(ctx))
	res, err := q.CreateStudentSubscription(ctx, sqlc.CreateStudentSubscriptionParams{
		StudentID: pgtype.UUID{Bytes: sub.StudentID, Valid: true},
		PlanID:    pgtype.UUID{Bytes: sub.PlanID, Valid: true},
		SubjectID: pgtype.UUID{Bytes: sub.SubjectID, Valid: true},
		StartDate: pgtype.Date{Time: sub.StartDate, Valid: true},
	})
	if err != nil {
		return err
	}
	sub.ID = res.ID.Bytes
	sub.CreatedAt = res.CreatedAt.Time
	return nil
}

func (r *SubscriptionRepository) GetStudentSubscriptions(ctx context.Context, studentID uuid.UUID) ([]*domain.StudentSubscriptionDetails, error) {
	q := sqlc.New(r.db(ctx))
	rows, err := q.GetStudentSubscriptions(ctx, pgtype.UUID{Bytes: studentID, Valid: true})
	if err != nil {
		return nil, err
	}

	var res []*domain.StudentSubscriptionDetails
	for _, row := range rows {
		res = append(res, &domain.StudentSubscriptionDetails{
			ID:        row.ID.Bytes,
			StudentID: row.StudentID.Bytes,
			StartDate: row.StartDate.Time,
			CreatedAt: row.CreatedAt.Time,
			Plan: domain.SubPlanDetails{
				ID:   row.PlanID.Bytes,
				Name: row.PlanName,
			},
			Subject: domain.SubSubjectDetails{
				ID:   row.SubjectID.Bytes,
				Name: row.SubjectName,
			},
		})
	}
	return res, nil
}

func (r *SubscriptionRepository) ValidatePlanSubject(ctx context.Context, planID, subjectID uuid.UUID) (bool, error) {
	q := sqlc.New(r.db(ctx))
	res, err := q.ValidatePlanSubject(ctx, sqlc.ValidatePlanSubjectParams{
		PlanID:    pgtype.UUID{Bytes: planID, Valid: true},
		SubjectID: pgtype.UUID{Bytes: subjectID, Valid: true},
	})
	if err != nil {
		return false, err
	}
	return res, nil
}

func (r *SubscriptionRepository) UpdatePlanStatus(ctx context.Context, planID uuid.UUID, status domain.EntityStatus) error {
	q := sqlc.New(r.db(ctx))
	return q.UpdatePlanStatus(ctx, sqlc.UpdatePlanStatusParams{
		ID:     pgtype.UUID{Bytes: planID, Valid: true},
		Status: sqlc.NullEntityStatus{EntityStatus: sqlc.EntityStatus(status), Valid: true},
	})
}

func (r *SubscriptionRepository) GetPlanByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	q := sqlc.New(r.db(ctx))
	row, err := q.GetPlanByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &domain.Plan{
		ID:       row.ID.Bytes,
		BranchID: row.BranchID.Bytes,
		Name:     row.Name,
		Type:     domain.PlanType(row.Type),
		Status:   domain.EntityStatus(row.Status.EntityStatus),
	}, nil
}
