package list

import (
	"context"
	"errors"
	"strings"

	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrBranchIDRequired = errors.New("branch_id is required")
)

type UseCase struct {
	studentRepo domain.StudentRepository
}

func NewUseCase(sr domain.StudentRepository) *UseCase {
	return &UseCase{
		studentRepo: sr,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (*Response, error) {
	if req.BranchID == uuid.Nil {
		return nil, ErrBranchIDRequired
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
		return nil, domain.ErrBranchAccessDenied
	}

	status, err := domain.ParseEntityStatus(req.Status)
	if err != nil {
		return nil, err
	}

	var students []*domain.Student
	if caller.Role == domain.RoleTeacher {
		students, err = uc.studentRepo.GetByBranchIDAndTeacherID(ctx, req.BranchID, caller.UserID, status)
	} else {
		students, err = uc.studentRepo.GetByBranchID(ctx, req.BranchID, status)
	}
	if err != nil {
		return nil, err
	}

	res := &Response{
		Students: make([]StudentResponse, 0, len(students)),
	}

	for _, s := range students {
		if req.Search != "" {
			search := strings.ToLower(req.Search)
			if !strings.Contains(strings.ToLower(s.FirstName), search) &&
				!strings.Contains(strings.ToLower(s.LastName), search) {
				continue
			}
		}

		res.Students = append(res.Students, StudentResponse{
			ID:        s.ID,
			FirstName: s.FirstName,
			LastName:  s.LastName,
			Phone:     s.Phone,
			Email:     s.Email,
			Status:    s.Status,
		})
	}

	return res, nil
}
