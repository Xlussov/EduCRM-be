package archive

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

func (uc *UseCase) Execute(ctx context.Context, role string, adminBranchIDs []uuid.UUID, teacherID uuid.UUID) (Response, error) {
	teacher, err := uc.userRepo.GetWithBranchesByID(ctx, teacherID)
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

	if !teacher.IsActive {
		return Response{}, domain.ErrAlreadyArchived
	}

	if err := uc.userRepo.UpdateUserStatus(ctx, teacherID, false); err != nil {
		return Response{}, err
	}

	return Response{Message: "success"}, nil
}
