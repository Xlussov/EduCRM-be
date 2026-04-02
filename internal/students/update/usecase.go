package update

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

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, role string, studentID uuid.UUID, req Request) (Response, error) {

	if role == "ADMIN" {
		branchID, err := uc.studentRepo.GetBranchID(ctx, studentID)
		if err != nil {
			return Response{}, err
		}

		branchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, userID)
		if err != nil {
			return Response{}, err
		}

		hasAccess := false
		for _, bid := range branchIDs {
			if bid == branchID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return Response{}, ErrBranchAccessDenied
		}
	}

	student := &domain.Student{
		ID:                 studentID,
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

	if err := uc.studentRepo.Update(ctx, student); err != nil {
		return Response{}, err
	}

	return Response{Message: "success"}, nil
}
