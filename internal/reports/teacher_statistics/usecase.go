package teacher_statistics

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	reportRepo domain.ReportRepository
	userRepo   domain.UserRepository
}

func NewUseCase(reportRepo domain.ReportRepository, userRepo domain.UserRepository) *UseCase {
	return &UseCase{
		reportRepo: reportRepo,
		userRepo:   userRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (*Response, error) {
	if caller.Role == domain.RoleTeacher {
		req.TeacherID = caller.UserID
	} else if caller.Role == domain.RoleAdmin {
		teacherBranchIDs, err := uc.userRepo.GetUserBranchIDs(ctx, req.TeacherID)
		if err != nil {
			return nil, err
		}

		hasAccess := false
		for _, tbID := range teacherBranchIDs {
			if domain.HasBranchAccess(caller.BranchIDs, tbID) {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			return nil, domain.ErrBranchAccessDenied
		}
	}

	stats, err := uc.reportRepo.GetTeacherStatistics(ctx, req.TeacherID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	return &Response{
		ScheduledLessons: stats.ScheduledLessons,
		CompletedLessons: stats.CompletedLessons,
		CancelledLessons: stats.CancelledLessons,
	}, nil
}
