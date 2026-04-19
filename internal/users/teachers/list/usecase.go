package list

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

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) ([]TeacherResponse, error) {
	var filter []uuid.UUID

	if caller.Role == domain.RoleSuperadmin {
		if req.BranchID != nil {
			filter = []uuid.UUID{*req.BranchID}
		}
	} else if domain.RequiresBranchAccess(caller.Role) {
		if len(caller.BranchIDs) == 0 {
			return []TeacherResponse{}, nil
		}

		if req.BranchID != nil {
			if !domain.HasBranchAccess(caller.BranchIDs, *req.BranchID) {
				return nil, domain.ErrBranchAccessDenied
			}
			filter = []uuid.UUID{*req.BranchID}
		} else {
			filter = caller.BranchIDs
		}
	}

	teachers, err := uc.userRepo.GetTeachers(ctx, filter)
	if err != nil {
		return nil, err
	}

	res := make([]TeacherResponse, 0, len(teachers))
	for _, teacher := range teachers {
		branches := make([]BranchResponse, 0, len(teacher.Branches))
		for _, b := range teacher.Branches {
			branches = append(branches, BranchResponse{ID: b.ID, Name: b.Name})
		}

		status := string(domain.StatusArchived)
		if teacher.IsActive {
			status = string(domain.StatusActive)
		}

		res = append(res, TeacherResponse{
			ID:        teacher.ID,
			FirstName: teacher.FirstName,
			LastName:  teacher.LastName,
			Phone:     teacher.Phone,
			Status:    status,
			Branches:  branches,
		})
	}

	return res, nil
}
