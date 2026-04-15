package update

import (
	"context"
	"time"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type UseCase struct {
	studentRepo domain.StudentRepository
}

func NewUseCase(sr domain.StudentRepository) *UseCase {
	return &UseCase{
		studentRepo: sr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, studentID uuid.UUID, req Request) (Response, error) {
	currentStudent, err := uc.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return Response{}, err
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, currentStudent.BranchID) {
		return Response{}, domain.ErrBranchAccessDenied
	}

	if currentStudent.Status == domain.StatusArchived {
		return Response{}, domain.ErrCannotEditArchived
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
