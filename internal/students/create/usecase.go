package create

import (
	"context"
	"errors"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchAccessDenied = errors.New("branch access denied")
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

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, req Request) (Response, error) {

	if role == "ADMIN" {
		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return Response{}, err
		}
		hasAccess := false
		for _, bid := range branchIDs {
			if bid == req.BranchID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return Response{}, ErrBranchAccessDenied
		}
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
