package list

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchAccessDenied = errors.New("branch access denied")
)

type UseCase struct {
	subRepo     domain.SubscriptionRepository
	userRepo    domain.UserRepository
	studentRepo domain.StudentRepository
}

func NewUseCase(sr domain.SubscriptionRepository, ur domain.UserRepository, std domain.StudentRepository) *UseCase {
	return &UseCase{
		subRepo:     sr,
		userRepo:    ur,
		studentRepo: std,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID, studentID uuid.UUID, role string) ([]SubscriptionResponse, error) {
	if role == "ADMIN" {
		branchID, err := uc.studentRepo.GetBranchID(ctx, studentID)
		if err != nil {
			return nil, err
		}

		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return nil, err
		}

		hasAccess := false
		for _, bid := range branchIDs {
			if bid == branchID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return nil, ErrBranchAccessDenied
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
