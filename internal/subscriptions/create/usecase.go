package create

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrBranchAccessDenied = errors.New("branch access denied")
	ErrCrossBranchData    = errors.New("student, plan and subject must belong to the same branch")
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
		return Response{}, domain.ErrArchivedReference
	}

	branchIDs, err := uc.subRepo.GetSubscriptionBranchIDs(ctx, studentID, req.PlanID, req.SubjectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Response{}, domain.ErrArchivedReference
		}
		return Response{}, err
	}

	if branchIDs.StudentBranchID != branchIDs.PlanBranchID || branchIDs.PlanBranchID != branchIDs.SubjectBranchID {
		return Response{}, ErrCrossBranchData
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
