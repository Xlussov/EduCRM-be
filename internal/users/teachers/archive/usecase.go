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

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, teacherID uuid.UUID) (Response, error) {
	teacher, err := uc.userRepo.GetWithBranchesByID(ctx, teacherID)
	if err != nil {
		return Response{}, err
	}

	if teacher.Role != domain.RoleTeacher {
		return Response{}, domain.ErrNotFound
	}

	if domain.RequiresBranchAccess(caller.Role) {
		hasAccess := false
		for _, teacherBranch := range teacher.Branches {
			if domain.HasBranchAccess(caller.BranchIDs, teacherBranch.ID) {
				hasAccess = true
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
