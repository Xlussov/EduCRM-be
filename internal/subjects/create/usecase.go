package create

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	subjectRepo domain.SubjectRepository
	branchRepo  domain.BranchRepository
}

func NewUseCase(sr domain.SubjectRepository, br domain.BranchRepository) *UseCase {
	return &UseCase{
		subjectRepo: sr,
		branchRepo:  br,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (*Response, error) {
	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
		return nil, domain.ErrBranchAccessDenied
	}

	isActive, err := uc.branchRepo.IsActive(ctx, req.BranchID)
	if err != nil {
		return nil, err
	}
	if !isActive {
		return nil, domain.ErrArchivedReference
	}

	subject := &domain.Subject{
		BranchID:    req.BranchID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := uc.subjectRepo.Create(ctx, subject); err != nil {
		return nil, err
	}

	return &Response{
		ID:          subject.ID.String(),
		BranchID:    subject.BranchID.String(),
		Name:        subject.Name,
		Description: subject.Description,
	}, nil
}
