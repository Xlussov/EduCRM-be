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

	updatedSubject, err := uc.subjectRepo.Update(ctx, subject)
	if err != nil {
		return Response{}, err
	}
	return Response{
		ID:          updatedSubject.ID.String(),
		Name:        updatedSubject.Name,
		Description: updatedSubject.Description,
		Status:      string(updatedSubject.Status),
		CreatedAt:   updatedSubject.CreatedAt,
		UpdatedAt:   updatedSubject.UpdatedAt,
	}, nil
}
