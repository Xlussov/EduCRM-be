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

func (uc *UseCase) Execute(ctx context.Context, studentID uuid.UUID) (Response, error) {
	err := uc.studentRepo.UpdateStatus(ctx, studentID, domain.StatusArchived)
	if err != nil {
		return Response{}, err
	}
	return Response{Message: "success"}, nil
}
