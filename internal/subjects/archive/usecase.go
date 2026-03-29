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

func (uc *UseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return uc.subjectRepo.UpdateStatus(ctx, id, domain.StatusArchived)
}
