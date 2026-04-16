package update

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	userRepo  domain.UserRepository
	txManager domain.TxManager
}

func NewUseCase(ur domain.UserRepository, tm domain.TxManager) *UseCase {
	return &UseCase{userRepo: ur, txManager: tm}
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
		hasCurrentAccess := false
		for _, teacherBranch := range existing.Branches {
			if domain.HasBranchAccess(caller.BranchIDs, teacherBranch.ID) {
				hasCurrentAccess = true
				break
			}
		}
		if !hasCurrentAccess {
			return Response{}, domain.ErrBranchAccessDenied
		}

		if !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
			return Response{}, domain.ErrBranchAccessDenied
		}
	}

	err = uc.txManager.Transaction(ctx, func(txCtx context.Context) error {
		isActive, err := uc.userRepo.IsBranchActive(txCtx, req.BranchID)
		if err != nil {
			return err
		}
		if !isActive {
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

		if err := uc.userRepo.AssignToBranches(txCtx, teacherID, []uuid.UUID{req.BranchID}); err != nil {
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
