package get

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

func (uc *UseCase) Execute(ctx context.Context, _ domain.Caller, req Request) (Response, error) {
	admin, err := uc.userRepo.GetWithBranchesByID(ctx, req.ID)
	if err != nil {
		return Response{}, err
	}

	branches := make([]BranchResponse, 0, len(admin.Branches))
	for _, b := range admin.Branches {
		branches = append(branches, BranchResponse{ID: b.ID, Name: b.Name})
	}

	status := string(domain.StatusArchived)
	if admin.IsActive {
		status = string(domain.StatusActive)
	}

	return Response{
		ID:        admin.ID,
		FirstName: admin.FirstName,
		LastName:  admin.LastName,
		Phone:     admin.Phone,
		Status:    status,
		Branches:  branches,
	}, nil
}
