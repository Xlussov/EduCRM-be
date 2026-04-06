package create

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	subjectRepo domain.SubjectRepository
}

func NewUseCase(repo domain.SubjectRepository) *UseCase {
	return &UseCase{subjectRepo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) (*Response, error) {
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
