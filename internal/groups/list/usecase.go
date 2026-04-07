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
	groupRepo domain.GroupRepository
	userRepo  domain.UserRepository
}

func NewUseCase(gr domain.GroupRepository, ur domain.UserRepository) *UseCase {
	return &UseCase{
		groupRepo: gr,
		userRepo:  ur,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, req Request) (Response, error) {
	if req.BranchID == uuid.Nil {
		return Response{}, ErrBranchIDRequired
	}

	if role == "ADMIN" {
		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return Response{}, err
		}

		hasAccess := false
		for _, bid := range branchIDs {
			if bid == req.BranchID {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			return Response{}, ErrBranchAccessDenied
		}
	}

	status, err := parseGroupStatus(req.Status)
	if err != nil {
		return Response{}, err
	}

	groups, err := uc.groupRepo.GetByBranchID(ctx, req.BranchID, status)
	if err != nil {
		return Response{}, err
	}

	var res []GroupResponse
	for _, g := range groups {
		res = append(res, GroupResponse{
			ID:            g.ID,
			Name:          g.Name,
			StudentsCount: g.StudentsCount,
			Status:        g.Status,
		})
	}

	if res == nil {
		res = []GroupResponse{}
	}

	return Response{Groups: res}, nil
}

func parseGroupStatus(raw string) (*domain.EntityStatus, error) {
	if raw == "" {
		return nil, nil
	}

	status := domain.EntityStatus(strings.ToUpper(raw))
	if status != domain.StatusActive && status != domain.StatusArchived {
		return nil, fmt.Errorf("%w: status must be ACTIVE or ARCHIVED", domain.ErrInvalidInput)
	}

	return &status, nil
}
