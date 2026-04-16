package list

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	subRepo     domain.SubscriptionRepository
	studentRepo domain.StudentRepository
}

func NewUseCase(sr domain.SubscriptionRepository, std domain.StudentRepository) *UseCase {
	return &UseCase{
		subRepo:     sr,
		studentRepo: std,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, studentID uuid.UUID) ([]SubscriptionResponse, error) {
	if domain.RequiresBranchAccess(caller.Role) {
		branchID, err := uc.studentRepo.GetBranchID(ctx, studentID)
		if err != nil {
			return nil, err
		}

		if !domain.HasBranchAccess(caller.BranchIDs, branchID) {
			return nil, domain.ErrBranchAccessDenied
		}
	}

	subs, err := uc.subRepo.GetStudentSubscriptions(ctx, studentID)
	if err != nil {
		return nil, err
	}

	res := make([]SubscriptionResponse, 0, len(subs))
	for _, s := range subs {
		res = append(res, SubscriptionResponse{
			ID:        s.ID,
			Plan:      PlanRef{ID: s.Plan.ID, Name: s.Plan.Name},
			Subject:   SubjectRef{ID: s.Subject.ID, Name: s.Subject.Name},
			StartDate: s.StartDate,
			CreatedAt: s.CreatedAt,
		})
	}

	return res, nil
}
