package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PlanType string

const (
	PlanTypeIndividual PlanType = "INDIVIDUAL"
	PlanTypeGroup      PlanType = "GROUP"
)

type Plan struct {
	ID        uuid.UUID
	BranchID  uuid.UUID
	Name      string
	Type      PlanType
	Status    EntityStatus
	CreatedAt time.Time
}

type PricingGrid struct {
	ID              uuid.UUID
	PlanID          uuid.UUID
	LessonsPerMonth int
	PricePerLesson  float64
}

type SubjectBase struct {
	ID   uuid.UUID
	Name string
}

type PlanDetails struct {
	Plan
	Subjects    []*SubjectBase
	PricingGrid []*PricingGrid
}

type StudentSubscription struct {
	ID        uuid.UUID
	StudentID uuid.UUID
	PlanID    uuid.UUID
	SubjectID uuid.UUID
	StartDate time.Time
	CreatedAt time.Time
}

type SubPlanDetails struct {
	ID   uuid.UUID
	Name string
}

type SubSubjectDetails struct {
	ID   uuid.UUID
	Name string
}

type StudentSubscriptionDetails struct {
	ID        uuid.UUID
	StudentID uuid.UUID
	Plan      SubPlanDetails
	Subject   SubSubjectDetails
	StartDate time.Time
	CreatedAt time.Time
}

type SubscriptionRepository interface {
	CreatePlan(ctx context.Context, plan *Plan, subjectIDs []uuid.UUID, grids []*PricingGrid) error
	GetPlansByBranchID(ctx context.Context, branchID uuid.UUID) ([]*PlanDetails, error)
	UpdatePlanStatus(ctx context.Context, planID uuid.UUID, status EntityStatus) error
	AssignToStudent(ctx context.Context, sub *StudentSubscription) error
	GetStudentSubscriptions(ctx context.Context, studentID uuid.UUID) ([]*StudentSubscriptionDetails, error)
	ValidatePlanSubject(ctx context.Context, planID, subjectID uuid.UUID) (bool, error)
	GetPlanByID(ctx context.Context, id uuid.UUID) (*Plan, error)
}
