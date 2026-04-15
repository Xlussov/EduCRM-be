package update

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

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, subjectID uuid.UUID, req Request) (Response, error) {
	currentSubject, err := uc.subjectRepo.GetByID(ctx, subjectID)
	if err != nil {
		return Response{}, err
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, currentSubject.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if currentSubject.BranchID != req.BranchID {
		if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
			return Response{}, domain.ErrBranchAccessDenied
		}
	}

	if currentSubject.Status == domain.StatusArchived {
		return Response{}, domain.ErrCannotEditArchived
	}

	subject := &domain.Subject{
		ID:          subjectID,
		BranchID:    req.BranchID,
		Name:        req.Name,
		Description: req.Description,
	}

	updatedSubject, err := uc.subjectRepo.Update(ctx, subject)
	if err != nil {
		return Response{}, err
	}
	return Response{
		ID:          updatedSubject.ID.String(),
		BranchID:    updatedSubject.BranchID.String(),
		Name:        updatedSubject.Name,
		Description: updatedSubject.Description,
		Status:      string(updatedSubject.Status),
		CreatedAt:   updatedSubject.CreatedAt,
		UpdatedAt:   updatedSubject.UpdatedAt,
	}, nil
}
