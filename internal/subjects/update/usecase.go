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

func (uc *UseCase) Execute(ctx context.Context, subjectID uuid.UUID, req Request) (Response, error) {
	subject := &domain.Subject{
		ID:          subjectID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := uc.subjectRepo.Update(ctx, subject); err != nil {
		return Response{}, err
	}
	return Response{Message: "success"}, nil
}
