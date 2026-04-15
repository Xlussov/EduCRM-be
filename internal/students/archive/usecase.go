package archive

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	studentRepo domain.StudentRepository
}

func NewUseCase(sr domain.StudentRepository) *UseCase {
	return &UseCase{
		studentRepo: sr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, studentID uuid.UUID) (Response, error) {
	student, err := uc.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return Response{}, err
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, student.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if student.Status == domain.StatusArchived {
		return Response{}, domain.ErrAlreadyArchived
	}

	err = uc.studentRepo.UpdateStatus(ctx, studentID, domain.StatusArchived)
	if err != nil {
		return Response{}, err
	}
	return Response{Message: "success"}, nil
}
