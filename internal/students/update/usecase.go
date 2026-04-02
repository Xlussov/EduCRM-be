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

	updatedStudent, err := uc.studentRepo.Update(ctx, student)
	if err != nil {
		return Response{}, err
	}

	var dobStr *string
	if updatedStudent.Dob != nil {
		d := updatedStudent.Dob.Format("2006-01-02")
		dobStr = &d
	}

	return Response{
		ID:                 updatedStudent.ID.String(),
		BranchID:           updatedStudent.BranchID.String(),
		FirstName:          updatedStudent.FirstName,
		LastName:           updatedStudent.LastName,
		Dob:                dobStr,
		Phone:              updatedStudent.Phone,
		Email:              updatedStudent.Email,
		Address:            updatedStudent.Address,
		ParentName:         updatedStudent.ParentName,
		ParentPhone:        updatedStudent.ParentPhone,
		ParentEmail:        updatedStudent.ParentEmail,
		ParentRelationship: updatedStudent.ParentRelationship,
		Status:             string(updatedStudent.Status),
		CreatedAt:          updatedStudent.CreatedAt,
	}, nil
}
