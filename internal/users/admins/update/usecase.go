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

func (uc *UseCase) Execute(ctx context.Context, adminID uuid.UUID, req Request) (Response, error) {
	existing, err := uc.userRepo.GetByID(ctx, adminID)
	if err != nil {
		return Response{}, err
	}
	if existing.Role != domain.RoleAdmin {
		return Response{}, domain.ErrNotFound
	}
	if !existing.IsActive {
		return Response{}, domain.ErrCannotEditArchived
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
			ID:        adminID,
			Phone:     req.Phone,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}); err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				return domain.ErrPhoneAlreadyExists
			}
			return err
		}

		if err := uc.userRepo.DeleteUserBranches(txCtx, adminID); err != nil {
			return err
		}

		if err := uc.userRepo.AssignToBranches(txCtx, adminID, req.BranchIDs); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Response{}, err
	}

	updated, err := uc.userRepo.GetWithBranchesByID(ctx, adminID)
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
