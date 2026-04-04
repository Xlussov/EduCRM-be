package addstudents

import (
	"context"
	"errors"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrGroupNotFound         = errors.New("group not found")
	ErrBranchAccessDenied    = errors.New("branch access denied")
	ErrStudentNotFound       = errors.New("student not found")
	ErrStudentBranchMismatch = errors.New("student branch does not match group branch")
)

type UseCase struct {
	groupRepo   domain.GroupRepository
	userRepo    domain.UserRepository
	studentRepo domain.StudentRepository
	txManager   domain.TxManager
}

func NewUseCase(gr domain.GroupRepository, ur domain.UserRepository, sr domain.StudentRepository, tx domain.TxManager) *UseCase {
	return &UseCase{
		groupRepo:   gr,
		userRepo:    ur,
		studentRepo: sr,
		txManager:   tx,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, groupID uuid.UUID, req Request) (Response, error) {
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

	for _, sID := range req.StudentIDs {
		sBranchID, err := uc.studentRepo.GetBranchID(ctx, sID)
		if err != nil {
			return Response{}, errors.Join(ErrStudentNotFound, err)
		}
		if sBranchID != group.BranchID {
			return Response{}, ErrStudentBranchMismatch
		}
	}

	err = uc.txManager.Transaction(ctx, func(txCtx context.Context) error {
		now := time.Now()
		for _, sID := range req.StudentIDs {
			if err := uc.groupRepo.AddStudent(txCtx, groupID, sID, now); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return Response{}, err
	}

	return Response{Message: "success"}, nil
}
