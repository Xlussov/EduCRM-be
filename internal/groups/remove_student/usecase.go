package removestudent

import (
	"context"
	"errors"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrBranchAccessDenied = errors.New("branch access denied")
	ErrStudentNotFound    = errors.New("student not found")
)

type UseCase struct {
	groupRepo   domain.GroupRepository
	userRepo    domain.UserRepository
	studentRepo domain.StudentRepository
}

func NewUseCase(gr domain.GroupRepository, ur domain.UserRepository, sr domain.StudentRepository) *UseCase {
	return &UseCase{
		groupRepo:   gr,
		userRepo:    ur,
		studentRepo: sr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, groupID uuid.UUID, studentID uuid.UUID) (Response, error) {
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

	_, err = uc.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return Response{}, errors.Join(ErrStudentNotFound, err)
	}

	err = uc.groupRepo.RemoveStudent(ctx, groupID, studentID, time.Now())
	if err != nil {
		return Response{}, err
	}

	return Response{Message: "success"}, nil
}
