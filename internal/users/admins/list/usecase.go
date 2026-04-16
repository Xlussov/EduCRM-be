package list

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	userRepo domain.UserRepository
}

func NewUseCase(ur domain.UserRepository) *UseCase {
	return &UseCase{userRepo: ur}
}

func (uc *UseCase) Execute(ctx context.Context, _ domain.Caller, _ Request) ([]AdminResponse, error) {
	admins, err := uc.userRepo.GetAdmins(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]AdminResponse, 0, len(admins))
	for _, admin := range admins {
		branches := make([]BranchResponse, 0, len(admin.Branches))
		for _, b := range admin.Branches {
			branches = append(branches, BranchResponse{ID: b.ID, Name: b.Name})
		}

		status := string(domain.StatusArchived)
		if admin.IsActive {
			status = string(domain.StatusActive)
		}

		res = append(res, AdminResponse{
			ID:        admin.ID,
			FirstName: admin.FirstName,
			LastName:  admin.LastName,
			Phone:     admin.Phone,
			Status:    status,
			Branches:  branches,
		})
	}

	return res, nil
}
