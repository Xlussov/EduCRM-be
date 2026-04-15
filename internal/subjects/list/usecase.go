package list

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	subjectRepo domain.SubjectRepository
}

func NewUseCase(repo domain.SubjectRepository) *UseCase {
	return &UseCase{subjectRepo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (*Response, error) {
	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
		return nil, domain.ErrBranchAccessDenied
	}

	status, err := domain.ParseEntityStatus(req.Status)
	if err != nil {
		return nil, err
	}

	subjects, err := uc.subjectRepo.GetAll(ctx, req.BranchID, status)
	if err != nil {
		return nil, err
	}

	res := &Response{
		Subjects: make([]SubjectResponse, 0, len(subjects)),
	}

	for _, s := range subjects {
		res.Subjects = append(res.Subjects, SubjectResponse{
			ID:          s.ID.String(),
			BranchID:    s.BranchID.String(),
			Name:        s.Name,
			Description: s.Description,
			Status:      s.Status,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}

	return res, nil
}
