package mocks

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SubscriptionRepository struct {
	mock.Mock
}

func (m *SubscriptionRepository) CreatePlan(ctx context.Context, plan *domain.Plan, subjectIDs []uuid.UUID, grids []*domain.PricingGrid) error {
	args := m.Called(ctx, plan, subjectIDs, grids)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	plan.ID = uuid.New()
	return nil
}

func (m *SubscriptionRepository) GetPlansByBranchID(ctx context.Context, branchID uuid.UUID) ([]*domain.PlanDetails, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) != nil {
		return args.Get(0).([]*domain.PlanDetails), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SubscriptionRepository) GetPlanDetailsByID(ctx context.Context, id uuid.UUID) (*domain.PlanDetails, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.PlanDetails), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SubscriptionRepository) AssignToStudent(ctx context.Context, sub *domain.StudentSubscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *SubscriptionRepository) GetStudentSubscriptions(ctx context.Context, studentID uuid.UUID) ([]*domain.StudentSubscriptionDetails, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) != nil {
		return args.Get(0).([]*domain.StudentSubscriptionDetails), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SubscriptionRepository) ValidatePlanSubject(ctx context.Context, planID, subjectID uuid.UUID) (bool, error) {
	args := m.Called(ctx, planID, subjectID)
	return args.Bool(0), args.Error(1)
}

func (m *SubscriptionRepository) CountSubjectsInBranch(ctx context.Context, branchID uuid.UUID, subjectIDs []uuid.UUID) (int, error) {
	args := m.Called(ctx, branchID, subjectIDs)
	return args.Int(0), args.Error(1)
}

func (m *SubscriptionRepository) GetSubscriptionBranchIDs(ctx context.Context, studentID, planID, subjectID uuid.UUID) (*domain.SubscriptionBranchIDs, error) {
	args := m.Called(ctx, studentID, planID, subjectID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.SubscriptionBranchIDs), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SubscriptionRepository) UpdatePlanStatus(ctx context.Context, planID uuid.UUID, status domain.EntityStatus) error {
	args := m.Called(ctx, planID, status)
	return args.Error(0)
}

func (m *SubscriptionRepository) GetPlanByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Plan), args.Error(1)
	}
	return nil, args.Error(1)
}
