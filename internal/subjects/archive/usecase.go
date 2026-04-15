package archive

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	subjectRepo domain.SubjectRepository
}

func NewUseCase(repo domain.SubjectRepository) *UseCase {
	return &UseCase{subjectRepo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, id uuid.UUID) (Response, error) {
	subject, err := uc.subjectRepo.GetByID(ctx, id)
	if err != nil {
		return Response{}, err
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, subject.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if subject.Status == domain.StatusArchived {
		return Response{}, domain.ErrAlreadyArchived
	}

	err = uc.subjectRepo.UpdateStatus(ctx, id, domain.StatusArchived)
	if err != nil {
		return Response{}, err
	}
	return Response{Message: "success"}, nil
}
