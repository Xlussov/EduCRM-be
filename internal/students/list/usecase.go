package list

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchIDRequired   = errors.New("branch_id is required")
	ErrBranchAccessDenied = errors.New("branch access denied")
)

type UseCase struct {
	studentRepo domain.StudentRepository
	userRepo    domain.UserRepository
}

func NewUseCase(sr domain.StudentRepository, ur domain.UserRepository) *UseCase {
	return &UseCase{
		studentRepo: sr,
		userRepo:    ur,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, req Request) (*Response, error) {
	if req.BranchID == uuid.Nil {
		return nil, ErrBranchIDRequired
	}

	if role == "ADMIN" {
		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return nil, err
		}
		hasAccess := false
		for _, bid := range branchIDs {
			if bid == req.BranchID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return nil, ErrBranchAccessDenied
		}
	}

	status, err := parseStudentStatus(req.Status)
	if err != nil {
		return nil, err
	}

	students, err := uc.studentRepo.GetByBranchID(ctx, req.BranchID, status)
	if err != nil {
		return nil, err
	}

	res := &Response{
		Students: make([]StudentResponse, 0, len(students)),
	}

	for _, s := range students {
		if req.Search != "" {
			search := strings.ToLower(req.Search)
			if !strings.Contains(strings.ToLower(s.FirstName), search) &&
				!strings.Contains(strings.ToLower(s.LastName), search) {
				continue
			}
		}

		res.Students = append(res.Students, StudentResponse{
			ID:        s.ID,
			FirstName: s.FirstName,
			LastName:  s.LastName,
			Status:    s.Status,
		})
	}

	return res, nil
}

func parseStudentStatus(raw string) (*domain.EntityStatus, error) {
	if raw == "" {
		return nil, nil
	}

	status := domain.EntityStatus(strings.ToUpper(raw))
	if status != domain.StatusActive && status != domain.StatusArchived {
		return nil, fmt.Errorf("%w: status must be ACTIVE or ARCHIVED", domain.ErrInvalidInput)
	}

	return &status, nil
}
