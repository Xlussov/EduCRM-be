package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	branchRepo domain.BranchRepository
}

func NewUseCase(br domain.BranchRepository) *UseCase {
	return &UseCase{
		branchRepo: br,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, req Request) ([]BranchResponse, error) {
	status, err := parseEntityStatus(req.Status)
	if err != nil {
		return nil, err
	}

	var branches []*domain.Branch

	if role == "SUPERADMIN" {
		branches, err = uc.branchRepo.GetAll(ctx, status)
	} else {
		// ADMIN
		branches, err = uc.branchRepo.GetByUserID(ctx, userID, status)
	}

	if err != nil {
		return nil, err
	}

	res := make([]BranchResponse, 0, len(branches))
	for _, b := range branches {
		res = append(res, BranchResponse{
			ID:      b.ID,
			Name:    b.Name,
			Address: b.Address,
			City:    b.City,
			Status:  string(b.Status),
		})
	}

	return res, nil
}

func parseEntityStatus(raw string) (*domain.EntityStatus, error) {
	if raw == "" {
		return nil, nil
	}

	status := domain.EntityStatus(strings.ToUpper(raw))
	if status != domain.StatusActive && status != domain.StatusArchived {
		return nil, fmt.Errorf("%w: status must be ACTIVE or ARCHIVED", domain.ErrInvalidInput)
	}

	return &status, nil
}
