package get

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
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

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, groupID uuid.UUID) (Response, error) {
	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return Response{}, err
	}

	if role == "ADMIN" {
		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return Response{}, err
		}
		hasAccess := false
		for _, bid := range branchIDs {
			if bid == group.BranchID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return Response{}, ErrBranchAccessDenied
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
		})
	}
	if students == nil {
		students = []StudentResponse{}
	}

	return Response{
		ID:       group.ID,
		Name:     group.Name,
		Students: students,
	}, nil
}
