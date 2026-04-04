package create

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchAccessDenied = errors.New("branch access denied")
	ErrInvalidSubject     = errors.New("subject is not allowed for this plan")
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

func (uc *UseCase) Execute(ctx context.Context, userID, studentID uuid.UUID, role string, req Request) (Response, error) {
	if role == "ADMIN" {
		branchID, err := uc.studentRepo.GetBranchID(ctx, studentID)
		if err != nil {
			return Response{}, err
		}

		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return Response{}, err
		}

		hasAccess := false
		for _, bid := range branchIDs {
			if bid == branchID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return Response{}, ErrBranchAccessDenied
		}
	}

	isValid, err := uc.subRepo.ValidatePlanSubject(ctx, req.PlanID, req.SubjectID)
	if err != nil {
		return Response{}, err
	}
	if !isValid {
		return Response{}, ErrInvalidSubject
	}

	sub := &domain.StudentSubscription{
		StudentID: studentID,
		PlanID:    req.PlanID,
		SubjectID: req.SubjectID,
		StartDate: req.StartDate,
	}

	if err := uc.subRepo.AssignToStudent(ctx, sub); err != nil {
		return Response{}, err
	}

	return Response{ID: sub.ID}, nil
}
