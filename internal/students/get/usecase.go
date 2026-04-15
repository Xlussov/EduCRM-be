package get

import (
	"context"
	"errors"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var ErrStudentNotFound = errors.New("student not found")

type UseCase struct {
	studentRepo domain.StudentRepository
}

func NewUseCase(sr domain.StudentRepository) *UseCase {
	return &UseCase{studentRepo: sr}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, studentID uuid.UUID) (*Response, error) {
	student, err := uc.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, ErrStudentNotFound
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, student.BranchID) {
		return nil, domain.ErrBranchAccessDenied
	}

	if caller.Role == domain.RoleTeacher {
		ok, err := uc.studentRepo.IsTeacherStudent(ctx, caller.UserID, studentID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, domain.ErrBranchAccessDenied
		}
	}

	res := &Response{
		ID:                 student.ID,
		BranchID:           student.BranchID,
		FirstName:          student.FirstName,
		LastName:           student.LastName,
		ParentName:         student.ParentName,
		ParentPhone:        student.ParentPhone,
		ParentEmail:        student.ParentEmail,
		ParentRelationship: student.ParentRelationship,
		Phone:              student.Phone,
		Email:              student.Email,
		Address:            student.Address,
		Status:             student.Status,
		CreatedAt:          student.CreatedAt,
	}

	if student.Dob != nil {
		formatted := student.Dob.Format("2006-01-02")
		res.Dob = &formatted
	}

	return res, nil
}
