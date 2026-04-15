package get

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
	s, err := uc.subjectRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &Response{
		Subject: SubjectResponse{
			ID:          s.ID.String(),
			BranchID:    s.BranchID.String(),
			Name:        s.Name,
			Description: s.Description,
			Status:      s.Status,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		},
	}, nil
}
