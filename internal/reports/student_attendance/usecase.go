package student_attendance

import (
	"context"

	"github.com/Xlussov/EduCRM-be/internal/domain"
)

type UseCase struct {
	reportRepo  domain.ReportRepository
	studentRepo domain.StudentRepository
}

func NewUseCase(reportRepo domain.ReportRepository, studentRepo domain.StudentRepository) *UseCase {
	return &UseCase{
		reportRepo:  reportRepo,
		studentRepo: studentRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, caller domain.Caller, req Request) (*Response, error) {
	student, err := uc.studentRepo.GetByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}

	if domain.RequiresBranchAccess(caller.Role) && !domain.HasBranchAccess(caller.BranchIDs, student.BranchID) {
		return nil, domain.ErrBranchAccessDenied
	}

	items, err := uc.reportRepo.GetStudentAttendanceHistory(ctx, req.StudentID, req.StartDate, req.EndDate, req.SubjectID)
	if err != nil {
		return nil, err
	}

	var total, attended, missed int
	var respItems []ReportItem

	for _, item := range items {
		respItems = append(respItems, ReportItem{
			Date:        item.Date,
			Time:        item.Time,
			SubjectName: item.SubjectName,
			IsPresent:   item.IsPresent,
			Notes:       item.Notes,
		})

		total++
		if item.IsPresent {
			attended++
		} else {
			missed++
		}
	}

	var percentage float64
	if total > 0 {
		percentage = float64(attended) / float64(total) * 100
	}

	return &Response{
		Items: respItems,
		Summary: Summary{
			TotalLessons:         total,
			Attended:             attended,
			Missed:               missed,
			AttendancePercentage: percentage,
		},
	}, nil
}
