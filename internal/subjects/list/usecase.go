package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	subjectRepo domain.SubjectRepository
}

func NewUseCase(repo domain.SubjectRepository) *UseCase {
	return &UseCase{subjectRepo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) (*Response, error) {
	status, err := parseSubjectStatus(req.Status)
	if err != nil {
		return nil, err
	}

	subjects, err := uc.subjectRepo.GetAll(ctx, req.BranchID, status)
	if err != nil {
		return nil, err
	}

	res := &Response{
		Subjects: make([]SubjectResponse, 0, len(subjects)),
	}

	for _, s := range subjects {
		res.Subjects = append(res.Subjects, SubjectResponse{
			ID:          s.ID.String(),
			BranchID:    s.BranchID.String(),
			Name:        s.Name,
			Description: s.Description,
			Status:      s.Status,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}

	return res, nil
}

func parseSubjectStatus(raw string) (*domain.EntityStatus, error) {
	if raw == "" {
		return nil, nil
	}

	status := domain.EntityStatus(strings.ToUpper(raw))
	if status != domain.StatusActive && status != domain.StatusArchived {
		return nil, fmt.Errorf("%w: status must be ACTIVE or ARCHIVED", domain.ErrInvalidInput)
	}

	return &status, nil
}
