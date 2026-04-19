package update

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	userRepo     domain.UserRepository
	scheduleRepo domain.ScheduleRepository
	txManager    domain.TxManager
}

func NewUseCase(ur domain.UserRepository, sr domain.ScheduleRepository, tm domain.TxManager) *UseCase {
	return &UseCase{userRepo: ur, scheduleRepo: sr, txManager: tm}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, teacherID uuid.UUID, req Request) (Response, error) {
	existing, err := uc.userRepo.GetWithBranchesByID(ctx, teacherID)
	if err != nil {
		return Response{}, err
	}

	if existing.Role != domain.RoleTeacher {
		return Response{}, domain.ErrNotFound
	}

	if !existing.IsActive {
		return Response{}, domain.ErrCannotEditArchived
	}

	if domain.RequiresBranchAccess(caller.Role) {
		for _, teacherBranch := range existing.Branches {
			if !domain.HasBranchAccess(caller.BranchIDs, teacherBranch.ID) {
				return Response{}, domain.ErrBranchAccessDenied
			}
		}
		for _, branchID := range req.BranchIDs {
			if !domain.HasBranchAccess(caller.BranchIDs, branchID) {
				return Response{}, domain.ErrBranchAccessDenied
			}
		}
	}

	removed := removedBranchIDs(existing.Branches, req.BranchIDs)
	for _, branchID := range removed {
		hasFutureLessons, err := uc.scheduleRepo.CheckTeacherFutureLessonsInBranch(ctx, teacherID, branchID)
		if err != nil {
			return Response{}, err
		}
		if hasFutureLessons {
			return Response{}, domain.ErrTeacherHasFutureLessons
		}

		hasActiveTemplates, err := uc.scheduleRepo.CheckTeacherActiveTemplatesInBranch(ctx, teacherID, branchID)
		if err != nil {
			return Response{}, err
		}
		if hasActiveTemplates {
			return Response{}, domain.ErrTeacherHasActiveTemplates
		}
	}

	err = uc.txManager.Transaction(ctx, func(txCtx context.Context) error {
		activeCount, err := uc.userRepo.CountActiveBranchesByIDs(txCtx, req.BranchIDs)
		if err != nil {
			return err
		}
		if activeCount != len(req.BranchIDs) {
			return domain.ErrArchivedReference
		}

		if err := uc.userRepo.UpdateUser(txCtx, &domain.User{
			ID:        teacherID,
			Phone:     req.Phone,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}); err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				return domain.ErrPhoneAlreadyExists
			}
			return err
		}

		if err := uc.userRepo.DeleteUserBranches(txCtx, teacherID); err != nil {
			return err
		}

		if err := uc.userRepo.AssignToBranches(txCtx, teacherID, req.BranchIDs); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Response{}, err
	}

	updated, err := uc.userRepo.GetWithBranchesByID(ctx, teacherID)
	if err != nil {
		return Response{}, err
	}

	branches := make([]BranchResponse, 0, len(updated.Branches))
	for _, b := range updated.Branches {
		branches = append(branches, BranchResponse{ID: b.ID, Name: b.Name})
	}

	status := string(domain.StatusArchived)
	if updated.IsActive {
		status = string(domain.StatusActive)
	}

	return Response{
		ID:        updated.ID,
		FirstName: updated.FirstName,
		LastName:  updated.LastName,
		Phone:     updated.Phone,
		Status:    status,
		Branches:  branches,
	}, nil
}

func removedBranchIDs(existing []domain.UserBranch, requested []uuid.UUID) []uuid.UUID {
	requestedSet := make(map[uuid.UUID]struct{}, len(requested))
	for _, id := range requested {
		requestedSet[id] = struct{}{}
	}

	removed := make([]uuid.UUID, 0)
	for _, b := range existing {
		if _, ok := requestedSet[b.ID]; !ok {
			removed = append(removed, b.ID)
		}
	}

	return removed
}
