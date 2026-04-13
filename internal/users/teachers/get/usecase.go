package get

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	userRepo domain.UserRepository
}

func NewUseCase(ur domain.UserRepository) *UseCase {
	return &UseCase{userRepo: ur}
}

func (uc *UseCase) Execute(ctx context.Context, role string, adminBranchIDs []uuid.UUID, req Request) (Response, error) {
	teacher, err := uc.userRepo.GetWithBranchesByID(ctx, req.ID)
	if err != nil {
		return Response{}, err
	}

	if teacher.Role != domain.RoleTeacher {
		return Response{}, domain.ErrNotFound
	}

	if role == "ADMIN" {
		hasAccess := false
		for _, adminBranchID := range adminBranchIDs {
			for _, teacherBranch := range teacher.Branches {
				if adminBranchID == teacherBranch.ID {
					hasAccess = true
					break
				}
			}
			if hasAccess {
				break
			}
		}
		if !hasAccess {
			return Response{}, domain.ErrBranchAccessDenied
		}
	}

	branches := make([]BranchResponse, 0, len(teacher.Branches))
	for _, b := range teacher.Branches {
		branches = append(branches, BranchResponse{ID: b.ID, Name: b.Name})
	}

	status := string(domain.StatusArchived)
	if teacher.IsActive {
		status = string(domain.StatusActive)
	}

	return Response{
		ID:        teacher.ID,
		FirstName: teacher.FirstName,
		LastName:  teacher.LastName,
		Phone:     teacher.Phone,
		Status:    status,
		Branches:  branches,
	}, nil
}
