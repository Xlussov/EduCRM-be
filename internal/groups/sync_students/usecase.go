package syncstudents

import (
	"context"
	"errors"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrStudentNotFound       = errors.New("student not found")
	ErrStudentBranchMismatch = errors.New("student branch does not match group branch")
	ErrStudentIDsRequired    = errors.New("student_ids is required")
)

type UseCase struct {
	groupRepo   domain.GroupRepository
	studentRepo domain.StudentRepository
	txManager   domain.TxManager
}

func NewUseCase(gr domain.GroupRepository, sr domain.StudentRepository, tx domain.TxManager) *UseCase {
	return &UseCase{
		groupRepo:   gr,
		studentRepo: sr,
		txManager:   tx,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, groupID uuid.UUID, req Request) (Response, error) {
	if req.StudentIDs == nil {
		return Response{}, ErrStudentIDsRequired
	}

	group, err := uc.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return Response{}, err
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, group.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	desiredIDs := uniqueUUIDs(req.StudentIDs)
	for _, studentID := range desiredIDs {
		studentBranchID, err := uc.studentRepo.GetBranchID(ctx, studentID)
		if err != nil {
			return Response{}, errors.Join(ErrStudentNotFound, err)
		}
		if studentBranchID != group.BranchID {
			return Response{}, ErrStudentBranchMismatch
		}
	}

	currentIDs, err := uc.groupRepo.GetActiveStudentIDs(ctx, groupID)
	if err != nil {
		return Response{}, err
	}

	toAdd, toRemove := diffStudents(currentIDs, desiredIDs)
	if len(toAdd) == 0 && len(toRemove) == 0 {
		return Response{Message: "success"}, nil
	}

	err = uc.txManager.Transaction(ctx, func(txCtx context.Context) error {
		now := time.Now()

		if len(toRemove) > 0 {
			if err := uc.groupRepo.RemoveStudents(txCtx, groupID, toRemove, now); err != nil {
				return err
			}
		}

		if len(toAdd) > 0 {
			if err := uc.groupRepo.AddStudents(txCtx, groupID, toAdd, now); err != nil {
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

func uniqueUUIDs(values []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(values))
	result := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}

func diffStudents(currentIDs, desiredIDs []uuid.UUID) ([]uuid.UUID, []uuid.UUID) {
	currentSet := make(map[uuid.UUID]struct{}, len(currentIDs))
	for _, id := range currentIDs {
		currentSet[id] = struct{}{}
	}

	desiredSet := make(map[uuid.UUID]struct{}, len(desiredIDs))
	for _, id := range desiredIDs {
		desiredSet[id] = struct{}{}
	}

	toAdd := make([]uuid.UUID, 0)
	for _, id := range desiredIDs {
		if _, exists := currentSet[id]; !exists {
			toAdd = append(toAdd, id)
		}
	}

	toRemove := make([]uuid.UUID, 0)
	for _, id := range currentIDs {
		if _, exists := desiredSet[id]; !exists {
			toRemove = append(toRemove, id)
		}
	}

	return toAdd, toRemove
}
