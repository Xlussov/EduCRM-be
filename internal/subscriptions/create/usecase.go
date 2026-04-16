package create

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrCrossBranchData = errors.New("student, plan and subject must belong to the same branch")
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

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, studentID uuid.UUID, req Request) (Response, error) {
	if domain.RequiresBranchAccess(caller.Role) {
		branchID, err := uc.studentRepo.GetBranchID(ctx, studentID)
		if err != nil {
			return Response{}, err
		}

		if !domain.HasBranchAccess(caller.BranchIDs, branchID) {
			return Response{}, domain.ErrBranchAccessDenied
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
