package create

import (
	"context"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	studentRepo domain.StudentRepository
	userRepo    domain.UserRepository
}

func NewUseCase(sr domain.StudentRepository, ur domain.UserRepository) *UseCase {
	return &UseCase{
		studentRepo: sr,
		userRepo:    ur,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (Response, error) {
	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	isActive, err := uc.userRepo.IsBranchActive(ctx, req.BranchID)
	if err != nil {
		return Response{}, err
	}
	if !isActive {
		return Response{}, domain.ErrArchivedReference
	}

	student := &domain.Student{
		BranchID:           req.BranchID,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		ParentName:         req.ParentName,
		ParentPhone:        req.ParentPhone,
		Phone:              req.Phone,
		Email:              req.Email,
		Address:            req.Address,
		ParentEmail:        req.ParentEmail,
		ParentRelationship: req.ParentRelationship,
	}

	if req.Dob != nil && *req.Dob != "" {
		parsed, _ := time.Parse("2006-01-02", *req.Dob)
		student.Dob = &parsed
	}

	if err := uc.studentRepo.Create(ctx, student); err != nil {
		return Response{}, err
	}

	return Response{ID: student.ID}, nil
}
