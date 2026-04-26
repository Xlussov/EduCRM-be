package branch_statistics

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	reportRepo domain.ReportRepository
}

func NewUseCase(reportRepo domain.ReportRepository) *UseCase {
	return &UseCase{
		reportRepo: reportRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (*Response, error) {
	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, req.BranchID) {
		return nil, domain.ErrBranchAccessDenied
	}

	stats, err := uc.reportRepo.GetBranchStatistics(ctx, req.BranchID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	return &Response{
		ActiveStudents:       stats.ActiveStudents,
		CompletedLessons:     stats.CompletedLessons,
		CancelledLessons:     stats.CancelledLessons,
		AttendancePercentage: stats.AttendancePercentage,
	}, nil
}
