package get

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrGroupNotFound = errors.New("group not found")
)

type UseCase struct {
	groupRepo domain.GroupRepository
}

func NewUseCase(gr domain.GroupRepository) *UseCase {
	return &UseCase{
		groupRepo: gr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, groupID uuid.UUID) (Response, error) {
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return Response{}, ErrGroupNotFound
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, group.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if caller.Role == domain.RoleTeacher {
		ok, err := uc.groupRepo.IsTeacherGroup(ctx, caller.UserID, groupID)
		if err != nil {
			return Response{}, err
		}
		if !ok {
			return Response{}, domain.ErrBranchAccessDenied
		}
	}

	domainStudents, err := uc.groupRepo.GetStudents(ctx, groupID)
	if err != nil {
		return Response{}, err
	}

	var students []StudentResponse
	for _, s := range domainStudents {
		students = append(students, StudentResponse{
			ID:        s.ID,
			FirstName: s.FirstName,
			LastName:  s.LastName,
			Status:    s.Status,
			Phone:     s.Phone,
			Email:     s.Email,
		})
	}
	if students == nil {
		students = []StudentResponse{}
	}

	return Response{
		ID:       group.ID,
		Name:     group.Name,
		Status:   group.Status,
		BranchID: group.BranchID,
		Students: students,
	}, nil
}
